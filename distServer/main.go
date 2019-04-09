package main

import (
    "distServer/distributor"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

const (
    FabricCheck = iota
    EthCheck
    BtcCheck
)

type RegistRequest struct {
    ServerType int
    ServerAddress string
}


type CheckResponse struct {
    result []byte
}

func main()  {
    go func() {
        startRegistServer()
    }()

    startServer()
}

func startServer()  {
    http.HandleFunc("/Check", Check)

    err := http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func startRegistServer()  {
    http.HandleFunc("/RegistServer", RegistServer)

    err := http.ListenAndServe(":6000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func Check(w http.ResponseWriter, r *http.Request)  {
    body, _ := ioutil.ReadAll(r.Body)

    var request =  distributor.CheckRequest{}

    if err := json.Unmarshal(body, &request); err == nil {
        result, err := distributor.SendRequestToServer(request)
        if err != nil {
            ResponseWithOrigin(w, r, http.StatusInternalServerError, []byte(err.Error()))
            return
        }
        ret := CheckResponse{result}
        resp, err := json.Marshal(ret)
        if err != nil{
            ResponseWithOrigin(w, r, http.StatusInternalServerError, []byte(err.Error()))
            return
        }
        ResponseWithOrigin(w, r, http.StatusOK, resp)
    } else {
        fmt.Println(err)
    }
}

func RegistServer(w http.ResponseWriter, r *http.Request)  {
    body, _ := ioutil.ReadAll(r.Body)
    var request =  RegistRequest{}

    if err := json.Unmarshal(body, &request); err == nil {
        distributor.ServerUpdate(request.ServerType, request.ServerAddress)
        ResponseWithOrigin(w, r, http.StatusOK, nil)
    } else {
        ResponseWithOrigin(w, r, http.StatusBadRequest, []byte(err.Error()))
    }

}


func ResponseWithOrigin(w http.ResponseWriter, r *http.Request, code int, json []byte) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}