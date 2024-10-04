package mocks

import (
	"pay-system/domain"

	"github.com/jinzhu/gorm"
)

type MockRepository struct {
	GetUserFunc       func(id string) (*domain.User, error)
	GetWalletFunc     func(userID string, tx *gorm.DB) (*domain.Wallet, error)
	UpdateWalletFunc  func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error)
	GetPaymentFunc    func(id string) (*domain.Payment, error)
	UpdatePaymentFunc func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error)
	CreatePaymentFunc func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error)
}

func (m *MockRepository) GetUser(id string) (*domain.User, error) {
	return m.GetUserFunc(id)
}

func (m *MockRepository) GetWallet(userID string, tx *gorm.DB) (*domain.Wallet, error) {
	return m.GetWalletFunc(userID, tx)
}

func (m *MockRepository) UpdateWallet(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
	return m.UpdateWalletFunc(wallet, tx)
}

func (m *MockRepository) GetPayment(id string) (*domain.Payment, error) {
	return m.GetPaymentFunc(id)
}

func (m *MockRepository) UpdatePayment(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
	return m.UpdatePaymentFunc(payment, tx)
}

func (m *MockRepository) CreatePayment(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
	return m.CreatePaymentFunc(payment, tx)
}
