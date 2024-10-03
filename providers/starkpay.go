package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pay-system/domain"
	"strings"
)

type StarkPayProvider struct {
	apiKey       string
	client       *http.Client
	Server       *httptest.Server
	paymentStore map[string]domain.PaymentDTO
}

func NewStarkPayProvider(apiKey string) *StarkPayProvider {
	// in-memory store for payment responses
	store := make(map[string]domain.PaymentDTO)

	// create mock http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handle get payment request
		segments := strings.Split(r.URL.Path, "/")
		if len(segments) > 2 {
			// handle get payment request
			paymentID := segments[2]
			paymentData := store[paymentID]
			response := fmt.Sprintf(`{
			  "reference": "%s",
			  "account_id": "%s",
			  "amount": %d
			}`, paymentID, paymentData.AccountID, paymentData.Amount)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
			return
		}

		// handle create payment request
		if r.URL.Path == "/third-party/payments" {
			// handle payment request
			var paymentRequest domain.PaymentDTO
			err := json.NewDecoder(r.Body).Decode(&paymentRequest)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			store[paymentRequest.Reference] = paymentRequest
			// mock response
			response := fmt.Sprintf(`{
			  "account_id": "%s",
			  "amount": %d,
			  "reference": "%s"
			}`, paymentRequest.AccountID, paymentRequest.Amount, paymentRequest.Reference)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
			return
		}

	}))
	return &StarkPayProvider{
		apiKey:       apiKey,
		Server:       server,
		client:       &http.Client{},
		paymentStore: store,
	}
}

func (p *StarkPayProvider) MakePayment(payment *domain.Payment) (*domain.PaymentDTO, error) {
	/**
	 *	create mock http POST request to /third-party/payments
	 */
	// headers
	contentType := "application/json"
	authorization := "Bearer " + p.apiKey

	// request body
	data := fmt.Sprintf(`{
	  "account_id": "%s",
	  "amount": %d,
	  "reference": "%s"
	}`, payment.UserID, payment.Amount, payment.ID)

	req, err := http.NewRequest("POST", p.Server.URL+"/third-party/payments", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// unmarshal response into PaymentDTO struct
	var PaymentDTO domain.PaymentDTO
	if err := json.NewDecoder(resp.Body).Decode(&PaymentDTO); err != nil {
		return nil, err
	}

	return &PaymentDTO, nil
}

func (p *StarkPayProvider) GetPayment(payment *domain.Payment) (*domain.PaymentDTO, error) {
	// Implement logic to create a credit payment using Stark Pay API
	return nil, nil
}
