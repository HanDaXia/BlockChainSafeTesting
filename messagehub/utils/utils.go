package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"errors"
)

func PostBytes(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Keep-Alive")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode<200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("receive bad response : %d", resp.StatusCode))
	}
	return ioutil.ReadAll(resp.Body)
}
