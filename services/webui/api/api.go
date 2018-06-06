package api

import (
	"github.com/opentracing/opentracing-go"
	"github.com/ycrxun/onion/services/account/proto"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/ycrxun/onion/tracing"
	"github.com/gorilla/mux"
)

type Server struct {
	tracer  opentracing.Tracer
	account account.AccountServiceClient
}

func NewServer(tr opentracing.Tracer, account account.AccountServiceClient) *Server {
	return &Server{
		account: account,
		tracer:  tr,
	}
}

func (s *Server) Run(port int) error {

	r := tracing.NewServeMux(s.tracer)
	r.Handle("/v1/accounts", http.HandlerFunc(s.accountsHandler))
	r.Handle("/v1/accounts/{id}", http.HandlerFunc(s.accountHandler))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (s *Server) accountsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()
	pageSize, pageToken := r.URL.Query().Get("pageSize"),
		r.URL.Query().Get("pageToken")

	if pageSize == "" || pageToken == "" {
		http.Error(w, "Please specify pageSize/pageToken params", http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(pageSize)
	if err != nil {
		http.Error(w, "Please specify pageSize params", http.StatusBadRequest)
		return
	}

	response, err := s.account.List(ctx, &account.ListAccountsRequest{
		PageSize:  int32(size),
		PageToken: pageToken,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts": response.Accounts,
	})
}

func (s *Server) accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	vars := mux.Vars(r)
	a, err := s.account.GetById(ctx, &account.GetByIdRequest{
		Id: vars["id"],
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(a)
}
