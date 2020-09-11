package api

import (
	"github.com/darcys22/godbledger-web/pkg/middleware"
)

const accountListing = `[
			{value: 0, text: "Guest"}
			{value: 1, text: "Service"}
			{value: 2, text: "Customer"}
			{value: 3, text: "Operator"}
			{value: 4, text: "Support"}
			{value: 5, text: "Admin"}
			]`

func GetAccountListing(c *middleware.Context) {
	c.JSON(200, accountlisting)
}
