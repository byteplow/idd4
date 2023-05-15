package hydraui

import (
	"context"
	"log"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/internal/util"

	"github.com/gin-gonic/gin"
	client "github.com/ory/hydra-client-go"
)

func GetLogin(c *gin.Context) {
	//crfs protection

	challenge, ok := c.GetQuery("login_challenge")
	if !ok {
		panic(ok)
	}

	loginRequest, _, err := container.HydraAdminClient.AdminApi.GetLoginRequest(context.Background()).LoginChallenge(challenge).Execute()
	if err != nil {
		log.Println(loginRequest)
		panic(err)
	}

	if loginRequest.Skip {
		loginAccept, _, err := container.HydraAdminClient.AdminApi.AcceptLoginRequest(context.Background()).LoginChallenge(challenge).AcceptLoginRequest(client.AcceptLoginRequest{
			Subject: loginRequest.Subject,
		}).Execute()
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusSeeOther, loginAccept.RedirectTo)
		return
	}

	loginUrl := util.UrlWithReturnTo(config.Config.Urls["login_url"], util.UriWithQuery("login_challenge", config.Config.Urls["hydra_login_url"], challenge))

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

	loginAccept, _, err := container.HydraAdminClient.AdminApi.AcceptLoginRequest(context.Background()).LoginChallenge(challenge).AcceptLoginRequest(client.AcceptLoginRequest{
		Subject: session.Identity.Id,
	}).Execute()
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusSeeOther, loginAccept.RedirectTo)
	return
}
