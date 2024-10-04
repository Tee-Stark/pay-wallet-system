package ports

import "pay-system/domain"

type IThirdPartyService interface {
	MakePayment(payment *domain.Payment) (*domain.PaymentDTO, error)
	GetPayment(payment *domain.Payment) (*domain.PaymentDTO, error)
	// possible function to act as webhook emitter
}

type IWalletService interface {
	CreditWallet(payment *domain.Payment) (*domain.Wallet, error)
	DebitWallet(payment *domain.Payment) (*domain.Wallet, error)
	HandleTransaction(payment *domain.Payment) (bool, error)
}
