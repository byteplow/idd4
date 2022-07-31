package routers

import (
	"html/template"
	"net/http"

	"github.com/byteplow/idd4/internal/util"
	"github.com/byteplow/idd4/routers/hydraui"
	"github.com/byteplow/idd4/routers/idd4"
	"github.com/byteplow/idd4/routers/kratosui"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.SetFuncMap(template.FuncMap{
		"toTemplateUrl": util.ToTemplateUrl,
		"toTemplateJs":  util.ToTemplateJs,
	})

	r.LoadHTMLGlob("./templates/*")

	r.StaticFS("/static", http.Dir("./static"))

	r.GET("/login", kratosui.GetLogin)
	r.GET("/error", kratosui.GetError)
	r.GET("/registration", kratosui.GetRegistration)
	r.GET("/welcome", kratosui.GetWelcome)
	r.GET("/settings", kratosui.GetSettings)
	r.GET("/", kratosui.GetWelcome)
	r.GET("/invite", idd4.GetInvite)
	r.POST("/flow/invite", idd4.PostInvite)
	r.GET("/consent", hydraui.GetConsent)
	r.GET("/flow/login", hydraui.GetLogin)
	r.POST("/self-service/registration", kratosui.PostRegistration)

	return r
}
