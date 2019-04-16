package signatureCheck

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/json"
	"github.com/flyinox/crypto/sm/sm2"
	"github.com/flyinox/crypto/x509"
	"github.com/pkg/errors"
)

type SignatureInfo struct {
	PublicKey     []byte
	Signature     []byte
	PlainText     []byte
	SignatureAlgo asn1.ObjectIdentifier
	PublicAlgo    asn1.ObjectIdentifier
	Curve         asn1.ObjectIdentifier
}

func VerifySignature(signInfo SignatureInfo) (bool, error) {

	signatureAlgo := x509.GetSignatureAlgorithmFromOID(signInfo.SignatureAlgo)

    pub, err := GetPublicKey(signatureAlgo, signInfo.PublicAlgo, signInfo.PublicKey, signInfo.Curve)
    if err != nil {
        return false, err
    }

	result := x509.CheckSignatureNor(signatureAlgo, signInfo.PlainText, signInfo.Signature, pub)

	if result != nil {
		return false, nil
	}

	return true, nil
}

func GetPublicKey(signAlgo x509.SignatureAlgorithm, pubOID asn1.ObjectIdentifier, pubBytes []byte, curve asn1.ObjectIdentifier) (pubKey interface{}, err error) {
	algo := x509.GetPublicKeyAlgorithmFromOID(pubOID)
	if algo == x509.ECDSA {
	    if signAlgo == x509.SM2WithSM3 || signAlgo == x509.SM2WithSHA1 || signAlgo == x509.SM2WithSHA256 {
	        algo = x509.SM2
        }
    }
	switch algo {
	case x509.SM2:
		namedCurve := x509.NamedCurveFromOID(curve)
		if namedCurve == nil {
			return nil, errors.New("x509: unsupported SM2 elliptic curve")
		}
		Key := &sm2.PublicKey{}
		_ = json.Unmarshal(pubBytes, Key)
		pubKey = &sm2.PublicKey{namedCurve, Key.X, Key.Y}
		return
	case x509.ECDSA:
		namedCurve := x509.NamedCurveFromOID(curve)
		if namedCurve == nil {
			return nil, errors.New("x509: unsupported ECDSA elliptic curve")
		}
		Key := &ecdsa.PublicKey{}
		_ = json.Unmarshal(pubBytes, Key)
		pubKey = &ecdsa.PublicKey{namedCurve, Key.X, Key.Y}

		return
	case x509.DSA:
		pubKey = &dsa.PublicKey{}
		err = json.Unmarshal(pubBytes, pubKey)
		return
	case x509.RSA:
		pubKey = &rsa.PublicKey{}
		err = json.Unmarshal(pubBytes, pubKey)
		return
	default:
		err = errors.New("Unknown public key algorithm!")
	}

	return
}
