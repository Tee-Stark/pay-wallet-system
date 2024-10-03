package repository

import (
	"pay-system/domain"
	"pay-system/utils"

	"github.com/jinzhu/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

// User methods
func (r *PostgresRepo) GetUser(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepo) GetPayment(id string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.Where("id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PostgresRepo) CreatePayment(p *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
	var newPayment *domain.Payment

	p.ID = utils.GenerateRef()
	p.Status = domain.PaymentStatusPending

	query := `
		INSERT INTO payments (id, user_id, amount, type, status)
		VALUES (?, ?, ?, ?, ?)
		RETURNING *
	`

	err := r.db.Raw(query, p.ID, p.UserID, p.Amount, p.Type, p.Status).Scan(&newPayment).Error

	if err != nil {
		return nil, err
	}
	return newPayment, nil
}

func (r *PostgresRepo) UpdatePayment(p *domain.Payment, tx *gorm.DB) (*domain.Payment, error) {
	var db *gorm.DB
	var updated *domain.Payment

	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	query := `
	UPDATE payments WHERE reference = ? SET status = ?
	RETURNING *
	`
	err := db.Raw(query, p.ID, p.Status).Scan(updated).Error
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *PostgresRepo) GetWallet(userID string, tx *gorm.DB) (*domain.Wallet, error) {
	var db *gorm.DB
	var wallet domain.Wallet
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	if err := db.Where("id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *PostgresRepo) UpdateWallet(w *domain.Wallet, tx *gorm.DB) (*domain.Wallet, error) {
	var db *gorm.DB
	var updated *domain.Wallet

	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	query := `
	UPDATE wallets WHERE user_id = ? SET balance = ?
	RETURNING *
	`
	if err := db.Raw(query, w.UserID, w.Balance).Error; err != nil {
		return nil, err
	}
	return updated, nil
}
