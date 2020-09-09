package api

import (
	"github.com/darcys22/godbledger-web/pkg/middleware"
)

func setIndexViewData(c *middleware.Context) error {
	//settings, err := getFrontendSettingsMap(c)
	//if err != nil {
	//return err
	return nil
}

//currentUser := &dtos.CurrentUser{
//IsSignedIn:     c.IsSignedIn,
//Login:          c.Login,
//Email:          c.Email,
//Name:           c.Name,
//LightTheme:     c.Theme == "light",
//OrgName:        c.OrgName,
//OrgRole:        c.OrgRole,
//GravatarUrl:    dtos.GetGravatarUrl(c.Email),
//IsGrafanaAdmin: c.IsGrafanaAdmin,
//}

//if len(currentUser.Name) == 0 {
//currentUser.Name = currentUser.Login
//}

//c.Data["User"] = currentUser
//c.Data["Settings"] = settings
//c.Data["AppUrl"] = setting.AppUrl
//c.Data["AppSubUrl"] = setting.AppSubUrl

//if setting.GoogleAnalyticsId != "" {
//c.Data["GoogleAnalyticsId"] = setting.GoogleAnalyticsId
//}

//return nil
//}

func Index(c *middleware.Context) {
	//if err := setIndexViewData(c); err != nil {
	//c.Handle(500, "Failed to get settings", err)
	//return
	//}

	c.HTML(200, "index")
}

func NotFound(c *middleware.Context) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	if err := setIndexViewData(c); err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}

	c.HTML(404, "index")
}
