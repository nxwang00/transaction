package router

import (
	"net/http"

	handler "github.com/server/transaction/handlers"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

// For domain
var myRoutes = Routes{
	Route{
		"CreateTransaction",
		"POST",
		"/trans",
		handler.CreateTransaction,
	},
	Route{
		"ReadTransactions",
		"GET",
		"/trans",
		handler.ReadTransactions,
	},
	// Route{
	// 	"ReadDomain",
	// 	"GET",
	// 	"/domains/{id}",
	// 	handler.ReadDomain,
	// },
}