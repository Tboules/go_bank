package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	store         Storage
}

func NewApiServer(listenAddress string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleGetAccountByID))

	log.Println("JSON API running on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)

}

func (s *ApiServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Anthony", "Boules")
	vars := mux.Vars(r)

	fmt.Printf("id: %v\n", vars["id"])

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountParams := new(CreateAccountParams)

	err := json.NewDecoder(r.Body).Decode(createAccountParams)
	if err != nil {
		return err
	}

	account := NewAccount(createAccountParams.FirstName, createAccountParams.LastName)

	dbAcc, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, dbAcc)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleTranfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}
	}
}
