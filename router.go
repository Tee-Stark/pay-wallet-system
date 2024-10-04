package main

import (
	"pay-system/rest"

	"github.com/gorilla/mux"
)

func SetupRoutes(app *App) *mux.Router {
	router := mux.NewRouter()

	ctrl := rest.NewWalletCtrl(app.WalletSvc)

	router.HandleFunc("/transaction", ctrl.HandleTransaction).Methods("POST")

	return router
}
