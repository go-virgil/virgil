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
 *
 */

package crypto

type CardCrypto interface {
	GenerateSignature(data []byte, privateKey PrivateKey) ([]byte, error)
	VerifySignature(data []byte, signature []byte, publicKey PublicKey) error
	ExportPublicKey(publicKey PublicKey) ([]byte, error)
	ImportPublicKey(publicKeySrc []byte) (PublicKey, error)
	GenerateSHA512(data []byte) []byte
}

type AccessTokenSigner interface {
	GenerateTokenSignature(data []byte, privateKey PrivateKey) ([]byte, error)
	VerifyTokenSignature(data []byte, signature []byte, publicKey PublicKey) error
	GetAlgorithm() string
}

type PrivateKeyExporter interface {
	ExportPrivateKey(privateKey PrivateKey) ([]byte, error)
	ImportPrivateKey(data []byte) (privateKey PrivateKey, err error)
}

type PrivateKey interface {
	Identifier() []byte
	PublicKey() PublicKey
}
type PublicKey interface {
	Export() ([]byte, error)
	Identifier() []byte
}
