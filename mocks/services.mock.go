package mocks

import "pay-system/domain"

type MockThirdPartyService struct {
	MakePaymentFunc func(payment *domain.Payment) (*domain.PaymentDTO, error)
	GetPaymentFunc  func(payment *domain.Payment) (*domain.PaymentDTO, error)
}

// func NewMockThirdParty() (*MockThirdPartyService, error) {

// }

func (m *MockThirdPartyService) MakePayment(payment *domain.Payment) (*domain.PaymentDTO, error) {
	return m.MakePaymentFunc(payment)
}

func (m *MockThirdPartyService) GetPayment(payment *domain.Payment) (*domain.PaymentDTO, error) {
	return m.GetPaymentFunc(payment)
}

// mock wallet service
type MockWalletService struct {
	CreditWalletFunc      func(payment *domain.Payment) (*domain.Wallet, error)
	DebitWalletFunc       func(payment *domain.Payment) (*domain.Wallet, error)
	HandleTransactionFunc func(payment *domain.PaymentDTO) (*domain.PaymentDTO, error)
}

func (m *MockWalletService) CreditWallet(payment *domain.Payment) (*domain.Wallet, error) {
	return m.CreditWalletFunc(payment)
}

func (m *MockWalletService) DebitWallet(payment *domain.Payment) (*domain.Wallet, error) {
	return m.DebitWalletFunc(payment)
}

func (m *MockWalletService) HandleTransaction(payment *domain.PaymentDTO) (*domain.PaymentDTO, error) {
	return m.HandleTransactionFunc(payment)
}
