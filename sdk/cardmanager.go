/*
 * Copyright (C) 2015-2018 Virgil Security Inc.
 *
 * Lead Maintainer: Virgil Security Inc. <support@virgilsecurity.com>
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 *   (1) Redistributions of source code must retain the above copyright
 *   notice, this list of conditions and the following disclaimer.
 *
 *   (2) Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions and the following disclaimer in
 *   the documentation and/or other materials provided with the
 *   distribution.
 *
 *   (3) Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived
 *   from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package sdk

import (
	"sync"

	"time"

	"gopkg.in/virgil.v5/cryptoapi"
	"gopkg.in/virgil.v5/errors"
)

type CardManager struct {
	ModelSigner                                  *ModelSigner
	Crypto                                       cryptoapi.CardCrypto
	AccessTokenProvider                          AccessTokenProvider
	CardVerifier                                 CardVerifier
	CardClient                                   *CardClient
	SignCallback                                 func(model *RawSignedModel) error
	onceClient, onceModelSigner, onceCheckParams sync.Once
	paramsError                                  error
}

func NewCardManager(params *CardManagerParams) (*CardManager, error) {
	client := params.CardClient

	if client == nil {
		client = NewCardsClient(params.ApiUrl)
	}
	mgr := &CardManager{
		Crypto:              params.Crypto,
		ModelSigner:         NewModelSigner(params.Crypto),
		SignCallback:        params.SignCallback,
		AccessTokenProvider: params.AccessTokenProvider,
		CardVerifier:        params.CardVerifier,
		CardClient:          client,
	}
	if err := mgr.selfCheck(); err != nil {
		return nil, err
	}
	return mgr, nil
}

func (c *CardManager) GenerateRawCard(cardParams *CardParams) (*RawSignedModel, error) {
	if err := c.selfCheck(); err != nil {
		return nil, err
	}
	if err := cardParams.Validate(false); err != nil {
		return nil, err
	}
	now := time.Now().UTC().Truncate(time.Second)

	model, err := GenerateRawCard(c.Crypto, cardParams, now)

	if err != nil {
		return nil, err
	}
	err = c.getModelSigner().SelfSign(model, cardParams.PrivateKey, cardParams.ExtraFields)

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (c *CardManager) PublishRawSignedModel(rawSignedModel *RawSignedModel, tokenContext *TokenContext, token AccessToken) (*Card, error) {
	if err := c.selfCheck(); err != nil {
		return nil, err
	}

	if c.SignCallback != nil {
		if err := c.SignCallback(rawSignedModel); err != nil {
			return nil, err
		}
	}
	rawCard, err := c.getClient().PublishCard(rawSignedModel, token.StringRepresentation())
	if err != nil {
		return nil, err
	}
	card, err := ParseRawCard(c.Crypto, rawCard, false)

	if err != nil {
		return nil, err
	}

	if err := c.verifyCards(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (c *CardManager) PublishCard(cardParams *CardParams) (*Card, error) {
	if err := c.selfCheck(); err != nil {
		return nil, err
	}
	if err := cardParams.Validate(false); err != nil {
		return nil, err
	}
	tokenContext := &TokenContext{Identity: cardParams.Identity, Operation: "publish"}
	token, err := c.AccessTokenProvider.GetToken(tokenContext)
	if err != nil {
		return nil, err
	}
	identity, err := token.Identity()

	if err != nil {
		return nil, err
	}
	rawSignedModel, err := c.GenerateRawCard(&CardParams{
		Identity:       identity,
		PrivateKey:     cardParams.PrivateKey,
		PublicKey:      cardParams.PublicKey,
		ExtraFields:    cardParams.ExtraFields,
		PreviousCardId: cardParams.PreviousCardId,
	})
	if err != nil {
		return nil, err
	}
	return c.PublishRawSignedModel(rawSignedModel, tokenContext, token)
}

func (c *CardManager) GetCard(cardId string) (*Card, error) {
	if err := c.selfCheck(); err != nil {
		return nil, err
	}
	tokenContext := &TokenContext{Identity: "my_default_identity", Operation: "get"}
	token, err := c.AccessTokenProvider.GetToken(tokenContext)
	if err != nil {
		return nil, err
	}

	rawCard, outdated, err := c.getClient().GetCard(cardId, token.StringRepresentation())

	if err != nil {
		return nil, err
	}
	card, err := ParseRawCard(c.Crypto, rawCard, outdated)
	if err != nil {
		return nil, err
	}
	err = c.verifyCards(card)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (c *CardManager) SearchCards(identity string) ([]*Card, error) {
	if err := c.selfCheck(); err != nil {
		return nil, err
	}
	tokenContext := &TokenContext{Identity: identity, Operation: "search"}
	token, err := c.AccessTokenProvider.GetToken(tokenContext)
	if err != nil {
		return nil, err
	}

	rawCards, err := c.getClient().SearchCards(identity, token.StringRepresentation())
	if err != nil {
		return nil, err
	}

	cards, err := ParseRawCards(c.Crypto, rawCards...)
	if err != nil {
		return nil, err
	}
	err = c.verifyCards(cards...)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c *CardManager) verifyCards(cards ...*Card) error {
	if c.CardVerifier == nil {
		return nil
	}

	for _, card := range cards {
		if err := c.CardVerifier.VerifyCard(card); err != nil {
			return err
		}
	}
	return nil
}

func (c *CardManager) getModelSigner() *ModelSigner {
	c.onceModelSigner.Do(func() {
		if c.ModelSigner == nil {
			c.ModelSigner = &ModelSigner{Crypto: c.Crypto}
		}
	})

	return c.ModelSigner
}

func (c *CardManager) getClient() *CardClient {
	c.onceClient.Do(func() {
		if c.CardClient == nil {
			c.CardClient = &CardClient{}
		}
	})

	return c.CardClient
}

func (c *CardManager) selfCheck() error {
	c.onceCheckParams.Do(func() {
		if c.Crypto == nil {
			c.paramsError = errors.New("Crypto must be set")
			return
		}

		if c.AccessTokenProvider == nil {
			c.paramsError = errors.New("AccessTokenProvider must be set")
			return
		}
	})
	return c.paramsError
}