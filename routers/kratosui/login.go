package kratosui

import (
	"context"
	"net/http"
	"log"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/internal/util"
	"github.com/gin-gonic/gin"
)

func GetLogin(c *gin.Context) {
	cookie := c.GetHeader("Cookie")
	id, ok := c.GetQuery("flow")
	if !ok {
		c.Redirect(http.StatusSeeOther, config.Config.Urls["welcome_url"])
		return
	}

	flow, _, err := container.KratosPublicClient.GetSelfServiceLoginFlow(context.Background()).Cookie(cookie).Id(id).Execute()
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusSeeOther, config.Config.Urls["welcome_url"])
		return
	}

	nodes := util.UiToNodes(flow.Ui)
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"ui":   nodes,
		"urls": config.Config.Urls,
	})
}
