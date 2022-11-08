module github.com/server/transaction

go 1.19

require (
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
)

require github.com/lib/pq v1.10.7 // indirect

replace github.com/server/transaction/router => ./router

replace github.com/server/transaction/handlers => ./handlers

replace github.com/server/transaction/model => ./model

replace github.com/server/transaction/db => ./db
