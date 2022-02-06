package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	model "github.com/mdolin/wallet/models"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

type Database struct {
	Connection *sql.DB
}

func Initialize(username string, password string, database string) (Database, error) {
	db := Database{}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return db, err
	}

	db.Connection = conn
	err = db.Connection.Ping()
	if err != nil {
		return db, err
	}

	log.Println("Conection established")

	return db, nil
}

func (db Database) Fetch(w http.ResponseWriter, r *http.Request) {
	response := model.Response{}

	rows, err := db.Connection.Query("SELEFT * FROM account")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var account model.Account
		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Currency,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			panic(err)
		}

		response.Accounts = append(response.Accounts, account)
		response.Type = "success"

		json.NewEncoder(w).Encode(response)
	}
}

func (db Database) FetchByName(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	response := model.Response{}

	query := `SELECT * FROM account WHERE name = $1;`
	row := db.Connection.QueryRow(query, name)

	var account model.Account
	err := row.Scan(
		&account.ID,
		&account.Name,
		&account.Currency,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		panic(err)
	}

	response.Accounts = append(response.Accounts, account)
	response.Type = "success"

	json.NewEncoder(w).Encode(response)
}

func (db Database) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	currency := r.FormValue("currency")
	balance := r.FormValue("balance")

	var id int
	var createdAt string
	var bal float64

	response := model.Response{}

	if balance == "" {
		bal, _ = strconv.ParseFloat("0", 64)
	} else {
		bal, _ = strconv.ParseFloat(balance, 64)
	}

	query := `INSERT INTO account (name, currency, balance) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Connection.QueryRow(
		query,
		name,
		currency,
		bal).Scan(&id, &createdAt)

	if err != nil {
		panic(err)
	}

	response.Type = "success"
	response.Message = "Account has been created"

	json.NewEncoder(w).Encode(response)
}

func (db Database) Deposit(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	balance, _ := strconv.ParseFloat(r.FormValue("balance"), 64)

	response := model.Response{}

	var account model.Account

	query := `UPDATE account SET balance = balance + $1 WHERE name = $2 RETURNING id, name, currency, balance, created_at;`
	err := db.Connection.QueryRow(query,
		balance,
		name,
	).Scan(
		&account.ID,
		&account.Name,
		&account.Currency,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		panic(err)
	}

	response.Accounts = append(response.Accounts, account)
	response.Type = "success"

	json.NewEncoder(w).Encode(response)
}

func (db Database) Withdrawal(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	balance, _ := strconv.ParseFloat(r.FormValue("balance"), 64)

	response := model.Response{}

	var account model.Account

	query := `UPDATE account SET balance = balance - $1 WHERE name = $2 RETURNING id, name, currency, balance, created_at;`
	err := db.Connection.QueryRow(query,
		balance,
		name,
	).Scan(
		&account.ID,
		&account.Name,
		&account.Currency,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		panic(err)
	}

	response.Accounts = append(response.Accounts, account)
	response.Type = "success"

	json.NewEncoder(w).Encode(response)
}

func (db Database) DeleteByName(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	response := model.Response{}

	query := `DELETE FROM account WHERE id = $1;`
	_, err := db.Connection.Exec(query, name)

	if err != nil {
		panic(err)
	}

	response.Type = "success"
	response.Message = "Account has been deleted"

}
