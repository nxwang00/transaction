module github.com/test/transaction

go 1.19

require (
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f
	github.com/gorilla/mux v1.8.0
)

replace github.com/test/transaction/router => ./router

replace github.com/test/transaction/handlers => ./handlers

replace github.com/test/transaction/model => ./model

replace github.com/test/transaction/db => ./db
