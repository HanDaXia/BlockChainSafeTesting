package ledgercheck

import "encoding/asn1"

type LedgerData struct {
	ChannelID string
	ConsusType string
	OrdererAddr string
	BlockTime string
	TxCount uint32
	TxMaxSize uint32
	BlockSize uint32
	Blocks []BlockData
	RandomResult []byte
	Err string
}

type VerifyInfo struct {
	SignAlgo string
	HashAlgo string
	PubAlgo string
	PubCurve string
	VerifyOk bool
}

type Endorser struct {
	VerifyInfo
	Name string
}

type TxData struct {
	VerifyInfo
	TxID string
	Endorsers []Endorser
}

type BlockData struct {
	VerifyInfo
	Height uint64
	Txs []TxData
}

type SourceData struct{
	LedgerData []byte
	RandomData []byte
}

type CheckRequest struct {
	CompanyName string
	CompanyID string
	CheckType int
	Data []byte
}

type VerifyResponse struct {
	VerifySuccess bool
	PublicKeyAlgo string
	SignatureAlgo string

}

type FabricResp struct {
	LedgerDetail []byte
	RandomDetail []byte
}

type FabricCheckResult struct {
	OtherResult []byte
	SignatureResult VerifyResponse
}

type CheckResponse struct {
	Status int
	Result FabricCheckResult
}

type RandomRequest struct {
	RandomCheckType int
	RandomData []byte
}

type RegisterReq struct {
	ServerType int
	ServerAddress string
}

type RegisterResp struct {
	Status int
}

type VerifySignature struct {
	PublicKey []byte
	Signature []byte
	PlainText []byte
	SignatureAlgo asn1.ObjectIdentifier
	PublicAlgo asn1.ObjectIdentifier
	Curve asn1.ObjectIdentifier
}

type VerifyRespFromServer struct {
	VerifySuccess bool
}
