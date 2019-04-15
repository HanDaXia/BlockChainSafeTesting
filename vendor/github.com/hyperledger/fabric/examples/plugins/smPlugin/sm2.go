/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		SPDX-License-Identifier: Apache-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"crypto/rand"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	sm "github.com/flyinox/crypto/sm/sm2"

	"github.com/hyperledger/fabric/bccsp"
)

type SM2Signature struct {
	R, S *big.Int
}

func MarshalSM2Signature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(SM2Signature{r, s})
}

func UnmarshalSM2Signature(raw []byte) (*big.Int, *big.Int, error) {
	// Unmarshal
	sig := new(SM2Signature)
	_, err := asn1.Unmarshal(raw, sig)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed unmashalling signature [%s]", err)
	}

	// Validate sig
	if sig.R == nil {
		return nil, nil, errors.New("Invalid signature. R must be different from nil.")
	}
	if sig.S == nil {
		return nil, nil, errors.New("Invalid signature. S must be different from nil.")
	}

	if sig.R.Sign() != 1 {
		return nil, nil, errors.New("Invalid signature. R must be larger than zero")
	}
	if sig.S.Sign() != 1 {
		return nil, nil, errors.New("Invalid signature. S must be larger than zero")
	}

	return sig.R, sig.S, nil
}

func signSM2(k *sm.PrivateKey, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	r, s, err := sm.Sign(rand.Reader, k, digest)
	if err != nil {
		return nil, err
	}
	return MarshalSM2Signature(r, s)
}

func verifySM2(k *sm.PublicKey, signature, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
	r, s, err := UnmarshalSM2Signature(signature)
	if err != nil {
		return false, fmt.Errorf("Failed unmashalling signature [%s]", err)
	}
	return sm.Verify(k, digest, r, s), nil
}
