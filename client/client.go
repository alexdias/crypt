package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type httpClient struct {
	url string
}

type StoreResponse struct {
	Key string `json:"key"`
	Err string `json:"error"`
}

type RetrieveResponse struct {
	Data string `json:"plaintext"`
	Err  string `json:"error"`
}

func NewHTTPClient(u string) Client {
	return &httpClient{
		url: u,
	}
}

var (
	ErrPerformingReq = errors.New("error performing request to server")
	ErrDecodingJson  = errors.New("error decoding json response")
)

func (c *httpClient) Store(id, payload []byte) (aesKey []byte, err error) {
	json_params := fmt.Sprintf(`{"id": "%v", "plaintext": "%v"}`, string(id), string(payload))
	json_bytes := []byte(json_params)
	endpoint_url := c.url + "store"
	req, err := http.NewRequest("POST", endpoint_url, bytes.NewBuffer(json_bytes))
	req.Header.Set("Content-Type", "application/json; content-type=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrPerformingReq
	}
	defer resp.Body.Close()
	var response StoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, ErrDecodingJson
	}
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	key_bytes := []byte(response.Key)
	return key_bytes, nil
}

func (c *httpClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	endpoint_url := c.url + "retrieve"
	req, err := http.NewRequest("GET", endpoint_url, nil)
	q := req.URL.Query()
	q.Add("id", string(id))
	q.Add("key", string(aesKey))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrPerformingReq
	}
	defer resp.Body.Close()
	var response RetrieveResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, ErrDecodingJson
	}
	if response.Err != "" {
		return nil, errors.New(response.Err)
	}
	return []byte(response.Data), nil
}
