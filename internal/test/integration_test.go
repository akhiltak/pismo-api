//go:build integration

package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/akhiltak/pismo-api/internal/storage/models"
	"github.com/akhiltak/pismo-api/pkg/api"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// test config
const baseURL = "http://localhost:2091"
const dbConnStr = "postgres://pismo:pismo@localhost:5433/pismo_test_db?sslmode=disable"

func TestMain(m *testing.M) {
	// Wait for the application to be ready
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Set up database connection using Bun
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbConnStr)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Run tests
	code := m.Run()

	// Clean up data after tests
	tearDown(db)
	defer db.Close()

	// Exit
	os.Exit(code)
}

func tearDown(db *bun.DB) {
	// Delete all data from tables
	ctx := context.Background()
	tables := []interface{}{
		(*models.Transaction)(nil),
		(*models.Account)(nil),
	}

	for _, table := range tables {
		_, err := db.NewDelete().Model(table).WhereOr("1=1").Exec(ctx)
		if err != nil {
			fmt.Printf("Failed to clean up table %T: %v\n", table, err)
		} else {
			fmt.Printf("Cleaned up table %T\n", table)
		}
	}
}

func TestCreateAccount(t *testing.T) {
	payload := api.CreateAccountRequest{
		DocNum: "12345678",
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(baseURL+"/accounts", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var account models.Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	assert.NoError(t, err)
	assert.Equal(t, "12345678", account.DocNum)
}

func TestGetAccount(t *testing.T) {
	// First, create an account
	createPayload := api.CreateAccountRequest{DocNum: "87654321"}
	jsonPayload, _ := json.Marshal(createPayload)
	createResp, err := http.Post(baseURL+"/accounts", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	var createdAccount models.Account
	json.NewDecoder(createResp.Body).Decode(&createdAccount)

	// Now, get the created account
	resp, err := http.Get(fmt.Sprintf("%s/accounts/%d", baseURL, createdAccount.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var account models.Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	assert.NoError(t, err)
	assert.Equal(t, createdAccount.ID, account.ID)
	assert.Equal(t, "87654321", account.DocNum)
}

func TestCreateTransaction(t *testing.T) {
	// First, create an account
	createAccountPayload := api.CreateAccountRequest{DocNum: "11223344"}
	jsonPayload, _ := json.Marshal(createAccountPayload)
	createResp, err := http.Post(baseURL+"/accounts", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	var createdAccount models.Account
	json.NewDecoder(createResp.Body).Decode(&createdAccount)

	// Now, create a transaction
	createTransactionPayload := api.CreateTransactionRequest{
		AccountID:       createdAccount.ID,
		OperationTypeID: 1, // debit type
		Amount:          decimal.NewFromFloat(100.50),
	}
	jsonPayload, _ = json.Marshal(createTransactionPayload)
	resp, err := http.Post(baseURL+"/transactions", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var transaction models.Transaction
	err = json.NewDecoder(resp.Body).Decode(&transaction)
	assert.NoError(t, err)
	assert.Equal(t, createdAccount.ID, transaction.AccountID)
	assert.Equal(t, int64(1), transaction.OperationTypeID)
	assert.Equal(t, decimal.NewFromFloat(-100.50), transaction.Amount)

	// Now, create a transaction - credit
	createTransactionPayload = api.CreateTransactionRequest{
		AccountID:       createdAccount.ID,
		OperationTypeID: 4, // credit type
		Amount:          decimal.NewFromFloat(1000.50),
	}
	jsonPayload, _ = json.Marshal(createTransactionPayload)
	resp, err = http.Post(baseURL+"/transactions", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&transaction)
	assert.NoError(t, err)
	assert.Equal(t, createdAccount.ID, transaction.AccountID)
	assert.Equal(t, int64(4), transaction.OperationTypeID)
	assert.Equal(t, decimal.NewFromFloat(1000.50), transaction.Amount)

	// Create transaction failure - incorrect account
	createTransactionPayload = api.CreateTransactionRequest{
		AccountID:       1000,
		OperationTypeID: 2, // debit type
		Amount:          decimal.NewFromFloat(100.50),
	}
	jsonPayload, _ = json.Marshal(createTransactionPayload)
	resp, err = http.Post(baseURL+"/transactions", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Create transaction failure - invalid operation type
	createTransactionPayload = api.CreateTransactionRequest{
		AccountID:       createdAccount.ID,
		OperationTypeID: 5, // invalid type
		Amount:          decimal.NewFromFloat(100.50),
	}
	jsonPayload, _ = json.Marshal(createTransactionPayload)
	resp, err = http.Post(baseURL+"/transactions", "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetTransactions(t *testing.T) {
	resp, err := http.Get(baseURL + "/transactions")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var transactions []models.Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	assert.NoError(t, err)
	assert.NotEmpty(t, transactions)
	assert.Equal(t, 2, len(transactions))
}
