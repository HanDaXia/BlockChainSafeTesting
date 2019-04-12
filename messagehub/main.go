package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/ledgercheck"
	"net/http"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/utils"
	"os"
)

var (
	localUrl string
	serverUrl string
	verifyUrl string
	randomUrl string
	regUrl string
)

func init() {
	localUrl = os.Getenv("LOCAL_URL")
	if len(localUrl) == 0 {
		panic("environment LOCAL_URL not found")
	}
	localUrl += "/Check"

	serverUrl = os.Getenv("SERVER_URL")
	if len(serverUrl) == 0 {
		panic("environment SERVER_URL not found")
	}
	verifyUrl = serverUrl + "/VerifySignature"
	randomUrl = serverUrl + "/RandomNumberCheck"

	regUrl = os.Getenv("REGISTER_URL")
	if len(regUrl) == 0 {
		panic("environment REGISTER_URL not found")
	}
	regUrl += "/RegistServer"

	regReq := ledgercheck.RegisterReq{ServerType:0, ServerAddress:localUrl}
	regBytes, err := json.Marshal(&regReq)
	if err != nil {
		fmt.Println("json.marshal error : ", err)
		panic(err)
	}
	_, err = utils.PostBytes(regUrl, regBytes)
	if err != nil {
		fmt.Println("register server error : ", err)
		panic(err)
	}
	fmt.Println("register result ok")
}

func FabricLedgerCheck(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		fmt.Println("body is nil")
		_, err := w.Write([]byte([]byte("{}")))
		if err != nil {
			fmt.Printf("write response error : %s", err)
		}
		return
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read data error : %s", err)
		w.WriteHeader(500)
		return
	}

	checkReq := &ledgercheck.CheckRequest{}
	err = json.Unmarshal(data, checkReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	var fbRet ledgercheck.FabricResp
	sd := &ledgercheck.SourceData{}
	err = json.Unmarshal(checkReq.Data, sd)
	if err != nil || (sd.LedgerData==nil && sd.RandomData==nil){
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	ckResp := ledgercheck.CheckResponse{}
	if len(sd.LedgerData) > 0 {
		var ledgerRes *ledgercheck.LedgerData
		ledgerRes, err = ledgercheck.AnalyzeChannelLedger(sd.LedgerData)
		if err != nil {
			fmt.Println("analyze ledger error : ", err)
			ledgerRes = &ledgercheck.LedgerData{Err: err.Error()}
		}
		fbRet.LedgerDetail, _ = json.Marshal(ledgerRes)
	}
	if len(sd.RandomData) > 0 {
		randomReq := ledgercheck.RandomRequest{RandomCheckType:1}
		randomReq.RandomData = sd.RandomData
		reqBytes, _ := json.Marshal(randomReq)
		randomResult, err := utils.PostBytes(randomUrl, reqBytes)
		fmt.Println(randomResult)
		if err != nil {
			fbRet.RandomDetail = []byte(err.Error())
		} else {
			fbRet.RandomDetail = randomResult
		}
	}

	ckResp.Result.OtherResult, err = json.Marshal(fbRet)
	retBytes, err := json.Marshal(&ckResp)
	w.Write(retBytes)
	w.WriteHeader(200)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Check", FabricLedgerCheck).Methods("POST")
	http.Handle("/", r)
	//ports := strings.Split()
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe error : ", err)
	}
}
