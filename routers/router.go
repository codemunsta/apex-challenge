package routers

import (
	handler "apex-challenge/handlers"
	mWare "apex-challenge/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/signin", handler.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/wallet/fund", mWare.KnownUser(handler.FundWallet)).Methods(http.MethodPost)
	router.HandleFunc("/wallet/balance", mWare.KnownUser(handler.GetWalletBalance)).Methods(http.MethodGet)
	router.HandleFunc("/start", mWare.KnownUser(handler.StartGame)).Methods(http.MethodPost)
	router.HandleFunc("/roll", mWare.KnownUser(mWare.ActiveSession(handler.RollDice))).Methods(http.MethodPost)
	router.HandleFunc("/stop", mWare.KnownUser(mWare.ActiveSession(handler.EndGame))).Methods(http.MethodPost)

	// bonus
	router.HandleFunc("/active", mWare.KnownUser(handler.ActiveSession)).Methods(http.MethodGet)
	router.HandleFunc("/get/transactions", mWare.KnownUser(handler.GetTransactionList)).Methods(http.MethodGet)
	return router
}
