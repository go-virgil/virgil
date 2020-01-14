 # Virgil Security Go SDK

[![Build Status](https://travis-ci.com/VirgilSecurity/virgil-sdk-go.svg?branch=master)](https://travis-ci.com/VirgilSecurity/virgil-sdk-go)
[![GitHub license](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](https://github.com/VirgilSecurity/virgil/blob/master/LICENSE)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/VirgilSecurity/virgil-sdk-go)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/VirgilSecurity/virgil-sdk-go)


[Introduction](#introduction) | [SDK Features](#sdk-features) | [Library purposes](#library-purposes) | [Installation](#installation) | [Configure SDK](#configure-sdk) | [Usage Examples](#usage-examples) | [Docs](#docs) | [Support](#support)



## Introduction

<a href="https://developer.virgilsecurity.com/docs"><img width="230px" src="https://cdn.virgilsecurity.com/assets/images/github/logos/virgil-logo-red.png" align="left" hspace="10" vspace="6"></a> [Virgil Security](https://virgilsecurity.com) provides a set of APIs for adding security to any application. In a few simple steps you can encrypt communication, securely store data, provide passwordless login, and ensure data integrity.

The Virgil SDK allows developers to get up and running with Virgil API quickly and add full end-to-end security to their existing digital solutions to become HIPAA and GDPR compliant and more.

## SDK Features
- Communicate with [Virgil Cards Service][_cards_service]
- Manage users' public keys
- Encrypt, sign, decrypt and verify data
- Store private keys in secure local storage
- Use Virgil [Crypto library][_virgil_crypto]
- Use your own Crypto

## Library purposes
* Asymmetric Key Generation
* Encryption/Decryption of data and streams
* Generation/Verification of digital signatures
* PFS (Perfect Forward Secrecy)

## Installation

The Virgil Go SDK is provided as a package named virgil. The package is distributed via github. Also in this guide, you will find one more package called Virgil Crypto (Virgil Crypto Library) that is used by the SDK to perform cryptographic operations.

The package is available for Go 1.12 or newer.

Installing the package:

- go get -u github.com/VirgilSecurity/virgil-sdk-go/v6/sdk

### Virgil Crypto Library

Virgil Crypto library is a wrapper for [C Crypto library](https://github.com/VirgilSecurity/virgil-crypto-c). Virgil Crypto library uses prebuilded C library for better experience for the most popular operation system:

- MacOS amd64 with X Code 11
- Windows amd64 with mingw64
- Linux amd64 gcc >= 5 and clang >=7
- For support older version linux amd64 gcc < 5 and clang < 7  with 2.14 Linux kernel use tag `legacy`

## Configure SDK

This section contains guides on how to set up Virgil Core SDK modules for authenticating users, managing Virgil Cards and storing private keys.

### Setup authentication

Set up user authentication with tokens that are based on the [JSON Web Token standard](https://jwt.io/) with some Virgil modifications.

In order to make calls to Virgil Services (for example, to publish user's Card on Virgil Cards Service), you need to have a JSON Web Token ("JWT") that contains the user's `identity`, which is a string that uniquely identifies each user in your application.

Credentials that you'll need:

|Parameter|Description|
|--- |--- |
|App ID|ID of your Application at [Virgil Dashboard](https://dashboard.virgilsecurity.com)|
|App Key ID|A unique string value that identifies your account at the Virgil developer portal|
|App Key|A Private Key that is used to sign API calls to Virgil Services. For security, you will only be shown the App Key when the key is created. Don't forget to save it in a secure location for the next step|

#### Set up JWT provider on Client side

Use these lines of code to specify which JWT generation source you prefer to use in your project:

```go
package main

import (
	"github.com/VirgilSecurity/virgil-sdk-go/v6/sdk"
)

func main() {
	authenticatedQueryToServerSide := func(context *sdk.TokenContext) (string, error) {
		// Get generated token from server-side
		return "eyJraWQiOiI3MGI0NDdlMzIxZjNhMGZkIiwidHlwIjoiSldUIiwiYWxnIjoiVkVEUzUxMiIsImN0eSI6InZpcmdpbC1qd3Q7dj0xIn0.eyJleHAiOjE1MTg2OTg5MTcsImlzcyI6InZpcmdpbC1iZTAwZTEwZTRlMWY0YmY1OGY5YjRkYzg1ZDc5Yzc3YSIsInN1YiI6ImlkZW50aXR5LUFsaWNlIiwiaWF0IjoxNTE4NjEyNTE3fQ.MFEwDQYJYIZIAWUDBAIDBQAEQP4Yo3yjmt8WWJ5mqs3Yrqc_VzG6nBtrW2KIjP-kxiIJL_7Wv0pqty7PDbDoGhkX8CJa6UOdyn3rBWRvMK7p7Ak", nil
	}

	// Setup AccessTokenProvider
	accessTokenProvider := sdk.NewCachingStringJwtProvider(authenticatedQueryToServerSide)
}

```

#### Generate JWT on Server side

Next, you'll need to set up the `JwtGenerator` and generate a JWT using the Virgil SDK.

Here is an example of how to generate a JWT:

```go
	// App Key (you got this Key at Virgil Dashboard)
	privateKeyStr := []byte("MIGhMF0GCSqGSIb3DQEFDTBQMC8GCSqGSIb3DQEFDDAiBBC7Sg/DbNzhJ/uakTva")

	// Crypto library imports a private key into a necessary format
	privateKey, err := crypto.ImportPrivateKey(privateKeyStr, "")
	if err != nil {
		//handle error
	}

	// use your App Credentials you got at Virgil Dashboard:
	appId := "be00e10e4e1f4bf58f9b4dc85d79c77a" // App ID
	appKeyId := "70b447e321f3a0fd"           	// App Key ID
	ttl := time.Hour                            // 1 hour (JWT's lifetime)

	// setup JWT generator with necessary parameters:
	jwtGenerator := sdk.NewJwtGenerator(privateKey, appKeyId, tokenSigner, appId, ttl)

	// generate JWT for a user
	// remember that you must provide each user with his unique JWT
	// each JWT contains unique user's identity (in this case - Alice)
	// identity can be any value: name, email, some id etc.
	identity := "Alice"
	token, err := jwtGenerator.GenerateToken(identity, nil)

	if err != nil {
		//handle error
	}

	// as result you get users JWT, it looks like this: "eyJraWQiOiI3MGI0NDdlMzIxZjNhMGZkIiwidHlwIjoiSldUIiwiYWxnIjoiVkVEUzUxMiIsImN0eSI6InZpcmdpbC1qd3Q7dj0xIn0.eyJleHAiOjE1MTg2OTg5MTcsImlzcyI6InZpcmdpbC1iZTAwZTEwZTRlMWY0YmY1OGY5YjRkYzg1ZDc5Yzc3YSIsInN1YiI6ImlkZW50aXR5LUFsaWNlIiwiaWF0IjoxNTE4NjEyNTE3fQ.MFEwDQYJYIZIAWUDBAIDBQAEQP4Yo3yjmt8WWJ5mqs3Yrqc_VzG6nBtrW2KIjP-kxiIJL_7Wv0pqty7PDbDoGhkX8CJa6UOdyn3rBWRvMK7p7Ak"
	// you can provide users with JWT at registration or authorization steps
	// Send a JWT to client-side
	jwtString := token.String()
```

For this subsection we've created a sample backend that demonstrates how you can set up your backend to generate the JWTs. To set up and run the sample backend locally, head over to your GitHub repo of choice:

[Node.js](https://github.com/VirgilSecurity/sample-backend-nodejs) | [Golang](https://github.com/VirgilSecurity/sample-backend-go) | [PHP](https://github.com/VirgilSecurity/sample-backend-php) | [Java](https://github.com/VirgilSecurity/sample-backend-java) | [Python](https://github.com/VirgilSecurity/virgil-sdk-python/tree/master#sample-backend-for-jwt-generation)
 and follow the instructions in README.
 
### Setup Card Verifier

Virgil Card Verifier helps you automatically verify signatures of a user's Card, for example when you get a Card from Virgil Cards Service.

By default, `VirgilCardVerifier` verifies only two signatures - those of a Card owner and Virgil Cards Service.

Set up `VirgilCardVerifier` with the following lines of code:

```go
	cardVerifier:= sdk.NewVirgilCardVerifier()
```

### Set up Card Manager

This subsection shows how to set up a Card Manager module to help you manage users' public keys.

With Card Manager you can:
- specify an access Token (JWT) Provider.
- specify a Crypto Library that you’re planning to use for crypto operations.
- specify a Card Verifier used to verify signatures of your users, your App Server, Virgil Services (optional).

Use the following lines of code to set up the Card Manager:

```go
cardManager := sdk.NewCardManager(accessTokenProvider)
```

### Setup Key storage to store private keys

This subsection shows how to set up a `VSSKeyStorage` using Virgil SDK in order to save private keys after their generation.

Here is an example of how to set up the `VSSKeyStorage` class:

```go
	// Generate a private key
	privateKey, err := crypto.GenerateKeypair()
	if err != nil {
		//handle error
	}

	// Setup PrivateKeyStorage
	privateKeyStorage := storage.NewVirgilPrivateKeyStorage(&storage.FileStorage{Root:"~/keys/"})

	// Store a private key with a name, for example Alice
	err = privateKeyStorage.Store(privateKey,"Alice", nil)
	if err != nil {
		//handle error
	}

	// To load Alice private key use the following code lines:
	key, meta, err := privateKeyStorage.Load("Alice")

	// Delete a private key
	err = privateKeyStorage.Delete("Alice")
```

## Usage Examples

Before you start practicing with the usage examples, make sure that the SDK is configured. Check out our [SDK configuration guides][#configure-sdk] for more information.

### Generate and publish user's Cards with public keys inside at Cards Service

Use the following lines of code to create a user's Card with a public key inside and publish it at Virgil Cards Service:

```go
import (
	"github.com/VirgilSecurity/virgil-sdk-go/v6/crypto"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/sdk"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/session"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/storage"
)

var (
	AppKey   = "{YOUR_APP_KEY}"
	AppKeyID = "{YOUR_APP_KEY_ID}"
	AppID    = "{YOU_APP_ID}"
)

func main() {
	var crypto crypto.Crypto
	const identity = "Alice"

	// generate a key pair
	keypair, err := crypto.GenerateKeypair()
	if err != nil {
		//handle error
	}

	privateKeyStorage := storage.NewVirgilPrivateKeyStorage(&storage.FileStorage{})
	// save a private key into key storage
	err = privateKeyStorage.Store(keypair, identity, nil)
	if err != nil {
		//handle error
	}

	appKey, err := crypto.ImportPrivateKey([]byte(AppKey))
	if err != nil {
		//handle error
	}

	cardManager := sdk.NewCardManager(session.NewGeneratorJwtProvider(session.JwtGenerator{
		AppKeyID: AppKeyID,
		AppKey:   appKey,
		AppID:    AppID,
	}))

	// publish user's on the Cards Service
	card, err := cardManager.PublishCard(&sdk.CardParams{
		PrivateKey: keypair,
		Identity:   identity,
	})
	if err != nil {
		//handle error
	}
}
```

### Sign then encrypt data

Virgil SDK allows you to use a user's private key and their Virgil Cards to sign and encrypt any kind of data.

In the following example, we load a private key from a customized Key Storage and get recipient's Card from the Virgil Cards Service. Recipient's Card contains a public key which we will use to encrypt the data and verify a signature.

```go
import (
	"github.com/VirgilSecurity/virgil-sdk-go/v6/sdk"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/sdk/crypto"
	"github.com/VirgilSecurity/virgil-sdk-go/v6/sdk/storage"
)


func main() {
	var crypto crypto.Crypto
	messageToEncrypt := []byte("Hello, Bob!")

	privateKeyStorage := storage.NewVirgilPrivateKeyStorage(&storage.FileStorage{})

	// prepare a user's private key from a device storage
	alicePrivateKey, _, err := privateKeyStorage.Load("Alice")
	if err != nil{
		//handle error
	}

	// using cardManager search for Bob's cards on Cards Service
	cards, err := cardManager.SearchCards("Bob")
	if err != nil{
		//handle error
	}

	// sign a message with a private key then encrypt using Bob's public keys
	encryptedMessage, err := crypto.SignThenEncrypt(messageToEncrypt, alicePrivateKey, cards.ExtractPublicKeys()...)
	if err != nil{
		//handle error
	}
}

```

### Decrypt data and verify signature

Once the user receives the signed and encrypted message, they can decrypt it with their own private key and verify the signature with the sender's Card:

```go
import "github.com/VirgilSecurity/virgil-sdk-go/v6/sdk/crypto"


func main() {
	var crypto crypto.Crypto

	privateKeyStorage := storage.NewVirgilPrivateKeyStorage(&storage.FileStorage{})

	// prepare a user's private key
	bobPrivateKey, err := privateKeyStorage.Load("Bob")
	if err != nil{
		//handle error
	}

	// using cardManager search for Alice's cards on Cards Service
	aliceCards, err := cardManager.SearchCards("Alice")
	if err != nil{
		//handle error
	}

	// decrypt with a private key and verify using one of Alice's public keys
	decryptedMessage, err := crypto.DecryptThenVerify(encryptedMessage, bobPrivateKey, cards.ExtractPublicKeys()...)
	if err != nil{
		//handle error
	}
}
```

### Get a Card by its ID

Use the following lines of code to get a user's card from Virgil Cloud by its ID:

```go
card, err := cardManager.GetCard("f4bf9f7fcbedaba0392f108c59d8f4a38b3838efb64877380171b54475c2ade8")
if err != nil {
	//handle error
}
```

### Get a Card by user's identity

For a single user, use the following lines of code to get a user's Card by a user's `identity`:

```go
cards, err := cardManager.SearchCards("Bob")
if err != nil {
	//handle error
}
```

## Docs

Virgil Security has a powerful set of APIs, and the [Developer Documentation](https://developer.virgilsecurity.com/) can get you started today.

In order to use the Virgil Core SDK with your application, you will need to first configure your application. By default, the SDK will attempt to look for Virgil-specific settings in your application but you can change it during SDK configuration.

## License

This library is released under the [3-clause BSD License](LICENSE).

## Support

Our developer support team is here to help you. Find out more information on our [Help Center](https://help.virgilsecurity.com/).

You can find us on [Twitter](https://twitter.com/VirgilSecurity) or send us email support@VirgilSecurity.com.

Also, get extra help from our support team on [Slack](https://virgilsecurity.com/join-community).

[_virgil_crypto]: https://github.com/VirgilSecurity/virgil-crypto-go/tree/master
[_cards_service]: https://developer.virgilsecurity.com/docs/api-reference/cards-service
