package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	model "github.com/server/transaction/model"
)

var DB_HOST = "host.docker.internal"
var DB_PORT = 5455
var DB_USER = "postgres"
var DB_PASSWORD = "Mystar123!@#"
var DB_NAME = "Transaction"
var DB_DISABLE_SSL = true

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
func InsertTransaction(transaction *model.TransactionReq) (string, error) {
	db := setupDB()

	var lastInsertID int
	query := `INSERT INTO trans (origin, user_id, amount, op_type, registered_at) VALUES($1, $2, $3, $4, $5) returning id`
	err := db.QueryRow(query, transaction.Origin, transaction.User_ID, transaction.Amount, transaction.Op_Type, transaction.Registered_At).Scan(&lastInsertID)
	if err != nil {
		return "error", err
	}

	// Select the inserted record and return
	return "One transaction inserted", nil
}

// Select returns the whole database
func SelectTransactions(pageInfo model.PageInfo, filterInfo model.TransactionReq) []*model.Transaction {
	pageNumber := pageInfo.Page_Number
	pageSize := pageInfo.Page_Size
	origin := filterInfo.Origin
	user_id := filterInfo.User_ID
	amount := filterInfo.Amount
	op_type := filterInfo.Op_Type
	registered_at := filterInfo.Registered_At

	db := setupDB()

	query := "SELECT id, origin, user_id, amount, op_type, registered_at FROM trans"
	isWhere := false

	if origin != "" {
		query += " WHERE origin='" + origin + "'"
		isWhere = true
	}

	if user_id != 0 {
		if isWhere == false {
			query += " WHERE user_id=" + strconv.Itoa(user_id)
			isWhere = true
		} else {
			query += " AND user_id=" + strconv.Itoa(user_id)
		}
	}

	if amount != "" {
		if isWhere == false {
			query += " WHERE amount='" + amount + "'"
			isWhere = true
		} else {
			query += " AND amount='" + amount + "'"
		}
	}

	if op_type != "" {
		if isWhere == false {
			query += " WHERE op_type='" + op_type + "'"
			isWhere = true
		} else {
			query += " AND op_type='" + op_type + "'"
		}
	}

	if registered_at != "" {
		if isWhere == false {
			query += " WHERE registered_at='" + registered_at + "'"
		} else {
			query += " AND registered_at='" + registered_at + "'"
		}
	}

	log.Printf("isWhere: %t", isWhere)

	if pageNumber != 0 || pageSize != 0 {
		query += " LIMIT " + strconv.Itoa(pageSize) + " OFFSET " + strconv.Itoa((pageNumber-1)*pageSize)
	}
	rows, err := db.Query(query)

	if err != nil {
		return nil
	}
	defer rows.Close()

	var transactions []*model.Transaction

	for rows.Next() {
		var id int
		var origin sql.NullString
		var user_id int
		var amount string
		var op_type sql.NullString
		var registered_at sql.NullString

		err = rows.Scan(&id, &origin, &user_id, &amount, &op_type, &registered_at)
		if err != nil {
			return nil
		}
		transactions = append(transactions, &model.Transaction{ID: id, Origin: origin.String, User_ID: user_id, Amount: amount, Op_Type: op_type.String, Registered_At: registered_at.String})
	}

	return transactions
}