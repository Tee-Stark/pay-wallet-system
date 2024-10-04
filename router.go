package main

import (
	"github.com/gorilla/mux"
	"pay-system/rest"
)

func SetupRoutes(app *App) {
	router := mux.NewRouter()

	ctrl := rest.NewWalletCtrl(app.WalletSvc)

	router.HandlerFunc("/transaction", ctrl.HandleTransaction).methods("POST")
}
