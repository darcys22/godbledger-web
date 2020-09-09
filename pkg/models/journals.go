package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrDashboardNotFound = errors.New("Account not found")
)

// Dashboard model
type Dashboard struct {
	Id      int64
	Slug    string
	OrgId   int64
	Version int

	Created time.Time
	Updated time.Time

	Title string
	Data  map[string]interface{}
}

// NewDashboard creates a new dashboard
func NewDashboard(title string) *Dashboard {
	dash := &Dashboard{}
	dash.Data = make(map[string]interface{})
	dash.Data["title"] = title
	dash.Title = title
	return dash
}

// GetTags turns the tags in data json into go string array
func (dash *Dashboard) GetTags() []string {
	jsonTags := dash.Data["tags"]
	if jsonTags == nil {
		return []string{}
	}

	arr := jsonTags.([]interface{})
	b := make([]string, len(arr))
	for i := range arr {
		b[i] = arr[i].(string)
	}
	return b
}

// GetDashboardModel turns the command into the savable model
//func (cmd *SaveDashboardCommand) GetDashboardModel() *Dashboard {
//dash := &Dashboard{}
//dash.Data = cmd.Dashboard
//dash.Title = dash.Data["title"].(string)
//dash.OrgId = cmd.OrgId

//if dash.Data["id"] != nil {
//dash.Id = int64(dash.Data["id"].(float64))

//if dash.Data["version"] != nil {
//dash.Version = int(dash.Data["version"].(float64))
//}
//} else {
//dash.Data["version"] = 0
//}

//return dash
//}

//
// COMMANDS
//

type PostJournalCommand struct {
	Dashboard map[string]interface{} `json:"dashboard" binding:"Required"`
	Overwrite bool                   `json:"overwrite"`
	OrgId     int64                  `json:"-"`

	Result *Dashboard
}

//
// QUERIES
//

type GetDashboardQuery struct {
	Slug  string
	OrgId int64

	Result *Dashboard
}
