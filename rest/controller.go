package rest

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
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return

	}

	resp, err := c.Svc.HandleTransaction(req)

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, resp)
		return
	}

	sendJSON(w, http.StatusOK, resp)

}



// sendJSON sends a JSON response with the specified status code and data
func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
