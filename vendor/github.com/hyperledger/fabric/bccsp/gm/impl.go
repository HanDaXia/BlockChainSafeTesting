/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package gm

import (
	"hash"

	"crypto/rand"
	"fmt"
	"github.com/flyinox/crypto/sm/sm2"
	"github.com/flyinox/crypto/sm/sm3"
	"github.com/flyinox/crypto/x509"
	"math/big"

	origx509 "crypto/x509"

	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/sw"
	"github.com/hyperledger/fabric/common/flogging"
	coreconfig "github.com/hyperledger/fabric/core/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	logger = flogging.MustGetLogger("plugin_sm")
)

type impl struct {
	sw bccsp.BCCSP
	ks bccsp.KeyStore
}

func getKeyStoreDir(keystore string) (string, error) {
	viperconfig := viper.New()
	viperconfig.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viperconfig.SetEnvKeyReplacer(replacer)

	if keystore == "" {
		switch strings.Split(filepath.Base(os.Args[0]), "_")[0] {
		case "peer":
			viperconfig.SetEnvPrefix("CORE")
			coreconfig.InitViper(viperconfig, "core")
			if err := viperconfig.ReadInConfig(); err != nil {
				return "", fmt.Errorf("Error reading configuration: %s", err)
			}
			return filepath.Join(viper.GetString("peer.mspConfigPath"), "keystore"), nil
		case "orderer":
			viperconfig.SetEnvPrefix("ORDERER")
			coreconfig.InitViper(viperconfig, "orderer")
			if err := viperconfig.ReadInConfig(); err != nil {
				return "", fmt.Errorf("Error reading configuration: %s", err)
			}
			return filepath.Join(viperconfig.GetString("General.LocalMSPDir"), "keystore"), nil
		default:
			return "", nil
		}
	}

	return keystore, nil
}

// New returns a new instance of the BCCSP implementation
func New(keystore string) (bccsp.BCCSP, error) {

	var (
		swCsp bccsp.BCCSP
		ks    bccsp.KeyStore
		err   error
	)

	keyStorePath, err := getKeyStoreDir(keystore)
	if err != nil {
		return  nil, err
	}

	if keyStorePath == ""{
		ks = sw.NewDummyKeyStore()
		swCsp,err = sw.NewDefaultSecurityLevelWithKeystore(ks)
		if err != nil {
			return nil, err
		}
		return &impl{sw: swCsp, ks: ks}, nil
	}

	swCsp, err = sw.NewDefaultSecurityLevel(keyStorePath)
	if err != nil {
		return nil, err
	}

	ks, err = NewFileBasedKeyStore(nil, keyStorePath, false)
	if err != nil {
		return nil, err
	}

	return &impl{sw: swCsp, ks: ks}, nil
}

// KeyGen generates a key using opts.
func (csp *impl) KeyGen(opts bccsp.KeyGenOpts) (k bccsp.Key, err error) {
	switch opts.(type) {
	case *bccsp.SM2KeyGenOpts:
		privKey, err := sm2.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("Failed generating SM2 key for : [%s]", err)
		}

		k = &sm2PrivateKey{privKey}

		if !opts.Ephemeral() {
			// Store the key
			err = csp.ks.StoreKey(k)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed storing imported key with opts [%v]", opts)
			}
		}

		return k, nil
	default:
		return csp.sw.KeyGen(opts)
	}

}

// KeyDeriv derives a key from k using opts.
// The opts argument should be appropriate for the primitive used.
func (csp *impl) KeyDeriv(k bccsp.Key, opts bccsp.KeyDerivOpts) (dk bccsp.Key, err error) {
	switch k := k.(type) {
	case *sm2PrivateKey:
		// Validate opts
		if opts == nil {
			return nil, errors.New("Invalid opts parameter. It must not be nil.")
		}

		sm2K := k

		switch opts.(type) {
		// Re-randomized an ECDSA private key
		case *bccsp.SM2ReRandKeyOpts:
			reRandOpts := opts.(*bccsp.SM2ReRandKeyOpts)
			tempSK := &sm2.PrivateKey{
				PublicKey: sm2.PublicKey{
					Curve: sm2K.privKey.Curve,
					X:     new(big.Int),
					Y:     new(big.Int),
				},
				D: new(big.Int),
			}

			var k = new(big.Int).SetBytes(reRandOpts.ExpansionValue())
			var one = new(big.Int).SetInt64(1)
			n := new(big.Int).Sub(sm2K.privKey.Params().N, one)
			k.Mod(k, n)
			k.Add(k, one)

			tempSK.D.Add(sm2K.privKey.D, k)
			tempSK.D.Mod(tempSK.D, sm2K.privKey.PublicKey.Params().N)

			// Compute temporary public key
			tempX, tempY := sm2K.privKey.PublicKey.ScalarBaseMult(k.Bytes())
			tempSK.PublicKey.X, tempSK.PublicKey.Y =
				tempSK.PublicKey.Add(
					sm2K.privKey.PublicKey.X, sm2K.privKey.PublicKey.Y,
					tempX, tempY,
				)

			// Verify temporary public key is a valid point on the reference curve
			isOn := tempSK.Curve.IsOnCurve(tempSK.PublicKey.X, tempSK.PublicKey.Y)
			if !isOn {
				return nil, errors.New("Failed temporary public key IsOnCurve check.")
			}

			return &sm2PrivateKey{tempSK}, nil
		default:
			return nil, fmt.Errorf("Unsupported 'KeyDerivOpts' provided [%v]", opts)
		}
	case *sm2PublicKey:
		// Validate opts
		if opts == nil {
			return nil, errors.New("Invalid opts parameter. It must not be nil.")
		}

		sm2K := k

		switch opts.(type) {
		// Re-randomized an ECDSA private key
		case *bccsp.SM2ReRandKeyOpts:
			reRandOpts := opts.(*bccsp.SM2ReRandKeyOpts)
			tempSK := &sm2.PublicKey{
				Curve: sm2K.pubKey.Curve,
				X:     new(big.Int),
				Y:     new(big.Int),
			}

			var k = new(big.Int).SetBytes(reRandOpts.ExpansionValue())
			var one = new(big.Int).SetInt64(1)
			n := new(big.Int).Sub(sm2K.pubKey.Params().N, one)
			k.Mod(k, n)
			k.Add(k, one)

			// Compute temporary public key
			tempX, tempY := sm2K.pubKey.ScalarBaseMult(k.Bytes())
			tempSK.X, tempSK.Y = tempSK.Add(
				sm2K.pubKey.X, sm2K.pubKey.Y,
				tempX, tempY,
			)

			// Verify temporary public key is a valid point on the reference curve
			isOn := tempSK.Curve.IsOnCurve(tempSK.X, tempSK.Y)
			if !isOn {
				return nil, errors.New("Failed temporary public key IsOnCurve check.")
			}

			return &sm2PublicKey{tempSK}, nil
		default:
			return nil, fmt.Errorf("Unsupported 'KeyDerivOpts' provided [%v]", opts)
		}
	default:
		return csp.sw.KeyDeriv(k, opts)
	}
}

// KeyImport imports a key from its raw representation using opts.
// The opts argument should be appropriate for the primitive used.
func (csp *impl) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	switch opts.(type) {
	case *bccsp.SM2PKIXPublicKeyImportOpts:
		der, ok := raw.([]byte)
		if !ok {
			return nil, errors.New("Invalid raw material. Expected byte array.")
		}

		if len(der) == 0 {
			return nil, errors.New("Invalid raw. It must not be nil.")
		}

		lowLevelKey, err := DERToPublicKey(der)
		if err != nil {
			return nil, fmt.Errorf("Failed converting PKIX to ECDSA public key [%s]", err)
		}

		sm2PK, ok := lowLevelKey.(*sm2.PublicKey)
		if !ok {
			return nil, errors.New("Failed casting to ECDSA public key. Invalid raw material.")
		}

		return &sm2PublicKey{sm2PK}, nil
	case *bccsp.SM2PrivateKeyImportOpts:
		der, ok := raw.([]byte)
		if !ok {
			return nil, errors.New("[SM2PrivateKeyImportOpts] Invalid raw material. Expected byte array.")
		}

		if len(der) == 0 {
			return nil, errors.New("[SM2PrivateKeyImportOpts] Invalid raw. It must not be nil.")
		}

		lowLevelKey, err := DERToPrivateKey(der)
		if err != nil {
			return nil, fmt.Errorf("Failed converting PKIX to SM2 private key [%s]", err)
		}

		sm2SK, ok := lowLevelKey.(*sm2.PrivateKey)
		if !ok {
			return nil, errors.New("Failed casting to SM2 private key. Invalid raw material.")
		}

		return &sm2PrivateKey{sm2SK}, nil
	case *bccsp.SM2GoPublicKeyImportOpts:
		lowLevelKey, ok := raw.(*sm2.PublicKey)
		if !ok {
			return nil, errors.New("Invalid raw material. Expected *sm2.PublicKey.")
		}

		return &sm2PublicKey{lowLevelKey}, nil
	case *bccsp.X509PublicKeyImportOpts:
		x509Cert, ok := raw.(*x509.Certificate)
		if !ok {
			return nil, errors.New("Invalid raw material. Expected *x509.Certificate.")
		}

		//if pk, ok := x509Cert.PublicKey.(*sm2.PublicKey); ok {
		//	return &sm2PublicKey{pk}, nil
		//} else {
		//	return csp.sw.KeyImport(raw, opts)
		//}
		pk := x509Cert.PublicKey

		switch pk.(type) {
		case *sm2.PublicKey:
			return &sm2PublicKey{pk.(*sm2.PublicKey)}, nil
		default:
			//convert to origin x509
			origCert, err := origx509.ParseCertificate(raw.(*x509.Certificate).Raw)
			if err != nil {
				return nil, errors.New("Invalid raw material. can't do x509 converting from flyinox x509 to origin x509.")
			}
			return csp.sw.KeyImport(origCert, opts)
		}
	default:
		return csp.sw.KeyImport(raw, opts)
	}
}

// GetKey returns the key this CSP associates to
// the Subject Key Identifier ski.
func (csp *impl) GetKey(ski []byte) (k bccsp.Key, err error) {
	if k, err := csp.ks.GetKey(ski); err == nil {
		return k, err
	} else {
		return csp.sw.GetKey(ski)
	}
}

// Hash hashes messages msg using options opts.
// If opts is nil, the default hash function will be used.
func (csp *impl) Hash(msg []byte, opts bccsp.HashOpts) (hash []byte, err error) {
	switch opts.(type) {
	case *bccsp.SM3Opts:
		h := sm3.New()
		h.Write(msg)
		return h.Sum(nil), nil
	case *bccsp.SM3SIGOpts:
		h := NewSM3Sig()
		h.Write(msg)
		return h.Sum(nil), nil
	default:
		return csp.sw.Hash(msg, opts)
	}
}

// GetHash returns and instance of hash.Hash using options opts.
// If opts is nil, the default hash function will be returned.
func (csp *impl) GetHash(opts bccsp.HashOpts) (h hash.Hash, err error) {
	switch opts.(type) {
	case *bccsp.SM3Opts:
		return sm3.New(), nil
	case *bccsp.SM3SIGOpts:
		return
	default:
		return csp.sw.GetHash(opts)
	}
}

// Sign signs digest using key k.
// The opts argument should be appropriate for the algorithm used.
//
// Note that when a signature of a hash of a larger message is needed,
// the caller is responsible for hashing the larger message and passing
// the hash (as digest).
func (csp *impl) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	switch k := k.(type) {
	case *sm2PrivateKey:
		return signSM2(k.privKey, digest, opts)
	default:
		return csp.sw.Sign(k, digest, opts)
	}
}

// Verify verifies signature against key k and digest
// The opts argument should be appropriate for the algorithm used.
func (csp *impl) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
	switch k := k.(type) {
	case *sm2PrivateKey:
		return verifySM2(&(k.privKey.PublicKey), signature, digest, opts)
	case *sm2PublicKey:
		return verifySM2(k.pubKey, signature, digest, opts)
	default:
		return csp.sw.Verify(k, signature, digest, opts)
	}
}

// Encrypt encrypts plaintext using key k.
// The opts argument should be appropriate for the algorithm used.
func (csp *impl) Encrypt(k bccsp.Key, plaintext []byte, opts bccsp.EncrypterOpts) (ciphertext []byte, err error) {
	return nil, nil
}

// Decrypt decrypts ciphertext using key k.
// The opts argument should be appropriate for the algorithm used.
func (csp *impl) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {
	return nil, nil
}
