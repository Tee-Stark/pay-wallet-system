package scripts

import (
	"log"
	"pay-system/domain"
	"pay-system/utils"

	"github.com/jinzhu/gorm"
)

func SeedUsers(db *gorm.DB) {
	// run process in an ACID transaction
	tx := db.Begin()

	// Defer a rollback in case of errors
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	users := []domain.User{
		{ID: utils.GenerateUUID(), Name: "Timmy O.", Email: "timmyo@example.com"},
		{ID: utils.GenerateUUID(), Name: "Jane S.", Email: "janes@example.com"},
		{ID: utils.GenerateUUID(), Name: "John D.", Email: "johnd@example.com"},
		{ID: utils.GenerateUUID(), Name: "Mary J.", Email: "maryj@example.com"},
	}

	// Create users and a wallet for each user
	for _, user := range users {
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to create user: %v", err)
		}

		// Create a wallet for each user
		wallet := domain.Wallet{
			UserID:  user.ID,
			Balance: 100,
		}

		if err := tx.Create(&wallet).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to create wallet for user %s: %v", user.Name, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Users and wallets seeded successfully!")
}
