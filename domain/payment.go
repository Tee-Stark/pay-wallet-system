package domain

import "time"

const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
)

const (
	PaymentTypeCredit = "credit"
	PaymentTypeDebit  = "debit"
)

// Payment struct - for transactions
type Payment struct {
	ID        string    `json:"id,omitempty" gorm:"primary_key"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Amount    uint64    `json:"amount" gorm:"not null"` // in minimum units e.g kobo
	Type      string    `json:"type" gorm:"not null"`
	Status    string    `json:"status,omitempty" gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Payment) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return nil
}

type PaymentDTO struct {
	Reference string `json:"reference,omitempty"`
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
	Type      string `json:"type"` // to track transaction type
}
