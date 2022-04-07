package main

import (
	"fmt"
	"net/http"
	"xendit/pkg/logger"
)

func (s *server) TransactionRecordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TransactionRecordHandler!")
}

func (s *server) TransactionGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TransactionGetHandler!")
}

func (s *server) TransactionFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TransactionFeedbackHandler!")
}

func (s *server) Run() error {
	// add necessary services
	s.addLogger()
	s.addTracer()
	s.addJwt()
	s.addRedis()
	s.addMySQL()
	s.addAuthenticationService()
	s.addTransactionService()

	// add apis
	s.mux.HandleFunc(TRANSACTION_RECORD, s.TransactionRecordHandler)
	s.mux.HandleFunc(TRANSACTION_FEEDBACK, s.TransactionFeedbackHandler)
	s.mux.HandleFunc(TRANSACTION_GET, s.TransactionGetHandler)

	s.logger.Debugf("App is listening %v--", s.cfg.Server.Port)
	return http.ListenAndServe(s.cfg.Server.Port, s.mux)
}

func (s *server) GetLogger() logger.Logger {
	return s.logger
}
