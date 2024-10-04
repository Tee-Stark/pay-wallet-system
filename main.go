package main

import (
	"log"
	"net/http"
	"pay-system/config"
	"pay-system/ports"
	"pay-system/providers"
	"pay-system/repository"
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

	// initialize payment provider
	starkPayProvider := providers.NewStarkPayProvider(config.STARK_PAY_API_KEY)
	defer starkPayProvider.Server.Close()

	// Initialize services
	walletService := service.NewWalletService(repo, db, starkPayProvider)

	// seed database
	// scripts.SeedUsers(db)

	// Initialize application
	app := NewApp(db, starkPayProvider, walletService)

	// Initialize routes
	r := SetupRoutes(app)

	// Create a server instance
	server := &http.Server{
		Addr:    ":3534",
		Handler: r,
	}

	log.Println("Server started on port 3534")
	log.Fatal(server.ListenAndServe())
}
