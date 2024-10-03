package main

import (
	"log"
	"net/http"
	"pay-system/config"
	"pay-system/ports"
	"pay-system/providers"
	"pay-system/repository"
	"pay-system/scripts"
	"pay-system/service"

	"github.com/jinzhu/gorm"
)

/**
* Application composition root
 */

type App struct {
	db              *gorm.DB
	PaymentProvider ports.IThirdPartyService
	WalletSvc       ports.IWalletService
}

func NewApp(db *gorm.DB, payProvider ports.IThirdPartyService, svc ports.IWalletService) *App {
	return &App{
		db:              db,
		PaymentProvider: payProvider,
		WalletSvc:       svc,
	}
}

func main() {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	// Initialize database
	db := config.NewPostgresDB()
	defer db.Close()

	// Initialize repositories
	repo := repository.NewPostgresRepo(db)

	// Initialize services
	walletService := service.NewWalletService(repo)
	// initialize payment provider
	starkPayProvider := providers.NewStarkPayProvider(config.STARK_PAY_API_KEY)
	defer starkPayProvider.Server.Close()

	// seed database
	scripts.SeedUsers(db)

	// Initialize application
	app := NewApp(db, starkPayProvider, walletService)

	// Initialize routes
	SetupRoutes(app)

	log.Println("Starting application server on port 3534")
	if err := http.ListenAndServe(":3534", nil); err != nil {
		log.Fatal(err)
	}
}
