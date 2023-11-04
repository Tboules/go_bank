package main

import "net/http"

type ApiServer struct {
	listenAddress string
}

func NewApiServer(listenAddress string) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
	}
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleTranfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
