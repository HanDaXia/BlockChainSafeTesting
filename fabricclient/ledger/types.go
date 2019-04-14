package ledger

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

type FabricCheckResult struct {
	OtherResult []byte
	SignatureResult VerifyResponse
}

type RandomResponse struct {
	Result string
}

type CheckResponse struct {
	Status int
	Result FabricCheckResult
}

type DistributeResult struct {
	Result []byte
}
