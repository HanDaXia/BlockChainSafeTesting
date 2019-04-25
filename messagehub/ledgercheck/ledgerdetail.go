package ledgercheck

import (
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil/cert"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil/x509"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/ledgercheck/blockfilemgr"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/utils"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/capabilities"
	"github.com/hyperledger/fabric/common/channelconfig"
	ledgerutil "github.com/hyperledger/fabric/common/ledger/util"
	"github.com/hyperledger/fabric/common/util"
	fabricmsp "github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/common"
	m "github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	putil "github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var (
	LocalUrl string
	ServerUrl string
	VerifyUrl string
	RandomUrl string
	RegUrl string
)

type ChannelInfo struct {
	Protos *channelconfig.ChannelProtos
	MspManager fabricmsp.MSPManager
	AppConfig         *ApplicationConfig
	OrdererConfig     *OrdererConfig
	ConsortiumsConfig *ConsortiumsConfig
}

// OrdererConfig holds the orderer configuration information
type OrdererConfig struct {
	Protos *channelconfig.OrdererProtos
	Orgs   map[string]*OrganizationConfig
	BatchTimeout time.Duration
}

// NewOrdererConfig creates a new instance of the orderer config
func NewOrdererConfig(ordererGroup *common.ConfigGroup, mspConfig *channelconfig.MSPConfigHandler) (*OrdererConfig, error) {
	oc := &OrdererConfig{
		Protos: &channelconfig.OrdererProtos{},
		Orgs:   make(map[string]*OrganizationConfig),
	}

	if err := channelconfig.DeserializeProtoValuesFromGroup(ordererGroup, oc.Protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	for orgName, orgGroup := range ordererGroup.Groups {
		var err error
		if oc.Orgs[orgName], err = NewOrganizationConfig(orgName, orgGroup, mspConfig); err != nil {
			return nil, err
		}
	}
	return oc, nil
}

// OrganizationConfig stores the configuration for an organization
type OrganizationConfig struct {
	Protos *channelconfig.OrganizationProtos

	MspConfigHandler *channelconfig.MSPConfigHandler
	Msp              fabricmsp.MSP
	MspID            string
	Name             string
}

// NewOrganizationConfig creates a new config for an organization
func NewOrganizationConfig(name string, orgGroup *common.ConfigGroup, mspConfigHandler *channelconfig.MSPConfigHandler) (*OrganizationConfig, error) {
	if len(orgGroup.Groups) > 0 {
		return nil, fmt.Errorf("organizations do not support sub-groups")
	}

	oc := &OrganizationConfig{
		Protos:           &channelconfig.OrganizationProtos{},
		Name:             name,
		MspConfigHandler: mspConfigHandler,
	}

	if err := channelconfig.DeserializeProtoValuesFromGroup(orgGroup, oc.Protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}
	return oc, nil
}

// ApplicationConfig implements the Application interface
type ApplicationConfig struct {
	ApplicationOrgs map[string]*ApplicationOrgConfig
	Protos          *channelconfig.ApplicationProtos
}

// NewApplicationConfig creates config from an Application config group
func NewApplicationConfig(appGroup *common.ConfigGroup, mspConfig *channelconfig.MSPConfigHandler) (*ApplicationConfig, error) {
	ac := &ApplicationConfig{
		ApplicationOrgs: make(map[string]*ApplicationOrgConfig),
		Protos:          &channelconfig.ApplicationProtos{},
	}

	if err := channelconfig.DeserializeProtoValuesFromGroup(appGroup, ac.Protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	var err error
	for orgName, orgGroup := range appGroup.Groups {
		ac.ApplicationOrgs[orgName], err = NewApplicationOrgConfig(orgName, orgGroup, mspConfig)
		if err != nil {
			return nil, err
		}
	}
	return ac, nil
}

// ApplicationOrgConfig defines the configuration for an application org
type ApplicationOrgConfig struct {
	*OrganizationConfig
	protos *channelconfig.ApplicationOrgProtos
	name   string
}

// NewApplicationOrgConfig creates a new config for an application org
func NewApplicationOrgConfig(id string, orgGroup *common.ConfigGroup, mspConfig *channelconfig.MSPConfigHandler) (*ApplicationOrgConfig, error) {
	if len(orgGroup.Groups) > 0 {
		return nil, fmt.Errorf("ApplicationOrg config does not allow sub-groups")
	}

	protos := &channelconfig.ApplicationOrgProtos{}
	orgProtos := &channelconfig.OrganizationProtos{}

	if err := channelconfig.DeserializeProtoValuesFromGroup(orgGroup, protos, orgProtos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	aoc := &ApplicationOrgConfig{
		name:   id,
		protos: protos,
		OrganizationConfig: &OrganizationConfig{
			Name:             id,
			Protos:           orgProtos,
			MspConfigHandler: mspConfig,
		},
	}
	return aoc, nil
}

// ConsortiumsConfig holds the consoritums configuration information
type ConsortiumsConfig struct {
	Consortiums map[string]*ConsortiumConfig
}

// NewConsortiumsConfig creates a new instance of the consoritums config
func NewConsortiumsConfig(consortiumsGroup *common.ConfigGroup, mspConfig *channelconfig.MSPConfigHandler) (*ConsortiumsConfig, error) {
	cc := &ConsortiumsConfig{
		Consortiums: make(map[string]*ConsortiumConfig),
	}

	for consortiumName, consortiumGroup := range consortiumsGroup.Groups {
		var err error
		if cc.Consortiums[consortiumName], err = NewConsortiumConfig(consortiumGroup, mspConfig); err != nil {
			return nil, err
		}
	}
	return cc, nil
}

// ConsortiumConfig holds the consoritums configuration information
type ConsortiumConfig struct {
	Protos *channelconfig.ConsortiumProtos
	Orgs   map[string]*OrganizationConfig
}

// NewConsortiumConfig creates a new instance of the consoritums config
func NewConsortiumConfig(consortiumGroup *common.ConfigGroup, mspConfig *channelconfig.MSPConfigHandler) (*ConsortiumConfig, error) {
	cc := &ConsortiumConfig{
		Protos: &channelconfig.ConsortiumProtos{},
		Orgs:   make(map[string]*OrganizationConfig),
	}

	if err := channelconfig.DeserializeProtoValuesFromGroup(consortiumGroup, cc.Protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	for orgName, orgGroup := range consortiumGroup.Groups {
		var err error
		if cc.Orgs[orgName], err = NewOrganizationConfig(orgName, orgGroup, mspConfig); err != nil {
			return nil, err
		}
	}

	return cc, nil
}

func getChannelGroupDetail(channelGroup *common.ConfigGroup) (*ChannelInfo, error){
	var ci ChannelInfo
	ci.Protos = &channelconfig.ChannelProtos{}
	if err := channelconfig.DeserializeProtoValuesFromGroup(channelGroup, ci.Protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	caps := capabilities.NewChannelProvider(ci.Protos.Capabilities.Capabilities)
	mspConfigHandler := channelconfig.NewMSPConfigHandler(caps.MSPVersion())
	var err error
	for groupName, group := range channelGroup.Groups {
		switch groupName {
		case channelconfig.ApplicationGroupKey:
			ci.AppConfig, err = NewApplicationConfig(group, mspConfigHandler)
		case channelconfig.OrdererGroupKey:
			ci.OrdererConfig, err = NewOrdererConfig(group, mspConfigHandler)
		case channelconfig.ConsortiumsGroupKey:
			ci.ConsortiumsConfig, err = NewConsortiumsConfig(group, mspConfigHandler)
		default:
			return nil, fmt.Errorf("Disallowed channel group: %s", group)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "could not create channel %s sub-group config", groupName)
		}
	}

	if ci.MspManager, err = mspConfigHandler.CreateMSPManager(); err != nil {
		return nil, err
	}
	return &ci, nil
}

func CheckDataSignature(crt *x509.Certificate, rawBytes, signature []byte) (vi VerifyInfo, err error) {
	checkSigErrr := crt.CheckSignature(crt.SignatureAlgorithm, rawBytes, signature)
	if checkSigErrr != nil {
		fmt.Println("check signature error:", checkSigErrr)
	}

	hashAndSig := strings.Split(crt.SignatureAlgorithm.String(), "-")
	if len(hashAndSig) != 2 {
		err = errors.New("error cert signaturealgorithm")
		return
	}

	vi.HashAlgo = hashAndSig[1]
	vi.SignAlgo = hashAndSig[0]
	vi.PubAlgo, vi.PubCurve = cert.GetPublicKeyAlgo(crt)
	//err = crt.CheckSignature(crt.SignatureAlgorithm, rawBytes, signature)

	tbsCert, err := crt.GetTBS()
	if err != nil {
		return
	}
	pkBytes, _ := json.Marshal(crt.PublicKey)

	paramsData := tbsCert.PublicKey.Algorithm.Parameters.FullBytes
	namedCurveOID := new(asn1.ObjectIdentifier)
	_, err = asn1.Unmarshal(paramsData, namedCurveOID)
	vs := VerifySignature{
		PublicKey: pkBytes,
		Signature: signature,
		PlainText: rawBytes,
		SignatureAlgo: tbsCert.SignatureAlgorithm.Algorithm,
		PublicAlgo: tbsCert.PublicKey.Algorithm.Algorithm,
		Curve: *namedCurveOID,
	}

	sendData, _:= json.Marshal(vs)
	resp, err := utils.PostBytes(VerifyUrl, sendData)
	if err != nil {
		vi.VerifyOk = false
		return
	}

	verifyResp := VerifyRespFromServer{}
	err = json.Unmarshal(resp, &verifyResp)
	if err != nil {
		vi.VerifyOk = false
		return
	}

	if verifyResp.VerifySuccess {
		vi.VerifyOk = true
	} else {
		vi.VerifyOk = false
	}

	return
}

func analyzeData(bd *BlockData, envelope *common.Envelope) error {
	payload, err := putil.UnmarshalPayload(envelope.Payload)
	if err != nil {
		return err
	}
	chdr, err := putil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return err
	}

	// verify tx signature
	var txData TxData
	txData.TxID = chdr.TxId
	shdr, err := putil.GetSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return err
	}
	sId := &m.SerializedIdentity{}
	err = proto.Unmarshal(shdr.Creator, sId)
	if err != nil {
		return err
	}
	txCert, err := x509.ParsePem(sId.IdBytes)
	if err != nil {
		return err
	}
	txData.VerifyInfo, err = CheckDataSignature(txCert, envelope.Payload, envelope.Signature)
	if err != nil {
		return err
	}

	// verify endorsement signature
	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION{
		// extract actions from the envelope message
		tx, err := putil.GetTransaction(payload.Data)
		if err != nil {
			return err
		}
		if len(tx.Actions) == 0 {
			return errors.New("no action in tx")
		}
		actionPayload, _, err := putil.GetPayloads(tx.Actions[0])

		// check all endorsements in tx
		for _, endorseMent := range actionPayload.Action.Endorsements {
			sId := &m.SerializedIdentity{}
			err := proto.Unmarshal(endorseMent.Endorser, sId)
			if err != nil {
				return err
			}
			endorseCert, err := x509.ParsePem(sId.IdBytes)
			if err != nil {
				return err
			}
			endorse := Endorser{}
			endorse.Name = sId.Mspid
			endorse.VerifyInfo, err = CheckDataSignature(endorseCert,
				append(actionPayload.Action.ProposalResponsePayload, endorseMent.Endorser...),
				endorseMent.Signature)
			if err != nil {
				return err
			}
			txData.Endorsers = append(txData.Endorsers, endorse)
		}
	}
	bd.Txs = append(bd.Txs, txData)
	return nil
}

func AnalyzeChannelLedger(ledgerBytes []byte) (*LedgerData, error) {
	blf, err := blockfilemgr.NewBlockStream(ledgerBytes)
	if err != nil {
		return nil, err
	}

	var ld *LedgerData
	tmpCounter := 0
	for {
		tmpCounter ++
		rawBytes, err := blf.NextBlockBytes()
		if err != nil || rawBytes == nil{
			return ld, nil
		}

		rawBlock, err := deserializeBlock(rawBytes)
		if err != nil {
			return nil, err
		}
		if tmpCounter == 1 {
			ld, err = GetLedgerInfo(rawBlock)
			if err != nil {
				return nil, err
			}
			continue
		}

		if len(ld.Blocks) == 6{
			return ld, nil
		}

		if len(rawBlock.Metadata.Metadata) < 1 {
			return nil, errors.New("metadata error")
		}

		var bd BlockData
		bd.Height = rawBlock.Header.Number
		metadata, err := putil.GetMetadataFromBlock(rawBlock, common.BlockMetadataIndex_SIGNATURES)
		if err != nil {
			return nil, err
		}
		if metadata.Signatures != nil {
			shdr, err := putil.GetSignatureHeader(metadata.Signatures[0].SignatureHeader)
			if err != nil {
				return nil, err
			}

			sId := &m.SerializedIdentity{}
			err = proto.Unmarshal(shdr.Creator, sId)
			if err != nil {
				return nil, err
			}

			blockCert, err := x509.ParsePem(sId.IdBytes)
			if err != nil {
				return nil, errors.Wrap(err, "parseCertificate failed")
			}

			payloadBytes := util.ConcatenateBytes(metadata.Value, metadata.Signatures[0].SignatureHeader, rawBlock.Header.Bytes())
			bd.VerifyInfo, err = CheckDataSignature(blockCert, payloadBytes, metadata.Signatures[0].Signature)
			if err != nil {
				return nil, err
			}
		}

		for i := 0; i < len(rawBlock.Data.Data); i++ {
			envelope, err := putil.ExtractEnvelope(rawBlock, i)
			if err != nil {
				return nil, err
			}
			err = analyzeData(&bd, envelope)
			if err != nil {
				return nil, err
			}
		}
		ld.Blocks = append(ld.Blocks, bd)
	}
}

func GetLedgerInfo(rawBlock *common.Block) (*LedgerData, error) {
	ld := &LedgerData{}
	envelope, err := putil.ExtractEnvelope(rawBlock, 0)
	if err != nil {
		return nil, err
	}
	payload, err := putil.UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, err
	}
	chdr, err := putil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, err
	}

	if common.HeaderType(chdr.Type) == common.HeaderType_CONFIG {
		configEnv := &common.ConfigEnvelope{}
		_, err = putil.UnmarshalEnvelopeOfType(envelope, common.HeaderType_CONFIG, configEnv)
		if err != nil {
			return nil, err
		}
		if configEnv.Config == nil || configEnv.Config.ChannelGroup == nil {
			return nil, errors.New("block struct error")
		}

		cgd, err := getChannelGroupDetail(configEnv.Config.ChannelGroup)
		if err != nil {
			return nil, err
		}
		ld.ChannelID = chdr.ChannelId
		ld.TxCount = cgd.OrdererConfig.Protos.BatchSize.MaxMessageCount
		ld.BlockSize = cgd.OrdererConfig.Protos.BatchSize.AbsoluteMaxBytes
		ld.BlockTime = cgd.OrdererConfig.Protos.BatchTimeout.Timeout
		ld.ConsusType = cgd.OrdererConfig.Protos.ConsensusType.Type
		ld.TxMaxSize = cgd.OrdererConfig.Protos.BatchSize.PreferredMaxBytes
		for _, addr := range cgd.Protos.OrdererAddresses.Addresses {
			ld.OrdererAddr += addr + "\n"
		}
		return ld, nil
	}
	return nil, errors.New("block format error")
}

type OrgAndAnchors struct {
	Name string
	Anchors []*peer.AnchorPeer
}

type locPointer struct {
	offset      int
	bytesLength int
}

//The order of the transactions must be maintained for history
type txindexInfo struct {
	txID string
	loc  *locPointer
}
func extractHeader(buf *ledgerutil.Buffer) (*common.BlockHeader, error) {
	header := &common.BlockHeader{}
	var err error
	if header.Number, err = buf.DecodeVarint(); err != nil {
		return nil, err
	}
	if header.DataHash, err = buf.DecodeRawBytes(false); err != nil {
		return nil, err
	}
	if header.PreviousHash, err = buf.DecodeRawBytes(false); err != nil {
		return nil, err
	}
	if len(header.PreviousHash) == 0 {
		header.PreviousHash = nil
	}
	return header, nil
}

func extractData(buf *ledgerutil.Buffer) (*common.BlockData, []*txindexInfo, error) {
	data := &common.BlockData{}
	var txOffsets []*txindexInfo
	var numItems uint64
	var err error

	if numItems, err = buf.DecodeVarint(); err != nil {
		return nil, nil, err
	}
	for i := uint64(0); i < numItems; i++ {
		var txEnvBytes []byte
		var txid string
		txOffset := buf.GetBytesConsumed()
		if txEnvBytes, err = buf.DecodeRawBytes(false); err != nil {
			return nil, nil, err
		}
		if txid, err = extractTxID(txEnvBytes); err != nil {
			return nil, nil, err
		}
		data.Data = append(data.Data, txEnvBytes)
		idxInfo := &txindexInfo{txid, &locPointer{txOffset, buf.GetBytesConsumed() - txOffset}}
		txOffsets = append(txOffsets, idxInfo)
	}
	return data, txOffsets, nil
}

func extractMetadata(buf *ledgerutil.Buffer) (*common.BlockMetadata, error) {
	metadata := &common.BlockMetadata{}
	var numItems uint64
	var metadataEntry []byte
	var err error
	if numItems, err = buf.DecodeVarint(); err != nil {
		return nil, err
	}
	for i := uint64(0); i < numItems; i++ {
		if metadataEntry, err = buf.DecodeRawBytes(false); err != nil {
			return nil, err
		}
		metadata.Metadata = append(metadata.Metadata, metadataEntry)
	}
	return metadata, nil
}

func extractTxID(txEnvelopBytes []byte) (string, error) {
	txEnvelope, err := putil.GetEnvelopeFromBlock(txEnvelopBytes)
	if err != nil {
		return "", err
	}
	txPayload, err := putil.GetPayload(txEnvelope)
	if err != nil {
		return "", nil
	}
	chdr, err := putil.UnmarshalChannelHeader(txPayload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	return chdr.TxId, nil
}

func deserializeBlock(serializedBlockBytes []byte) (*common.Block, error) {
	block := &common.Block{}
	var err error
	b := ledgerutil.NewBuffer(serializedBlockBytes)
	if block.Header, err = extractHeader(b); err != nil {
		return nil, err
	}
	if block.Data, _, err = extractData(b); err != nil {
		return nil, err
	}
	if block.Metadata, err = extractMetadata(b); err != nil {
		return nil, err
	}
	return block, nil
}
