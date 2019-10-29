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

package cryptogo

import (
	"encoding/binary"
	"hash"
)

//Kdf2 derives length crypto bytes from key and a hash function
func kdf2(key []byte, length int, h func() hash.Hash) []byte {
	kdfHash := h()
	outLen := kdfHash.Size()

	cThreshold := (length + outLen - 1) / outLen
	var counter uint32 = 1
	outOff := 0
	res := make([]byte, length)
	b := make([]byte, 4)
	for i := 0; i < cThreshold; i++ {
		_, err := kdfHash.Write(key)
		if err != nil {
			panic(err)
		}
		binary.BigEndian.PutUint32(b, counter)
		_, err = kdfHash.Write(b)
		if err != nil {
			panic(err)
		}
		counter++
		digest := kdfHash.Sum(nil)

		if length > outLen {
			copy(res[outOff:], digest[:])
			outOff += outLen
			length -= outLen
		} else {
			copy(res[outOff:], digest[:length])
		}
		kdfHash.Reset()
	}
	return res

}
