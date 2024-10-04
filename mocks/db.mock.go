package mocks

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type MockDB struct {
	DB   *gorm.DB
	Mock sqlmock.Sqlmock
}

func NewMockDB() (*MockDB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, err
	}

	return &MockDB{
		DB:   gormDB,
		Mock: mock,
	}, nil
}

// for transactions
func (m *MockDB) Begin() *gorm.DB {
	tx := m.DB.Begin()
	return tx
}

func (m *MockDB) Commit() *gorm.DB {
	return m.DB.Commit()
}

func (m *MockDB) Rollback() *gorm.DB {
	return m.DB.Rollback()
}

func (m *MockDB) ExpectBegin() {
	m.Mock.ExpectBegin()
}

func (m *MockDB) ExpectCommit() {
	m.Mock.ExpectCommit()
}

func (m *MockDB) ExpectRollback() {
	m.Mock.ExpectRollback()
}

func (m *MockDB) ExpectGetWallet(userID string, balance int, ID uint64) {
	m.Mock.ExpectQuery("SELECT * FROM `wallets` WHERE user_id = ?").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance", "id"}).AddRow(userID, balance, ID))
}

func (m *MockDB) ExpectUpdateWallet(userID string, newBalance int) {
	m.Mock.ExpectExec("UPDATE `wallets` SET `balance`=? WHERE user_id = ?").
		WithArgs(newBalance, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func (m *MockDB) ExpectUpdatePayment(paymentID string, status string) {
	m.Mock.ExpectExec("UPDATE `payments` SET `status`=? WHERE id = ?").
		WithArgs(status, paymentID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func (m *MockDB) CheckExpectations() error {
	return m.Mock.ExpectationsWereMet()
}
