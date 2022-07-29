package hydraui

import (
	"context"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/gin-gonic/gin"
	client "github.com/ory/hydra-client-go"
)

func GetConsent(c *gin.Context) {
	challenge, ok := c.GetQuery("consent_challenge")
	if !ok {
		c.Data(http.StatusBadRequest, "text/plain", []byte("consent_challenge missing"))
		return
	}

	consentRequest, _, err := container.HydraAdminClient.AdminApi.GetConsentRequest(context.Background()).ConsentChallenge(challenge).Execute()
	if err != nil {
		panic(err)
	}

	if *consentRequest.Skip {
		consentAccept, _, err := container.HydraAdminClient.AdminApi.AcceptConsentRequest(context.Background()).ConsentChallenge(challenge).AcceptConsentRequest(client.AcceptConsentRequest{
			GrantScope:               consentRequest.RequestedScope,
			GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
		}).Execute()
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusSeeOther, consentAccept.RedirectTo)
		return
	}

	remember := true
	rememberFor := config.Config.Hydra.Session.RememberFor

	consentAccept, _, err := container.HydraAdminClient.AdminApi.AcceptConsentRequest(context.Background()).ConsentChallenge(challenge).AcceptConsentRequest(client.AcceptConsentRequest{
		GrantScope:               consentRequest.RequestedScope,
		GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
		Remember:                 &remember,
		RememberFor:              &rememberFor,
		Session:                  client.NewConsentRequestSession(),
	}).Execute()
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusSeeOther, consentAccept.RedirectTo)
	return
}
