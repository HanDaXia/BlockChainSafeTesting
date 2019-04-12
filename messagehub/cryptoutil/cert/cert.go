package cert

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil/sm/sm2"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil/x509"
)

var (
	oidNamedCurveP224    = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256    = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384    = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521    = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
	oidNamedCurveP256SM2 = asn1.ObjectIdentifier{1, 2, 156, 10197, 1, 301}
)

func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case sm2.P256Sm2():
		return oidNamedCurveP256SM2, true
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	}
	return nil, false
}

func curveNameFromOid(oid asn1.ObjectIdentifier) string {
	switch {
	case oid.Equal(oidNamedCurveP256SM2):
		return "sm2.p256"
	case oid.Equal(oidNamedCurveP224):
		return "elliptic.P224"
	case oid.Equal(oidNamedCurveP256):
		return "elliptic.P256"
	case oid.Equal(oidNamedCurveP384):
		return "elliptic.P384"
	case oid.Equal(oidNamedCurveP521):
		return "elliptic.P521"
	}
	return "unknown"
}

type CertFormatInfo struct {
	SignatureOid string
	SignatureName string
	PkAlgo string
	Err error
}

func GetPemDetail(rawPem []byte) (pr CertFormatInfo){
	block, _ := pem.Decode(rawPem)
	if block == nil {
		fmt.Println("pem.decode error")
		pr.Err = errors.New("decode pem error")
	}
	pr = GetCertDetail(block.Bytes)
	return
}

func GetCertDetail(rawCert []byte) (pr CertFormatInfo){
	//pemPath := "G:/certs/crypto-config-sm2-syl/peerOrganizations/org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem"
	//pemPath := "G:/goproj/src/github.com/hyperledger/fabric-sdk-go-master/pkg/core/config/testdata/certs/client_sdk_go.pem"
	x509Cert, err := x509.ParseCertificate(rawCert)
	if err != nil {
		fmt.Println("x509.ParseCertificate error : ", err)
		pr.Err = errors.New("parse x509 field failed")
		return
	}

	tbsCert, err := x509Cert.GetTBS()
	if err != nil {
		fmt.Println("cert.GetSignatureAlgorithm error : ", err)
		pr.Err = errors.New("解析TBS证书失败")
		return
	}

	pr.SignatureOid = tbsCert.SignatureAlgorithm.Algorithm.String()
	pr.SignatureName = x509Cert.SignatureAlgorithm.String()

	switch x509Cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		pr.PkAlgo = "ecdsa"
		//oid, ok := oidFromNamedCurve(pk.Curve)
		//fmt.Println(ok, oid)
	case *sm2.PublicKey:
		pr.PkAlgo = "sm2"
		//oid, ok := oidFromNamedCurve(pk.Curve)
		//fmt.Println(ok, oid)
	case *rsa.PublicKey:
		pr.PkAlgo = "rsa"
	case *dsa.PublicKey:
		pr.PkAlgo = "dsa"
	default:
		pr.PkAlgo = "unknown"
	}
	return
}

func GetPublicKeyAlgo(cert *x509.Certificate) (algoName, curveName string){
	switch pk := cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		algoName = "ecdsa"
		oid, _ := oidFromNamedCurve(pk.Curve)
		curveName = curveNameFromOid(oid)
	case *sm2.PublicKey:
		algoName = "sm2"
		oid, _ := oidFromNamedCurve(pk.Curve)
		curveName = curveNameFromOid(oid)
	case *rsa.PublicKey:
		algoName = "rsa"
	case *dsa.PublicKey:
		algoName = "dsa"
	default:
		algoName = "unknown"
	}
	return
}
