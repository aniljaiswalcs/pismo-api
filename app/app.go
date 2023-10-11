package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/aniljaiswalcs/pismo/handler"
	"github.com/aniljaiswalcs/pismo/repository/adapter"
	"github.com/gorilla/mux"
)

func Start() {

	db := getNewPullConnectionDb()
	defer db.Close()

	accountRepositoryPostgres := adapter.NewAccountRepositoryPostgres(db)
	transactionRepositoryPostgres := adapter.NewTransactionRepositoryPostgres(db)

	accountHandler := handler.NewAccountHandler(accountRepositoryPostgres)
	transactionHandler := handler.NewTransactionHandler(transactionRepositoryPostgres)

	//port := ":" + os.Getenv("API_PORT")
	port := ":3000"

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	// routes to accounts
	accountMux := router.PathPrefix("/accounts").Subrouter()
	accountMux.HandleFunc("", accountHandler.CreateAccount).Methods("POST")
	accountMux.HandleFunc("/{accountId:[0-9]+}", accountHandler.GetAccount).Methods("GET")

	// routes to transaction
	transactionMux := router.PathPrefix("/transactions").Subrouter()
	transactionMux.HandleFunc("", transactionHandler.CreateTransaction).Methods("POST")
	transactionMux.HandleFunc("/{transactionid:[0-9]+}", transactionHandler.GetAccount).Methods("GET")

	fmt.Println("Server: localhost" + port)

	http.Handle("/", router)
	fmt.Println(http.ListenAndServe(port, nil))
}

func getNewPullConnectionDb() *sql.DB {

	connStr := os.Getenv("POSTGRESQL_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
