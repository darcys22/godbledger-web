package middleware

import (
	//"strconv"
	"strings"

	"github.com/darcys22/godbledger-web/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "middleware")

type Context struct {
	*gin.Context
	//*m.SignedInUser

	//Session session.Store

	IsSignedIn     bool
	AllowAnonymous bool
}

//func GetContextHandler() gin.HandlerFunc {
//return func(c *Context) {
//ctx := &Context{
//Context: c,
////Session:        sess,
////SignedInUser:   &m.SignedInUser{},
//IsSignedIn:     true,
//AllowAnonymous: true,
//}

//c.Map(ctx)
//}
//}

// Handle handles and logs error by given status.
func (ctx *Context) Handle(status int, title string, err error) {
	//if err != nil {
	//log.Error(4, "%s: %v", title, err)
	//if setting.Env != setting.PROD {
	//ctx.Data["ErrorMsg"] = err
	//}
	//}

	//ctx.Data["Title"] = title
	ctx.HTML(status, "https", gin.H{
		"status": "success",
	})
}

func (ctx *Context) JsonOK(message string) {
	resp := make(map[string]interface{})

	resp["message"] = message

	ctx.JSON(200, resp)
}

func (ctx *Context) IsApiRequest() bool {
	return strings.HasPrefix(ctx.Request.URL.Path, "/api")
}

func (ctx *Context) JsonApiErr(status int, message string, err error) {
	resp := make(map[string]interface{})

	if err != nil {
		log.Error(4, "%s: %v", message, err)
		if setting.Env != setting.PROD {
			resp["error"] = err.Error()
		}
	}

	switch status {
	case 404:
		resp["message"] = "Not Found"
	case 500:
		resp["message"] = "Internal Server Error"
	}

	if message != "" {
		resp["message"] = message
	}

	ctx.JSON(status, resp)
}
