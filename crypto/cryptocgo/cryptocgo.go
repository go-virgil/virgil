/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2015-2018, Virgil Security, Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 *  Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 *  Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package cryptocgo

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"io"

	"github.com/VirgilSecurity/virgil-sdk-go/crypto"
	"github.com/VirgilSecurity/virgil-sdk-go/crypto/cryptocgo/internal/foundation"
)

type cryptoCgo struct {
	keyType               crypto.KeyType
	UseSha256Fingerprints bool
}

func NewVirgilCrypto() *cryptoCgo {
	return &cryptoCgo{}
}

var (
	signatureKey = []byte("VIRGIL-DATA-SIGNATURE")
	signerIDKey  = []byte("VIRGIL-DATA-SIGNER-ID")
)

func (c *cryptoCgo) SetKeyType(keyType crypto.KeyType) error {
	if _, ok := keyTypeMap[keyType]; !ok {
		return ErrUnsupportedKeyType
	}
	c.keyType = keyType
	return nil
}

func (c *cryptoCgo) generateKeypair(t keyAlg, rnd foundation.Random) (crypto.Keypair, error) {
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	if t.AlgID() == foundation.AlgIdRsa {
		// if keyAlg doesn't cast to rsaKeyAlg we catch panic
		rsaAlg, ok := t.(rsaKeyAlg)
		if !ok {
			panic("incorrect used RSA algorithm")
		}
		kp.SetRsaParams(rsaAlg.Len())
	}
	kp.SetRandom(rnd)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}

	sk, err := kp.GeneratePrivateKey(t.AlgID())
	if err != nil {
		return nil, err
	}

	pk, err := sk.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	id, err := c.calculateFingerprint(pk)
	if err != nil {
		return nil, err
	}

	return &keypair{
		privateKey: privateKey{receiverID: id, key: sk},
		publicKey:  publicKey{receiverID: id, key: pk},
	}, nil
}

func (c *cryptoCgo) GenerateKeypairForType(t crypto.KeyType) (crypto.Keypair, error) {
	keyType, ok := keyTypeMap[t]
	if !ok {
		return nil, ErrUnsupportedKeyType
	}
	return c.generateKeypair(keyType, random)
}

func (c *cryptoCgo) GenerateKeypair() (crypto.Keypair, error) {
	return c.GenerateKeypairForType(c.keyType)
}

func (c *cryptoCgo) GenerateKeypairFromKeyMaterialForType(t crypto.KeyType, keyMaterial []byte) (crypto.Keypair, error) {
	l := uint32(len(keyMaterial))
	if l < foundation.KeyMaterialRngKeyMaterialLenMin || l > foundation.KeyMaterialRngKeyMaterialLenMax {
		return nil, ErrInvalidSeedSize
	}

	keyType, ok := keyTypeMap[t]
	if !ok {
		return nil, ErrUnsupportedKeyType
	}

	rnd := foundation.NewKeyMaterialRng()
	rnd.ResetKeyMaterial(keyMaterial)
	defer delete(rnd)

	return c.generateKeypair(keyType, rnd)
}

func (c *cryptoCgo) GenerateKeypairFromKeyMaterial(keyMaterial []byte) (crypto.Keypair, error) {
	return c.GenerateKeypairFromKeyMaterialForType(c.keyType, keyMaterial)
}

func (c *cryptoCgo) Random(len int) ([]byte, error) {
	return random.Random(uint32(len))
}

func (c *cryptoCgo) ImportPrivateKey(data []byte) (crypto.PrivateKey, error) {
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	kp.SetRandom(random)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}

	sk, err := kp.ImportPrivateKey(data)
	if err != nil {
		return nil, err
	}

	pk, err := sk.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	id, err := c.calculateFingerprint(pk)
	if err != nil {
		return nil, err
	}

	return privateKey{receiverID: id, key: sk}, nil
}

func (c *cryptoCgo) ImportPublicKey(data []byte) (crypto.PublicKey, error) {
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	kp.SetRandom(random)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}
	pk, err := kp.ImportPublicKey(data)
	if err != nil {
		return nil, err
	}

	id, err := c.calculateFingerprint(pk)
	if err != nil {
		return nil, err
	}

	return publicKey{receiverID: id, key: pk}, nil
}

func (c *cryptoCgo) ExportPrivateKey(key crypto.PrivateKey) ([]byte, error) {
	sk, ok := key.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	kp.SetRandom(random)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}
	return kp.ExportPrivateKey(sk.key)
}

func (c *cryptoCgo) ExportPublicKey(key crypto.PublicKey) ([]byte, error) {
	pk, ok := key.(publicKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	kp.SetRandom(random)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}

	return kp.ExportPublicKey(pk.key)
}

func (c *cryptoCgo) calculateFingerprint(key foundation.PublicKey) ([]byte, error) {
	kp := foundation.NewKeyProvider()
	defer delete(kp)

	kp.SetRandom(random)
	if err := kp.SetupDefaults(); err != nil {
		return nil, err
	}

	data, err := kp.ExportPublicKey(key)
	if err != nil {
		return nil, err
	}

	if c.UseSha256Fingerprints {
		hash := sha256.Sum256(data)
		return hash[:], nil
	}

	hash := sha512.Sum512(data)
	return hash[:8], nil
}

func (c *cryptoCgo) Encrypt(data []byte, recipients ...crypto.PublicKey) ([]byte, error) {
	cipher, err := c.setupEncryptCipher(recipients)
	if err != nil {
		return nil, err
	}
	defer delete(cipher)

	if err := cipher.StartEncryption(); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.Grow(len(data))

	dst := NewEncryptWriter(NopWriteCloser(buffer), cipher)
	src := bytes.NewReader(data)
	if err := copyClose(dst, src); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *cryptoCgo) Decrypt(data []byte, key crypto.PrivateKey) ([]byte, error) {
	sk, ok := key.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	cipher := c.setupCipher()
	defer delete(cipher)

	if err := cipher.StartDecryptionWithKey(sk.receiverID, sk.key, nil); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	dr := NewDecryptReader(bytes.NewReader(data), cipher)
	if _, err := io.Copy(buf, dr); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *cryptoCgo) EncryptStream(in io.Reader, out io.Writer, recipients ...crypto.PublicKey) (err error) {
	cipher, err := c.setupEncryptCipher(recipients)
	if err != nil {
		return err
	}
	defer delete(cipher)

	if err := cipher.StartEncryption(); err != nil {
		return err
	}

	dst := NewEncryptWriter(NopWriteCloser(out), cipher)
	if err := copyClose(dst, in); err != nil {
		return err
	}

	return nil
}

func (c *cryptoCgo) DecryptStream(in io.Reader, out io.Writer, key crypto.PrivateKey) (err error) {
	sk, ok := key.(privateKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	cipher := c.setupCipher()
	defer delete(cipher)

	if err = cipher.StartDecryptionWithKey(sk.receiverID, sk.key, nil); err != nil {
		return err
	}

	dr := NewDecryptReader(in, cipher)
	if _, err := io.Copy(out, dr); err != nil {
		return err
	}

	return nil
}

func (c *cryptoCgo) Sign(data []byte, signer crypto.PrivateKey) ([]byte, error) {
	sk, ok := signer.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	s := foundation.NewSigner()
	h := foundation.NewSha512()
	defer delete(s, h)

	s.SetRandom(random)
	s.SetHash(h)
	s.Reset()
	s.AppendData(data)

	return s.Sign(sk.key)
}

func (c *cryptoCgo) VerifySignature(data []byte, signature []byte, key crypto.PublicKey) error {
	pk, ok := key.(publicKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	v := foundation.NewVerifier()
	defer delete(v)

	if err := v.Reset(signature); err != nil {
		return err
	}
	v.AppendData(data)

	if v.Verify(pk.key) {
		return nil
	}
	return ErrSignVerification
}

func (c *cryptoCgo) SignStream(in io.Reader, signer crypto.PrivateKey) ([]byte, error) {
	sk, ok := signer.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	s := foundation.NewSigner()
	h := foundation.NewSha512()
	defer delete(s, h)

	s.SetRandom(random)
	s.SetHash(h)
	s.Reset()
	if _, err := io.Copy(appenderWriter{s}, in); err != nil {
		return nil, err
	}

	return s.Sign(sk.key)
}

func (c *cryptoCgo) VerifyStream(in io.Reader, signature []byte, key crypto.PublicKey) error {
	pk, ok := key.(publicKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	v := foundation.NewVerifier()
	defer delete(v)

	if err := v.Reset(signature); err != nil {
		return err
	}
	if _, err := io.Copy(appenderWriter{v}, in); err != nil {
		return err
	}

	if v.Verify(pk.key) {
		return nil
	}
	return ErrSignVerification
}

func (c *cryptoCgo) SignThenEncrypt(data []byte, signer crypto.PrivateKey, recipients ...crypto.PublicKey) ([]byte, error) {
	sk, ok := signer.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	cipher, err := c.setupEncryptCipher(recipients)
	if err != nil {
		return nil, err
	}
	h := foundation.NewSha512()

	defer delete(cipher, h)

	cipher.SetSignerHash(h)
	if err = cipher.AddSigner(sk.receiverID, sk.key); err != nil {
		return nil, err
	}
	if err = cipher.StartSignedEncryption(uint32(len(data))); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	dst := NewEncryptWriter(NopWriteCloser(buffer), cipher)
	src := bytes.NewReader(data)
	if err = copyClose(dst, src); err != nil {
		return nil, err
	}

	buf, err := cipher.PackMessageInfoFooter()
	if err != nil {
		return nil, err
	}
	if _, err = buffer.Write(buf); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *cryptoCgo) DecryptThenVerify(
	data []byte,
	decryptionKey crypto.PrivateKey,
	verifierKeys ...crypto.PublicKey,
) (_ []byte, err error) {
	sk, ok := decryptionKey.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	cipher := c.setupCipher()
	defer delete(cipher)

	if err := cipher.StartDecryptionWithKey(sk.receiverID, sk.key, nil); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, NewDecryptReader(bytes.NewReader(data), cipher)); err != nil {
		return nil, err
	}
	if err := c.verifyCipherSign(cipher, verifierKeys); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *cryptoCgo) SignThenEncryptStream(
	in io.Reader,
	out io.Writer,
	streamSize int,
	signer crypto.PrivateKey,
	recipients ...crypto.PublicKey,
) (err error) {
	if streamSize < 0 {
		return ErrStreamSizeIncorrect
	}

	sk, ok := signer.(privateKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	var (
		cipher *foundation.RecipientCipher
		h      foundation.Hash
	)
	defer delete(cipher, h)

	cipher, err = c.setupEncryptCipher(recipients)
	if err != nil {
		return err
	}

	h = foundation.NewSha512()
	cipher.SetSignerHash(h)
	if err = cipher.AddSigner(sk.Identifier(), sk.key); err != nil {
		return err
	}
	if err = cipher.StartSignedEncryption(uint32(streamSize)); err != nil {
		return err
	}

	dst := NewEncryptWriter(NopWriteCloser(out), cipher)
	if err = copyClose(dst, in); err != nil {
		return err
	}

	buf, err := cipher.PackMessageInfoFooter()
	if err != nil {
		return err
	}
	if _, err = out.Write(buf); err != nil {
		return err
	}

	return nil
}

func (c *cryptoCgo) DecryptThenVerifyStream(
	in io.Reader,
	out io.Writer,
	decryptionKey crypto.PrivateKey,
	verifierKeys ...crypto.PublicKey,
) error {
	sk, ok := decryptionKey.(privateKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	cipher := c.setupCipher()
	defer delete(cipher)

	if err := cipher.StartDecryptionWithKey(sk.Identifier(), sk.key, nil); err != nil {
		return err
	}

	if _, err := io.Copy(out, NewDecryptReader(in, cipher)); err != nil {
		return err
	}

	return c.verifyCipherSign(cipher, verifierKeys)
}

func (c *cryptoCgo) SignAndEncrypt(data []byte, signer crypto.PrivateKey, recipients ...crypto.PublicKey) (_ []byte, err error) {
	var (
		cipher *foundation.RecipientCipher
		params *foundation.MessageInfoCustomParams
	)
	defer delete(cipher, params)

	cipher, err = c.setupEncryptCipher(recipients)
	if err != nil {
		return nil, err
	}

	sign, err := c.Sign(data, signer)
	if err != nil {
		return nil, err
	}

	params = cipher.CustomParams()
	params.AddData(signatureKey, sign)
	params.AddData(signerIDKey, signer.Identifier())

	if err := cipher.StartEncryption(); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.Grow(len(data))

	dst := NewEncryptWriter(NopWriteCloser(buffer), cipher)
	src := bytes.NewReader(data)
	if err := copyClose(dst, src); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *cryptoCgo) DecryptAndVerify(data []byte, decryptionKey crypto.PrivateKey, verifierKeys ...crypto.PublicKey) (_ []byte, err error) {
	sk, ok := decryptionKey.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	var (
		cipher *foundation.RecipientCipher
		params *foundation.MessageInfoCustomParams
	)
	defer delete(cipher, params)

	cipher = c.setupCipher()
	if err = cipher.StartDecryptionWithKey(sk.Identifier(), sk.key, nil); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(buffer, NewDecryptReader(bytes.NewReader(data), cipher)); err != nil {
		return nil, err
	}

	params = cipher.CustomParams()
	signerID, err := params.FindData(signerIDKey)
	if err != nil {
		return nil, err
	}

	sign, err := params.FindData(signatureKey)
	if err != nil {
		return nil, err
	}

	k, err := findVerifyKey(signerID, verifierKeys)
	if err != nil {
		return nil, err
	}

	if err := c.VerifySignature(buffer.Bytes(), sign, k); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *cryptoCgo) ExtractPublicKey(key crypto.PrivateKey) (crypto.PublicKey, error) {
	sk, ok := key.(privateKey)
	if !ok {
		return nil, ErrUnsupportedParameter
	}

	pk, err := sk.key.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	return publicKey{sk.Identifier(), pk}, nil
}

func (c *cryptoCgo) Hash(data []byte, t HashType) ([]byte, error) {
	hf, ok := hashMap[t]
	if !ok {
		return nil, ErrUnsupportedHashType
	}
	hash := hf().Hash(data)

	return hash, nil
}

func (c *cryptoCgo) verifyCipherSign(cipher *foundation.RecipientCipher, verifierKeys []crypto.PublicKey) error {
	var (
		signerInfoList *foundation.SignerInfoList
		signInfo       *foundation.SignerInfo
	)
	defer delete(signerInfoList, signInfo)

	if !cipher.IsDataSigned() {
		return ErrSignNotFound
	}

	signerInfoList = cipher.SignerInfos()
	if !signerInfoList.HasItem() {
		return ErrSignNotFound
	}

	signInfo = signerInfoList.Item()
	k, err := findVerifyKey(signInfo.SignerId(), verifierKeys)
	if err != nil {
		return err
	}
	pk, ok := k.(publicKey)
	if !ok {
		return ErrUnsupportedParameter
	}

	if cipher.VerifySignerInfo(signInfo, pk.key) {
		return nil
	}
	return ErrSignVerification
}

func findVerifyKey(signerID []byte, verifierKeys []crypto.PublicKey) (crypto.PublicKey, error) {
	for _, r := range verifierKeys {
		//TODO: check that it's really need
		if subtle.ConstantTimeCompare(signerID, r.Identifier()) == 1 {
			return r, nil
		}
	}
	return nil, ErrSignNotFound
}

func (c *cryptoCgo) setupCipher() *foundation.RecipientCipher {
	aesGcm := foundation.NewAes256Gcm()
	cipher := foundation.NewRecipientCipher()
	defer delete(aesGcm)

	cipher.SetEncryptionCipher(aesGcm)
	cipher.SetRandom(random)

	return cipher
}

func (c *cryptoCgo) setupEncryptCipher(recipients []crypto.PublicKey) (*foundation.RecipientCipher, error) {
	cipher := c.setupCipher()

	if err := c.setupRecipients(cipher, recipients); err != nil {
		return nil, err
	}
	return cipher, nil
}

func (c *cryptoCgo) setupRecipients(cipher *foundation.RecipientCipher, recipients []crypto.PublicKey) error {
	for _, r := range recipients {
		pk, ok := r.(publicKey)
		if !ok {
			return ErrUnsupportedParameter
		}
		cipher.AddKeyRecipient(r.Identifier(), pk.key)
	}
	if err := cipher.StartEncryption(); err != nil {
		return err
	}
	return nil
}

func copyClose(dst io.WriteCloser, src io.Reader) error {
	_, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	return dst.Close()
}

type appender interface {
	AppendData(b []byte)
}

type appenderWriter struct {
	a appender
}

func (aw appenderWriter) Write(d []byte) (int, error) {
	aw.a.AppendData(d)
	return len(d), nil
}