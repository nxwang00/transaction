package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	model "github.com/test/transaction/model"
)

var DB_HOST = "localhost"
var DB_PORT = 5432
var DB_USER = "postgres"
var DB_PASSWORD = "Mystar123!@#"
var DB_NAME = "Transaction"
var DB_DISABLE_SSL = false

var dbHandle *sql.DB

// DB set up
func setupDB() *sql.DB {
	if dbHandle == nil {
		dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

		if DB_DISABLE_SSL {
			dbinfo += " sslmode=disable"
		}

		dbHandle, _ = sql.Open("postgres", dbinfo)
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		if err := dbHandle.PingContext(ctx); err != nil {
			log.Fatalf("unable to connect to database: %v", err)
		} else {
			log.Printf("DB: %s Successful\n", dbinfo)
		}
	}

	return dbHandle
}

// Insert allows populating database
func InsertTransaction(transaction *model.TransactionReq) (*model.Transaction, error) {
	db := setupDB()

	var lastInsertID int
	query := `INSERT INTO trans (origin, user_id, amount, op_type, registered_at) VALUES($1, $2, $3, $4, $5) returning id`
	err := db.QueryRow(query, transaction.Origin, transaction.User_ID, transaction.Amount, transaction.Op_Type, transaction.Registered_At).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	// Select the inserted record and return
	return SelectTransaction(lastInsertID), nil
}

// Select the record with the id
func SelectTransaction(transactionId int) *model.Transaction {
	db := setupDB()
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT id, origin, user_id, amount, op_type, registered_At FROM trans WHERE id=$1", transactionId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var transaction *model.Transaction

	var id int
	var origin sql.NullString
	var user_id int
	var amount float64
	var op_type sql.NullString
	var registered_at sql.NullString

	if rows.Next() {
		err = rows.Scan(&id)
		if err == nil {
			transaction = &model.Transaction{ID: id, Origin: origin.String, User_ID: user_id, Amount: amount, Op_Type: op_type.String, Registered_At: registered_at.String}
		}
	}

	return transaction
}