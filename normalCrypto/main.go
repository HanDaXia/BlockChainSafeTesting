package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"normalCrypto/randomCheck"
	"normalCrypto/signatureCheck"
)

type VerifySignatureResult struct {
	VerifySuccess bool
}

type RandomCheckResult struct {
	result string
}

func main() {
	startServer()
}

func startServer() {
	http.HandleFunc("/VerifySignature", VerifySignature)
	http.HandleFunc("/RandomNumberCheck", RandomNumberCheck)

	err := http.ListenAndServe("0.0.0.0:8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func VerifySignature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var signInfo signatureCheck.SignatureInfo

		if err := json.Unmarshal(body, &signInfo); err == nil {
			result, err := signatureCheck.VerifySignature(signInfo)
			if err != nil {
				errMsg, _ := json.Marshal(err.Error())
				ResponseWithOrigin(w, r, http.StatusInternalServerError, errMsg)
				return
			}
			ret := VerifySignatureResult{result}
			resp, err := json.Marshal(ret)
			fmt.Printf("result = %s", string(resp))
			if err != nil {
				errMsg, _ := json.Marshal("Marshal result failed")
				ResponseWithOrigin(w, r, http.StatusInternalServerError, errMsg)
				return
			}
			ResponseWithOrigin(w, r, http.StatusOK, resp)
		} else {
			errMsg, _ := json.Marshal("Unmarshal request failed")
			ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
		}
	} else {
		errMsg, _ := json.Marshal("Just surport Post request")
		ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
	}

}

func RandomNumberCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var randomCheckInfo randomCheck.RandomCheckInfo

		if err := json.Unmarshal(body, &randomCheckInfo); err == nil {
			if len(randomCheckInfo.RandomData) < randomCheck.MINRAMDOMDATA{
				errMsg, _ := json.Marshal(fmt.Sprintf("RandomData must larger than %d", randomCheck.MINRAMDOMDATA))
				ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
				return
			}
			result := randomCheck.DealRandomCheck(randomCheckInfo.RandomCheckType, randomCheckInfo.RandomData)
			ret := RandomCheckResult{result}
			resp, err := json.Marshal(ret)
			if err != nil {
				errMsg, _ := json.Marshal("Marshal result failed")
				ResponseWithOrigin(w, r, http.StatusInternalServerError, errMsg)
				return
			}
			ResponseWithOrigin(w, r, http.StatusOK, resp)
		} else {
			errMsg, _ := json.Marshal("Unmarshal request failed")
			ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
		}
	} else {
		errMsg, _ := json.Marshal("Just surport Post request")
		ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
	}

}

func ResponseWithOrigin(w http.ResponseWriter, r *http.Request, code int, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}
