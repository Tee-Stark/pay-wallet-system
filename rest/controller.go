package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"pay-system/domain"
	"pay-system/ports"
	"pay-system/utils"
)

// wallet controller
type WalletCtrl struct {
	Svc ports.IWalletService
}

func NewWalletCtrl(svc ports.IWalletService) *WalletCtrl {
	return &WalletCtrl{
		Svc: svc,
	}
}

func (c *WalletCtrl) HandleTransaction(w http.ResponseWriter, r *http.Request) {

	var req domain.PaymentDTO

	//decode request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		utils.HandleResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return

	}

	resp, err := c.Svc.HandleTransaction(&req)

	if err != nil {
		log.Println(err)
		utils.HandleResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.HandleResponse(w, http.StatusOK, resp)
}
