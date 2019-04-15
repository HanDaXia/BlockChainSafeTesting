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
package bccsp

const (
	SM2 = "SM2"
	SM3 = "SM3"
	//This is used for indicating hashopts while doing sm3 before sm2 signature
	// through which hash will actually do nothing
	SM3SIG = "SM3SIG"
	// SM2ReRand SM2 key re-randomization
	SM2ReRand = "SM2"
)

// SHA256Opts contains options relating to SHA-256.
type SM3Opts struct {
}

// Algorithm returns the hash algorithm identifier (to be used).
func (opts *SM3Opts) Algorithm() string {
	return SM3
}

// SHA256Opts contains options relating to SHA-256.
type SM3SIGOpts struct {
}

// Algorithm returns the hash algorithm identifier (to be used).
func (opts *SM3SIGOpts) Algorithm() string {
	return SM3SIG
}

// SM2KeyGenOpts contains options for SM2 key generation.
type SM2KeyGenOpts struct {
	Temporary bool
}

// Algorithm returns the key generation algorithm identifier (to be used).
func (opts *SM2KeyGenOpts) Algorithm() string {
	return SM2
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *SM2KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}

// SM2PKIXPublicKeyImportOpts contains options for SM2 public key importation in PKIX format
type SM2PKIXPublicKeyImportOpts struct {
	Temporary bool
}

// Algorithm returns the key importation algorithm identifier (to be used).
func (opts *SM2PKIXPublicKeyImportOpts) Algorithm() string {
	return SM2
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *SM2PKIXPublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}

// SM2PrivateKeyImportOpts contains options for SM2 secret key importation in DER format
// or PKCS#8 format.
type SM2PrivateKeyImportOpts struct {
	Temporary bool
}

// Algorithm returns the key importation algorithm identifier (to be used).
func (opts *SM2PrivateKeyImportOpts) Algorithm() string {
	return SM2
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *SM2PrivateKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}

// SM2GoPublicKeyImportOpts contains options for SM2 key importation from ecdsa.PublicKey
type SM2GoPublicKeyImportOpts struct {
	Temporary bool
}

// Algorithm returns the key importation algorithm identifier (to be used).
func (opts *SM2GoPublicKeyImportOpts) Algorithm() string {
	return SM2
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *SM2GoPublicKeyImportOpts) Ephemeral() bool {
	return opts.Temporary
}

// SM2ReRandKeyOpts contains options for SM2 key re-randomization.
type SM2ReRandKeyOpts struct {
	Temporary bool
	Expansion []byte
}

// Algorithm returns the key derivation algorithm identifier (to be used).
func (opts *SM2ReRandKeyOpts) Algorithm() string {
	return SM2ReRand
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *SM2ReRandKeyOpts) Ephemeral() bool {
	return opts.Temporary
}

// ExpansionValue returns the re-randomization factor
func (opts *SM2ReRandKeyOpts) ExpansionValue() []byte {
	return opts.Expansion
}
