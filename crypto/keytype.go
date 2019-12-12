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

import (
	"github.com/VirgilSecurity/virgil-sdk-go/crypto/internal/foundation"
)

type KeyType int

// nolint: golint
const (
	DefaultKeyType KeyType = iota
	RSA_2048
	RSA_3072
	RSA_4096
	RSA_8192
	EC_SECP256R1
	EC_SECP384R1
	EC_SECP521R1
	EC_BP256R1
	EC_BP384R1
	EC_BP512R1
	EC_SECP256K1
	EC_CURVE25519
	FAST_EC_X25519
	FAST_EC_ED25519
	PQC
)

var keyTypeMap = map[KeyType]keyAlg{
	DefaultKeyType:  keyType(foundation.AlgIdEd25519),
	RSA_2048:        rsaKeyType{foundation.AlgIdRsa, 2048},
	RSA_3072:        rsaKeyType{foundation.AlgIdRsa, 3072},
	RSA_4096:        rsaKeyType{foundation.AlgIdRsa, 4096},
	RSA_8192:        rsaKeyType{foundation.AlgIdRsa, 8192},
	EC_SECP256R1:    keyType(foundation.AlgIdSecp256r1),
	EC_CURVE25519:   keyType(foundation.AlgIdCurve25519),
	FAST_EC_ED25519: keyType(foundation.AlgIdEd25519),
	PQC:             keyType(foundation.AlgIdPostQuantum),
}

type keyAlg interface {
	AlgID() foundation.AlgId
}

type keyType foundation.AlgId

func (t keyType) AlgID() foundation.AlgId {
	return foundation.AlgId(t)
}

type rsaKeyAlg interface {
	keyAlg
	Len() uint32
}

type rsaKeyType struct {
	t   foundation.AlgId
	len uint32
}

func (t rsaKeyType) AlgID() foundation.AlgId {
	return t.t
}
func (t rsaKeyType) Len() uint32 {
	return t.len
}
