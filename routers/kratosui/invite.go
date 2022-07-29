package kratosui

import (
	"context"
	"log"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/internal/kratos"
	"github.com/byteplow/idd4/internal/util"
	"github.com/gin-gonic/gin"
)

func GetInvite(c *gin.Context) {
	loginUrl := util.UrlWithReturnTo(config.Config.Urls["login_url"], config.Config.Urls["invite_url"])

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

	invite, err := kratos.CreateInvite()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "invite.tmpl", gin.H{
		"invite": invite,
		"urls":   config.Config.Urls,
	})
}
