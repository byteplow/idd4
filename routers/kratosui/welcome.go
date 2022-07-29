package kratosui

import (
	"context"
	"log"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/gin-gonic/gin"
	client "github.com/ory/kratos-client-go"
)

func GetWelcome(c *gin.Context) {
	logoutUrl := ""
	var session *client.Session
	signedIn := false

	//try to get session cookie
	cookie, err := c.Request.Cookie("ory_kratos_session")
	if err == nil {
		//try to get session
		session, _, err = container.KratosPublicClient.ToSession(context.Background()).Cookie(cookie.String()).Execute()
		if err != nil {
			//we can handle an error
			log.Println(err)
		} else {
			if *session.Active {
				signedIn = true
				logout, _, err := container.KratosPublicClient.CreateSelfServiceLogoutFlowUrlForBrowsers(context.Background()).Cookie(cookie.String()).Execute()
				if err != nil {
					log.Println(err)
					c.Redirect(http.StatusSeeOther, config.Config.Urls["login_url"])
					return
				}

				logoutUrl = logout.LogoutUrl
			}
		}
	}

	welcomeMessage := "Welcome"
	if signedIn && session != nil {
		username := session.Identity.Traits.(map[string]interface{})["username"]
		if username != nil {
			welcomeMessage = "Welcome " + username.(string)
		}
	}

	c.HTML(http.StatusOK, "welcome.tmpl", gin.H{
		"welcomeMessage": welcomeMessage,
		"signedIn":       signedIn,
		"logout_url":     logoutUrl,
		"urls":           config.Config.Urls,
	})
}
