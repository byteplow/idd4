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
	loginUrl := util.UrlWithReturnTo(config.Config.Urls["login_url"], config.Config.Urls["settings_url"])

	cookie, err := c.Request.Cookie("ory_kratos_session")
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, loginUrl)
		return
	}

	session, _, err := container.KratosPublicClient.ToSession(context.Background()).Cookie(cookie.String()).Execute()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, loginUrl)
		return
	}

	if !*session.Active {
		c.Redirect(http.StatusSeeOther, loginUrl)
		return
	}

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
