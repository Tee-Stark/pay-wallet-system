package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(segments)
		if len(segments) > 3 {
			// handle get payment request
			paymentID := segments[3]
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
		} else if r.URL.Path == "/third-party/payments" {
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
			  "reference": "%s",
			  "type": "%s"
			}`, paymentRequest.AccountID, paymentRequest.Amount, paymentRequest.Reference, paymentRequest.Type)

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
	// headers
	contentType := "application/json"
	authorization := "Bearer " + p.apiKey

	// request body
	data := fmt.Sprintf(`{
	  "account_id": "%s",
	  "amount": %d,
	  "reference": "%s",
	  "type": "%s"
	}`, payment.UserID, (payment.Amount * 100), payment.ID, payment.Type)

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
	var paymentDTO domain.PaymentDTO
	if err := json.NewDecoder(resp.Body).Decode(&paymentDTO); err != nil {
		return nil, err
	}

	return &paymentDTO, nil
}

func (p *StarkPayProvider) GetPayment(payment *domain.Payment) (*domain.PaymentDTO, error) {
	// headers
	contentType := "application/json"
	authorization := "Bearer " + p.apiKey

	req, err := http.NewRequest("GET", p.Server.URL+"/third-party/payments"+payment.ID, nil)
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
