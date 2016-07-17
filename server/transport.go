package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var (
	ErrBadRouting               = errors.New("bad mapping between route and handler")
	ErrMissingParameters        = errors.New("query parameters are missing")
	ErrMissingIdArgument        = errors.New("id POST argument must be specified")
	ErrMissingPlaintextArgument = errors.New("plaintext POST argument must be specified")
)

type postDataRequest struct {
	PlainTextData PlainTextData
}

type getDataRequest struct {
	ID  string
	Key []byte
}

type postDataResponse struct {
	Key string `json:"key"`
	Err error  `json:"err,omitempty"`
}

type getDataResponse struct {
	Data string `json:"plaintext"`
	Err  error  `json:"err,omitempty"`
}

func SetUpHTTPHandlers(ctx context.Context, s StoreService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		ctx,
		makePostDataEndpoint(s),
		decodePostDataRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		ctx,
		makeGetDataEndpoint(s),
		decodeGetDataRequest,
		encodeResponse,
		options...,
	))
	return r
}

func makePostDataEndpoint(s StoreService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postDataRequest)
		key, e := s.PostData(ctx, req.PlainTextData)
		string_key := string(key)
		return postDataResponse{Key: string_key, Err: e}, nil
	}
}

func makeGetDataEndpoint(s StoreService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDataRequest)
		data, e := s.GetData(ctx, req.ID, req.Key)
		data_string := string(data)
		return getDataResponse{Data: data_string, Err: e}, nil
	}
}

func decodePostDataRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postDataRequest
	if e := json.NewDecoder(r.Body).Decode(&req.PlainTextData); e != nil {
		return nil, e
	}
	if req.PlainTextData.ID == "" {
		return nil, ErrMissingIdArgument
	}
	if req.PlainTextData.Data == "" {
		return nil, ErrMissingPlaintextArgument
	}
	return req, nil
}

func decodeGetDataRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	id := r.FormValue("id")
	key := r.FormValue("key")
	if id == "" || key == "" {
		return nil, ErrMissingParameters
	}
	byte_key := []byte(key)
	return getDataRequest{
		ID:  id,
		Key: byte_key,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
