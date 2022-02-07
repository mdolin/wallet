## 
The goal was to write a Wallet app where we can store the funds for each user. The Wallet app is responsible to store the funds and provide functionality for manipulating the balance.

The application is using "GET", "POST" and "DELETE" methods for creating, adding, removing funds and queri for the current state of the Wallet. All the data is written to postgres database.

I need to get familiar how Go programming language is communicating with postgres so I was using different resources. Resources are listed below in the [Useful resources](#useful-resources) section.

## Main bits of the project
* PostgreSQL database
* Go PostgreSQL driver for handling the database
* Eequest router for matching incoming requests
* Create, Fetch, and Delete operations

## Structure of the project
```
.
├── database
│   ├── db.go
│   └── tables
│       ├── insert.sql
│       └── table.down.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── models
│   └── model.go
└── README.md


```

## Requirements
To run the project and tests you will need
* [Go Programming language](https://go.dev/doc/install)
* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [lib/pq](https://pkg.go.dev/github.com/lib/pq) - to interact with PostgreSQL
* [gorilla/mux](https://pkg.go.dev/github.com/gorilla/mux) - for URL matcher and routing

## Usage
Build the image and start the Wallet app and PostgreSQL database.

### Building and running the image 
```
docker-compose up --build
```

Verify if they are running
```
docker-compose ps
```

Output:
```
     Name                    Command              State                    Ports                  
--------------------------------------------------------------------------------------------------
wallet            /usr/bin/wallet                 Up      0.0.0.0:8000->8000/tcp,:::8000->8000/tcp
wallet_database   docker-entrypoint.sh postgres   Up      0.0.0.0:5432->5432/tcp,:::5432->5432/tcp

```

If encounter an error 
```
Error starting userland proxy: listen tcp4 0.0.0.0:5432: bind: address already in use
```
it means that where you are trying to run the app PostgreSQL is already listening on that port. Similar can happen for port 8000.

You can shut down the services running on that port or change the port for the app.

### Populating the database
Connect to the database and add the table.

First, go into the container.
```
docker exec -it wallet_database /bin/sh
```

Connect to the database.
```
psql -U postgres
```

Add a table.
```Sql
CREATE TABLE IF NOT EXISTS account(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
currency VARCHAR(3), /* Currency follows ISO 4217 standard */
balance NUMERIC,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

constraint balance_non_negative check (balance >= 0.0)
);
```

Alternatively, we can use [golang-migrate](https://github.com/golang-migrate/migrate) to manage our database migrations. 

## Test Wallet API
In a separate terminal, for graphical user interface [Postman](https://www.postman.com/) can be used, or [curl](https://curl.se/) command-line tool.

### Add new user
```
curl -X POST '127.0.0.1:8000/wallet/' -d '{"name": "Alice"}'
```

### Add funds to a wallet for the user.
```
curl -X POST '127.0.0.1:8000/wallet/deposit/' -d '{"name": "Alice", "balance": 450}'
```

### Remove funds from a wallet for the user.
```
curl -X POST '127.0.0.1:8000/wallet/withdraw/' -d '{"name": "Alice", "balance": 150}'
```

### Query the current state of a wallet.
#### For a particular user.
```
curl -X GET '127.0.0.1:8000/wallet/name/' -d '{"name": "Alice"}'
```

#### For all the users.
```
curl -X GET '127.0.0.1:8000/wallet/'
```

## Stop the services
```
docker-compose stop
```

## Endpoints
All the routes that are exposed in the app.
```Go
router.HandleFunc("/wallet/", database.Fetch).Methods("GET")
router.HandleFunc("/wallet/name/", database.FetchByName).Methods("GET")
router.HandleFunc("/wallet/", database.Create).Methods("POST")
router.HandleFunc("/wallet/deposit/", database.Deposit).Methods("POST")
router.HandleFunc("/wallet/withdraw/", database.Withdraw).Methods("POST")
router.HandleFunc("/wallet/delete/", database.DeleteByName).Methods("DELETE")

```


## Useful resources
* https://go.dev/doc/tutorial/
* https://pkg.go.dev/std
* https://www.practical-go-lessons.com/
* https://pkg.go.dev/github.com/lib/pq
* https://pkg.go.dev/github.com/gorilla/mux
* https://www.sohamkamani.com/golang/json/
* https://mholt.github.io/json-to-go/
