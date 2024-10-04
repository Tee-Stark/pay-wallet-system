package service

import (
	"errors"
	"log"
	"pay-system/domain"
	"pay-system/mocks"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestWalletService_DebitWallet(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"

	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	mockDB.ExpectGetWallet(userID, 1000, 1)
	mockDB.ExpectUpdateWallet(userID, 500)
	mockDB.ExpectUpdatePayment("some-payment-id", domain.PaymentStatusCompleted)

	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return &domain.Wallet{UserID: userID, Balance: 1000}, nil
		},
		UpdateWalletFunc: func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
			return wallet, nil
		},
		UpdatePaymentFunc: func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
			return payment, nil
		},
	}

	mockThirdParty := &mocks.MockThirdPartyService{
		MakePaymentFunc: func(payment *domain.Payment) (*domain.PaymentDTO, error) {
			dto := &domain.PaymentDTO{
				AccountID: payment.UserID,
				Reference: payment.ID,
				Amount:    payment.Amount,
				Type:      payment.Type,
			}
			return dto, nil
		},
	}

	service := NewWalletService(mockRepo, mockDB.DB, mockThirdParty)
	payment := &domain.Payment{UserID: userID, Amount: 500}

	wallet, err := service.DebitWallet(payment)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if wallet.Balance != 500 {
		t.Errorf("Expected wallet balance to be 500, got %v", wallet.Balance)
	}
}

func TestWalletService_DebitWallet_InsufficientBalance(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"

	mockDB.ExpectBegin()
	mockDB.ExpectGetWallet(userID, 300, 1)
	mockDB.ExpectRollback() // Rollback should occur

	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return &domain.Wallet{UserID: userID, Balance: 300, ID: 1}, nil
		},
		UpdateWalletFunc: func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
			return wallet, nil
		},
	}

	service := NewWalletService(mockRepo, mockDB.DB, nil)
	payment := &domain.Payment{UserID: userID, Amount: 500} // Attempt to debit 500

	wallet, err := service.DebitWallet(payment)

	if wallet != nil {
		t.Fatalf("Expected no wallet returned, got %v", wallet)
	}
	if err.Error() != "insufficient balance" {
		t.Fatalf("Expected insufficient balance error, got %v", err)
	}
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestWalletService_DebitWallet_GetWalletError(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"

	mockDB.ExpectBegin()
	mockDB.ExpectGetWallet(userID, 300, 1)
	mockDB.ExpectRollback()
	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return nil, errors.New("failed to get wallet")
		},
	}

	service := NewWalletService(mockRepo, mockDB.DB, nil)
	payment := &domain.Payment{UserID: userID, Amount: 500}

	wallet, err := service.DebitWallet(payment)

	if wallet != nil {
		t.Fatalf("Expected no wallet returned, got %v", wallet)
	}
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if err.Error() != "failed to get wallet" {
		t.Fatalf("Expected error 'failed to get wallet', got: %v", err)
	}
}

func TestWalletService_DebitWallet_UpdateWalletError(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"

	mockDB.ExpectBegin()
	mockDB.ExpectGetWallet(userID, 1000, 1)
	mockDB.ExpectRollback() // Rollback should occur
	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return &domain.Wallet{UserID: userID, Balance: 1000}, nil
		},
		UpdateWalletFunc: func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
			return nil, errors.New("failed to update wallet")
		},
	}

	service := NewWalletService(mockRepo, mockDB.DB, nil)
	payment := &domain.Payment{UserID: userID, Amount: 500}

	wallet, err := service.DebitWallet(payment)

	if wallet != nil {
		t.Fatalf("Expected no wallet returned, got %v", wallet)
	}
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if err.Error() != "failed to update wallet" {
		t.Fatalf("Expected error 'failed to update wallet', got: %v", err)
	}
}

func TestWalletService_DebitWallet_UpdatePaymentError(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"

	mockDB.ExpectBegin()
	mockDB.ExpectGetWallet(userID, 1000, 1)
	mockDB.ExpectUpdateWallet(userID, 500)
	mockDB.ExpectUpdatePayment("pref_mock", domain.PaymentStatusCompleted)
	mockDB.ExpectRollback() // Rollback should occur
	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return &domain.Wallet{UserID: userID, Balance: 1000, ID: 1}, nil
		},
		UpdateWalletFunc: func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
			return wallet, nil
		},
		UpdatePaymentFunc: func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
			return nil, errors.New("failed to update payment")
		},
	}

	service := NewWalletService(mockRepo, mockDB.DB, nil)
	payment := &domain.Payment{ID: "pref_mock", UserID: userID, Amount: 500}

	wallet, err := service.DebitWallet(payment)

	if wallet != nil {
		t.Fatalf("Expected no wallet returned, got %v", wallet)
	}
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if err.Error() != "failed to update payment" {
		t.Fatalf("Expected error 'failed to get wallet', got: %v", err)
	}
}

func TestHandleTransaction_Success(t *testing.T) {
	mockDB, err := mocks.NewMockDB()
	if err != nil {
		t.Fatalf("Failed to start mock DB")
	}
	defer mockDB.DB.Close()

	userID := "5f47c1ab-e6f2-4b5b-b43a-5f182fe8f69e"
	payID := "pref_mock"

	paymentDto := &domain.PaymentDTO{
		AccountID: userID,
		Amount:    1000,
		Reference: payID,
		Type:      domain.PaymentTypeCredit,
	}

	mockDB.ExpectBegin()
	mockDB.ExpectCommit()
	mockDB.ExpectGetWallet(userID, 1000, 1)
	mockDB.ExpectUpdatePayment(payID, "completed")
	mockDB.ExpectUpdateWallet(userID, 1000)

	mockRepo := &mocks.MockRepository{
		GetWalletFunc: func(userID string, tx *gorm.DB) (*domain.Wallet, error) {
			return &domain.Wallet{UserID: userID, Balance: 1000, ID: 1}, nil
		},
		UpdatePaymentFunc: func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
			return payment, nil
		},
		UpdateWalletFunc: func(wallet *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
			return wallet, nil
		},
		CreatePaymentFunc: func(payment *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
			return payment, nil
		},
	}

	mockPayProvider := &mocks.MockThirdPartyService{
		MakePaymentFunc: func(payment *domain.Payment) (*domain.PaymentDTO, error) {
			return &domain.PaymentDTO{
				AccountID: payment.UserID,
				Reference: payment.ID,
				Amount:    payment.Amount,
				Type:      payment.Type,
			}, err
		},
		GetPaymentFunc: func(payment *domain.Payment) (*domain.PaymentDTO, error) {
			return &domain.PaymentDTO{
				AccountID: payment.UserID,
				Reference: payment.ID,
				Amount:    payment.Amount,
				Type:      payment.Type,
			}, err
		},
	}

	walletSvc := NewWalletService(mockRepo, mockDB.DB, mockPayProvider)

	resp, err := walletSvc.HandleTransaction(paymentDto)
	log.Println(resp)
	if err != nil {
		t.Fatalf("Expected no errors but got: %v", err)
	}
	// if !reflect.DeepEqual(resp, paymentDto) {
	// 	t.Fatalf("Incorrect response from transaction: %v", resp)
	// }
}
