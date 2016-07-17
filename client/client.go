package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient struct {
	url string
}

type StoreResponse struct {
	Key string `json:"key"`
	Err string `json:"err"`
}

func NewHTTPClient() Client {
	return &httpClient{
		url: "http://localhost:8080/",
	}
}

func (c *httpClient) Store(id, payload []byte) (aesKey []byte, err error) {
	json_parameters := fmt.Sprintf(`{"id": "%v", "plaintext": "%v"}`, string(id), string(payload))
	json_bytes := []byte(json_parameters)
	url := c.url + "store"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_bytes))
	req.Header.Set("Content-Type", "application/json; content-type=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response StoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	key_bytes := []byte(response.Key)
	return key_bytes, nil
}

func (c *httpClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	return nil, nil
}
