package service

import (
	"log"
	"pay-system/domain"
	"pay-system/ports"

	"github.com/jinzhu/gorm"
)

// wallet factory: main functions are CreateWallet and DebitWallet
type WalletService struct {
	repo ports.IRepository
	db   *gorm.DB
}

func NewWalletService(repo ports.IRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) DebitWallet(payment *domain.Payment) (*domain.Wallet, error) {
	var wallet *domain.Wallet

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// create a payment
	payment, err := s.repo.CreatePayment(payment, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to create payment: %v", err)
	}
	log.Printf("payment created: %v", payment)

	// update wallet balance
	wallet, err = s.repo.GetWallet(payment.UserID, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to get wallet: %v", err)
	}

	wallet.Balance = wallet.Balance - int(payment.Amount)
	updated, err := s.repo.UpdateWallet(wallet, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to update wallet: %v", err)
	}
	log.Printf("wallet updated %v", updated)

	// mark payment as completed
	payment.Status = domain.PaymentStatusCompleted
	updatedPayment, err := s.repo.UpdatePayment(payment, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to update payment status: %v", err)
	}
	log.Printf("payment updated: %v", updatedPayment)

	return wallet, nil
}

// for credit transactions
func (s *WalletService) CreditWallet(payment *domain.Payment) (*domain.Wallet, error) {
	var wallet *domain.Wallet

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// create a payment
	payment, err := s.repo.CreatePayment(payment, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to create payment: %v", err)
	}
	log.Printf("payment created: %v", payment)

	// update wallet balance
	wallet, err = s.repo.GetWallet(payment.UserID, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to get wallet: %v", err)
	}

	wallet.Balance = wallet.Balance + int(payment.Amount)
	updated, err := s.repo.UpdateWallet(wallet, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to update wallet balance: %v", err)
	}
	log.Printf("wallet updated %v", updated)

	// mark payment as completed
	payment.Status = domain.PaymentStatusCompleted
	updatedPayment, err := s.repo.UpdatePayment(payment, tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("failed to update payment status: %v", err)
	}
	log.Printf("payment updated: %v", updatedPayment)

	return wallet, nil
}
