package kratosui

import (
	"context"
	"log"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/internal/util"
	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	if session := util.RequireActiveKratosSession(c); session == nil {
		c.Redirect(http.StatusSeeOther, util.UrlWithReturnTo(config.Config.Urls["login_url"], config.Config.Urls["settings_url"]))
		return
	}

	cookie, _ := c.Request.Cookie("ory_kratos_session")

	id, ok := c.GetQuery("flow")
	if !ok {
		c.Redirect(http.StatusSeeOther, config.Config.Urls["settings_url"])
		return
	}

	flow, _, err := container.KratosPublicClient.GetSelfServiceSettingsFlow(context.Background()).Cookie(cookie.String()).Id(id).Execute()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, config.Config.Urls["settings_url"])
		return
	}

	nodes := util.UiToNodes(flow.Ui)

	c.HTML(http.StatusOK, "settings.tmpl", gin.H{
		"ui":   nodes,
		"urls": config.Config.Urls,
	})
}
