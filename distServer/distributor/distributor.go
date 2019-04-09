package distributor

import (
    "bytes"
    "crypto"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

type CheckRequest struct {
    CompanyName string
    CompanyID string
    CheckType int
    Data []byte
}

type serverRequest struct {
    Txid []byte
    CheckRequest
}

var serverList = make(map[int]string, 0)

func SendRequestToServer(request CheckRequest) ([]byte, error) {
    serverType := request.CheckType
    server, ok := serverList[serverType]
    if !ok {
        return nil, errors.New("Unknown server type!")
    }
    var hashType = crypto.SHA256
    h := hashType.New()
    hashData := append(request.Data, byte(time.Now().UnixNano()))
    h.Write(hashData)
    txid := h.Sum(nil)
    newRequestData := &serverRequest{txid, request}
    return send(newRequestData, server)
}

func send(req *serverRequest, server string) ([]byte, error) {
    postData, err := json.Marshal(req)
    resp, err := http.Post(server,
        "application/x-www-form-urlencoded",
        bytes.NewReader(postData))
    if err != nil {
        fmt.Println(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

func ServerUpdate(serverType int, add string)  {
    serverList[serverType] = add
}


