package domain

type Wallet struct {
	ID      uint64 `json:"id" gorm:"primaryKey;autoincrement"`
	User    User   `json:"user" gorm:"foreignKey:UserID"`
	UserID  string `json:"user_id" gorm:"not null;unique;"`
	Balance int    `json:"balance" gorm:"not null"`
}
