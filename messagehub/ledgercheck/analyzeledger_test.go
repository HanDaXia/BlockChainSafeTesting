package ledgercheck

import (
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/protos/common"
	m "github.com/hyperledger/fabric/protos/msp"
	putil "github.com/hyperledger/fabric/protos/utils"
	"io/ioutil"
	"ledgercheck/blockfilemgr"
	"net/http"
	"testing"
)


func TestAnalizeLedger(t *testing.T) {
	vs := VerifySignature{PublicKey: []byte("aaaa"),
		Signature: []byte("aaaa"),
		PlainText: []byte("aaaa"),
		SignatureAlgo: asn1.ObjectIdentifier{1,1},
		PublicAlgo: asn1.ObjectIdentifier{1,1},
		Curve: asn1.ObjectIdentifier{1,1},}

	bytes, err := json.Marshal(vs)
	fmt.Println(string(bytes))
	fmt.Println(bytes, err)

	var res VerifySignature
	err = json.Unmarshal(bytes, &res)
	fmt.Println(err)
}

func PostBytes(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil{
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Keep-Alive")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func TestAnalyzeLedger(t *testing.T) {
	ledgerBytes := []byte("ledgerdata")
	blf, err := blockfilemgr.NewBlockStream(ledgerBytes)
	if err != nil {
		return
	}

	tmpCounter := 0
	for {
		tmpCounter ++
		rawBytes, err := blf.NextBlockBytes()
		if err != nil {
			return
		}
		rawBlock, err := deserializeBlock(rawBytes)
		if err != nil {
			return
		}
		envelope, err := putil.ExtractEnvelope(rawBlock, 0)
		if err != nil {
			return
		}
		payload, err := putil.UnmarshalPayload(envelope.Payload)
		if err != nil {
			return
		}
		chdr, err := putil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return
		}
		shdr, err := putil.GetSignatureHeader(payload.Header.SignatureHeader)
		if err != nil {
			return
		}

		sId := &m.SerializedIdentity{}
		err = proto.Unmarshal(shdr.Creator, sId)
		if err != nil {
			return
		}

		cert1, err := cryptoutil.ParsePem(sId.IdBytes)
		if err != nil {
			return
		}

		_, err = cert1.GetTBS()
		if err != nil {
			return
		}
		//pkBytes, _ := json.Marshal(cert1.PublicKey)
		//
		//paramsData := tbsCert.PublicKey.Algorithm.Parameters.FullBytes
		//namedCurveOID := new(asn1.ObjectIdentifier)
		//_, err = asn1.Unmarshal(paramsData, namedCurveOID)
		//vs := VerifySignature{
		//	PublicKey: pkBytes,
		//	Signature: envelope.Signature,
		//	PlainText: envelope.Payload,
		//	SignatureAlgo: tbsCert.SignatureAlgorithm.Algorithm,
		//	PublicAlgo: tbsCert.PublicKey.Algorithm.Algorithm,
		//	Curve: *namedCurveOID,
		//}
		//
		//sendData, _ := json.Marshal(vs)
		//resp, err := PostBytes("http://172.16.0.250:8081/VerifySignature", sendData)
		//
		//fmt.Println(string(resp), err)

		err = cert1.CheckSignature(cert1.SignatureAlgorithm, envelope.Payload, envelope.Signature)
		fmt.Printf("check signature result : %s, creator : %s\n", err, cert1.Subject.CommonName)

		if common.HeaderType(chdr.Type) == common.HeaderType_CONFIG{
			configEnv := &common.ConfigEnvelope{}
			_, err = putil.UnmarshalEnvelopeOfType(envelope, common.HeaderType_CONFIG, configEnv)
			if err != nil {
				return
			}
			if configEnv.Config == nil || configEnv.Config.ChannelGroup == nil {
				return
			}

			cgd, err := getChannelGroupDetail(configEnv.Config.ChannelGroup)
			if err != nil {
				return
			}
			fmt.Println(cgd)

		} else if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
			// extract actions from the envelope message
			tx, err := putil.GetTransaction(payload.Data)
			if err != nil {
				fmt.Println("GetTransaction error: ", err)
				return
			}
			if len(tx.Actions) == 0 {
				return
			}

			actionPayload, respPayload, err := putil.GetPayloads(tx.Actions[0])
			//signPayload := actionPayload.ChaincodeProposalPayload
			for _, endorseMent := range actionPayload.Action.Endorsements {
				sId := &m.SerializedIdentity{}
				err := proto.Unmarshal(endorseMent.Endorser, sId)
				if err != nil {
					return
				}
				txCert, err := cryptoutil.ParsePem(sId.IdBytes)
				if err != nil {
					return
				}
				err = txCert.CheckSignature(txCert.SignatureAlgorithm,
					append(actionPayload.Action.ProposalResponsePayload, endorseMent.Endorser...),
					endorseMent.Signature)
				fmt.Println(err)
			}

			txRWSet := &rwsetutil.TxRwSet{}
			if err = txRWSet.FromProtoBytes(respPayload.Results); err != nil {
				fmt.Println("Failed obtaining TxRwSet from ChaincodeAction's results", err)
				continue
			}
			fmt.Println(txRWSet.NsRwSets)
		}
	}
}
