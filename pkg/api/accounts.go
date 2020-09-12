package api

import (
	"encoding/json"
	"fmt"
	"github.com/darcys22/godbledger-web/pkg/middleware"
)

type Results struct {
	Results []Account `json:"results"`
}

type Account struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

const accountListing = `{"results": [{"id": 0, "text": "Guest"},{"id": 1, "text": "Service"}]}`

func GetAccountListing(c *middleware.Context) {
	arr := Results{}
	err := json.Unmarshal([]byte(accountListing), &arr)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, &arr)
}