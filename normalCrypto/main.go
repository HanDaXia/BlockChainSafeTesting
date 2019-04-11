package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"normalCrypto/randomCheck"
	"normalCrypto/signatureCheck"
	"os"
)

type VerifySignatureResult struct {
	VerifySuccess bool
}

type RandomCheckResult struct {
	Result string
}

func TestRandom()  {
    var bytes []byte
    //for i := 0; i < 10000000; i++{
    //    bit := rand.Int() % 2
    //    bytes = append(bytes, byte(bit + 48))
    //}
    filePath := "/home/hanhu/Rand_Number_Assess/data/data.pi"
    f, err := os.Open(filePath)
    if err != nil {
    	fmt.Println(err.Error())
	}
    bytes, err = ioutil.ReadAll(f)
    if err != nil {
		fmt.Println(err.Error())
	}
    result, _ := randomCheck.DealRandomCheck(1, bytes)
    fmt.Println(result)
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
			result, err := randomCheck.DealRandomCheck(randomCheckInfo.RandomCheckType, randomCheckInfo.RandomData)
			fmt.Println("randomCheck result = " + result)
			if err != nil {
                errMsg, _ := json.Marshal(err.Error())
                ResponseWithOrigin(w, r, http.StatusBadRequest, errMsg)
                return
            }
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
