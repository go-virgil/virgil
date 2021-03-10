//  @license
// --------------------------------------------------------------------------
//  Copyright (C) 2015-2020 Virgil Security, Inc.
//
//  All rights reserved.
//
//  Redistribution and use in source and binary forms, with or without
//  modification, are permitted provided that the following conditions are
//  met:
//
//      (1) Redistributions of source code must retain the above copyright
//      notice, this list of conditions and the following disclaimer.
//
//      (2) Redistributions in binary form must reproduce the above copyright
//      notice, this list of conditions and the following disclaimer in
//      the documentation and/or other materials provided with the
//      distribution.
//
//      (3) Neither the name of the copyright holder nor the names of its
//      contributors may be used to endorse or promote products derived from
//      this software without specific prior written permission.
//
//  THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS OR
//  IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
//  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
//  DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
//  INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
//  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
//  SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
//  HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
//  STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
//  IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
//  POSSIBILITY OF SUCH DAMAGE.
//
//  Lead Maintainer: Virgil Security Inc. <support@virgilsecurity.com>
// --------------------------------------------------------------------------

#ifndef MBEDTLS_CONFIG_H
#define MBEDTLS_CONFIG_H

#include <stddef.h> // size_t

//
//  Common
//
#define MBEDTLS_ERROR_C
#define MBEDTLS_PLATFORM_C

void *vsc_calloc(size_t count, size_t size);
void *vsc_dealloc(size_t size);

#define MBEDTLS_PLATFORM_MEMORY
#define MBEDTLS_PLATFORM_CALLOC vsc_calloc
#define MBEDTLS_PLATFORM_FREE vsc_dealloc

//
//  Required by library vsc::foundation
//
#define MBEDTLS_SHA256_C
#define MBEDTLS_SHA512_C
#define MBEDTLS_CIPHER_C
#define MBEDTLS_AES_C
#define MBEDTLS_GCM_C
#define MBEDTLS_MD_C
#define MBEDTLS_BIGNUM_C
#define MBEDTLS_PKCS1_V21
#define MBEDTLS_OID_C
#define MBEDTLS_RSA_C
#define MBEDTLS_ASN1_PARSE_C
#define MBEDTLS_ASN1_WRITE_C
#define MBEDTLS_GENPRIME
#define MBEDTLS_PLATFORM_ENTROPY
#define MBEDTLS_TIMING_C
#define MBEDTLS_HAVEGE_C
#define MBEDTLS_BASE64_C
#define MBEDTLS_CIPHER_MODE_CBC
#define MBEDTLS_CIPHER_MODE_WITH_PADDING
#define MBEDTLS_CIPHER_PADDING_PKCS7

#if !defined(MBEDTLS_PLATFORM_ENTROPY)
#   define MBEDTLS_NO_PLATFORM_ENTROPY
#endif

#define MBEDTLS_THREADING_C
/* #undef MBEDTLS_THREADING_PTHREAD */
#define MBEDTLS_THREADING_SRWLOCK

//
//  Required by library vsc::pythia
//
#define MBEDTLS_CTR_DRBG_C
#define MBEDTLS_ENTROPY_C

//
//  Required by library vsc::phe and vsc::foundation
//
#define MBEDTLS_ECP_DP_SECP256R1_ENABLED
#define MBEDTLS_ECP_C
#define MBEDTLS_ECDH_C
#define MBEDTLS_ECDSA_C

#if defined(MBEDTLS_ECDSA_C)
#    define MBEDTLS_ECDSA_DETERMINISTIC
#    define MBEDTLS_HMAC_DRBG_C
#endif

//
//  Alternative implementations
//
/* #undef MBEDTLS_SHA256_ALT */
/* #undef MBEDTLS_SHA512_ALT */
/* #undef MBEDTLS_AES_ALT */
/* #undef MBEDTLS_GCM_ALT */

//
//  Non configurable options
//
#define MBEDTLS_NO_DEFAULT_ENTROPY_SOURCES

//
// EC optimizations
//
#define MBEDTLS_ECP_NIST_OPTIM 1
#define MBEDTLS_ECP_FIXED_POINT_OPTIM 1

#include "check_config.h"

#endif /* MBEDTLS_CONFIG_H */
