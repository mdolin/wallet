package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	rows, err := db.Connection.Query("SELECT * FROM account")

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
	}
	response.Type = "success"
	json.NewEncoder(w).Encode(response)
}

func (db Database) FetchByName(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}
	response := model.Response{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM account WHERE name = $1;`
	row := db.Connection.QueryRow(query, account.Name)

	err = row.Scan(
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
	account := model.Account{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO account (name, currency, balance) VALUES ($1, $2, $3) RETURNING id, created_at;`
	_, err = db.Connection.Exec(
		query,
		account.Name,
		account.Currency,
		account.Balance,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (db Database) Deposit(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE account SET balance = balance + $1 WHERE name = $2 RETURNING id, name, currency, balance, created_at;`
	_, err = db.Connection.Exec(query,
		account.Balance,
		account.Name,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (db Database) Withdraw(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE account SET balance = balance - $1 WHERE name = $2 RETURNING id, name, currency, balance, created_at;`
	_, err = db.Connection.Exec(query,
		account.Balance,
		account.Name,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (db Database) DeleteByName(w http.ResponseWriter, r *http.Request) {
	account := model.Account{}

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `DELETE FROM account WHERE name = $1;`
	_, err = db.Connection.Exec(query, account.Name)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
