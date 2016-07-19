package test

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/alexnvdias/crypt/client"
	"github.com/alexnvdias/crypt/server"
)

func TestMain(m *testing.M) {
	logger := log.NewLogfmtLogger(os.Stderr)
	go server.Run(":8080", logger)
	os.Exit(m.Run())
}

func setupClient() client.Client {
	return client.NewHTTPClient("http://localhost:8080/")
}

func TestBasic(t *testing.T) {
	cl := setupClient()
	payload := []byte("data")
	id := []byte("1")
	key, err := cl.Store(id, payload)
	if err != nil {
		t.Error("Error when storing payload")
	}

	retrieved_data, err_retrieve := cl.Retrieve(id, key)
	if err_retrieve != nil {
		t.Error("Error when retrieving data")
	}
	if string(retrieved_data) != string(payload) {
		t.Error("Retrieved data does not equal original payload")
	}
}

func TestCantOverwrite(t *testing.T) {
	cl := setupClient()
	payload1 := []byte("data1")
	id1 := []byte("2")
	_, err := cl.Store(id1, payload1)
	if err != nil {
		t.Error("Error when storing payload")
	}

	payload2 := []byte("data2")
	id2 := []byte("2")
	_, err2 := cl.Store(id2, payload2)
	if err2 != client.ErrAlreadyPresent {
		t.Error(err2)
	}
}

func TestOneKeyCantDecryptAnother(t *testing.T) {
	cl := setupClient()
	payload1 := []byte("data1")
	id3 := []byte("3")
	key1, err := cl.Store(id3, payload1)
	if err != nil {
		t.Error("Error when storing payload")
	}

	id4 := []byte("4")
	key2, err2 := cl.Store(id4, payload1)
	if err2 != nil {
		t.Error("Error when storing payload")
	}

	_, retrieve_err3 := cl.Retrieve(id3, key2)
	if retrieve_err3 != client.ErrDecryptingData {
		t.Error(retrieve_err3)
	}

	_, retrieve_err4 := cl.Retrieve(id4, key1)
	if retrieve_err4 != client.ErrDecryptingData {
		t.Error(retrieve_err4)
	}

	data1, err3 := cl.Retrieve(id3, key1)
	if string(data1) != string(payload1) {
		if err3 != nil {
			t.Error(err3)
		}
		t.Error("Retrieved data does not match")
	}
	data2, err4 := cl.Retrieve(id4, key2)
	if string(data2) != string(payload1) {
		if err4 != nil {
			t.Error(err4)
		}
		t.Error("Retrieved data does not match")
	}
}

func TestSpecialCharacters(t *testing.T) {
	cl := setupClient()
	payload := []byte("ÇÇóà~~âôâô")
	id := []byte("5")
	key, err := cl.Store(id, payload)
	if err != nil {
		t.Error("Error when storing payload")
	}

	retrieved_data, err_retrieve := cl.Retrieve(id, key)
	if err_retrieve != nil {
		t.Error("Error when retrieving data")
	}
	if string(retrieved_data) != string(payload) {
		t.Error("Retrieved data does not equal original payload")
	}
}
