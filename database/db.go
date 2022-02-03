package db

import (
	"database/sql"
	"fmt"
	"log"
	model "wallet/models"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

// If row doesn't exist return NoMatch
var NoMatch = fmt.Errorf("no matching record")

type Database struct {
	Connection *sql.DB
}

func Initialize(username string, password string, database string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	conn, err := sql.Open("postgres", dsn)
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

func (db Database) FetchData() (*model.AccountDataList, error) {
	list := &model.AccountDataList{}

	rows, err := db.Connection.Query("SELEFT * FROM account_data")

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var data model.AccountData
		err := rows.Scan(&data.ID, &data.Name, &data.Currency, &data.Balance, &data.CreatedAt)

		if err != nil {
			return list, err
		}

		list.Accounts = append(list.Accounts, data)
	}

	return list, nil
}

func (db Database) CreateData(account *model.AccountData) error {
	var id int
	var createdAt string

	query := `INSERT INTO account_data (name, currency, balance) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Connection.QueryRow(query, account.Name, account.Currency, account.Balance).Scan(&id, &createdAt)

	if err != nil {
		return err
	}

	account.ID = id
	account.CreatedAt = createdAt

	return nil
}

func (db Database) FetchDataByID(id int) (model.AccountData, error) {
	data := model.AccountData{}
	query := `SELECT * FROM account_data WHERE name = $1;`

	row := db.Connection.QueryRow(query, id)
	err := row.Scan(&data.ID, &data.Name, &data.Currency, &data.Balance, &data.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return data, NoMatch
		}
		return data, err
	}

	return data, nil
}

func (db Database) DeleteData(id int) error {
	query := `DELETE FROM account_data WHERE id = $1;`
	_, err := db.Connection.Exec(query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return NoMatch
		}
		return err
	}

	return nil
}

func (db Database) UpdateData(id int, data model.AccountData) (model.AccountData, error) {
	item := model.AccountData{}
	query := `UPDATE account_data SET balance=$1 WHERE name=$4 RETURNING id, name, description, created_at;`
	err := db.Connection.QueryRow(query, data.Balance, data.Name, data.Currency, id).Scan(&item.ID, &item.Name, &item.Currency, &item.Balance, &item.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return item, NoMatch
		}
		return item, err
	}
	return item, nil
}
