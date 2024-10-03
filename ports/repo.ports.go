package ports

import (
	"pay-system/domain"

	"github.com/jinzhu/gorm"
)

type IRepository interface {
	IUserRepository
	IPaymentRepository
	IWalletRepository
}

type IUserRepository interface {
	GetUser(id string) (*domain.User, error)
	// CreateUser(user *domain.User) error
}

type IPaymentRepository interface {
	GetPayment(id string) (*domain.Payment, error)
	CreatePayment(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error)
}

type IWalletRepository interface {
	GetWallet(userID string, tx *gorm.DB) (*domain.Wallet, error)
	UpdateWallet(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error)
}
