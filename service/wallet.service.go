package service

import (
	"log"
	"errors"
	"pay-system/domain"
	"pay-system/ports"

	"github.com/jinzhu/gorm"
)

// wallet factory: main functions are CreateWallet and DebitWallet
type WalletService struct {
	repo ports.IRepository
	db   *gorm.DB
	paymentProvider ports.IThirdPartyService
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

func (s *WalletService) HandleTransaction(dto *domain.PaymentDTO) (bool, error) {
	wallet, err := s.repo.GetWallet(dto.UserID, nil)
	success := false

	if err != nil {
		log.Println("Error getting wallet: %v". err)
		return success, err
	}

	if dto.Type == domain.PaymentTypeDebit && wallet.Balance - dto.Amount < 0 {
		return success, errors.New("Insufficient Balance")
	}
	// transaction process
	paymentData := &domain.Payment{
		UserID: dto.UserID,
		Amount: dto.Amount,
		Type: dto.Type
	}
	// create a payment
	payment, err := s.repo.CreatePayment(payment, nil)
	if err != nil {
		return success, err
	}
	log.Printf("payment created: %v", payment)
	
	// make request to third party
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.paymentProvider.MakePayment(payment)



	success = true
	return success, nil
}