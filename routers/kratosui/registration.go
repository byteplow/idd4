package kratosui

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/internal/invite"
	"github.com/byteplow/idd4/internal/util"
	"github.com/gin-gonic/gin"
	kratos "github.com/ory/kratos-client-go"
)

var Endpoint_registration = "./self-service/registration"

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func GetRegistration(c *gin.Context) {
	cookie := c.GetHeader("Cookie")
	id, ok := c.GetQuery("flow")
	if !ok {
		c.Redirect(http.StatusSeeOther, config.Config.Urls["registration_url"])
		return
	}

	flow, _, err := container.KratosPublicClient.GetSelfServiceRegistrationFlow(context.Background()).Cookie(cookie).Id(id).Execute()
	if err != nil {
		c.Redirect(http.StatusSeeOther, config.Config.Urls["registration_url"])
		return
	}

	u, err := url.Parse(flow.RequestUrl)
	if err != nil {
		panic(err)
	}

	invite := u.Query().Get("invite")
	errName, hasError := c.GetQuery("error")

	if invite == "" || hasError {
		errText := config.Config.Messages.NoInviteLinkErrorMessage
		if errName == "invalid_invite" {
			errText = config.Config.Messages.InvalidInviteLinkErrorMessage
		}

		flow.Ui.Messages = append(flow.Ui.Messages, kratos.UiText{
			Text: errText,
			Type: "error",
		})

		flow.Ui.Nodes = []kratos.UiNode{}
	}

	nodes := util.UiToNodes(flow.Ui)

	c.HTML(http.StatusOK, "registration.tmpl", gin.H{
		"ui":   nodes,
		"urls": config.Config.Urls,
	})
}

func PostRegistration(c *gin.Context) {
	cookie := c.GetHeader("Cookie")
	id, ok := c.GetQuery("flow")
	if !ok {
		panic(ok)
	}

	flow, _, err := container.KratosPublicClient.GetSelfServiceRegistrationFlow(context.Background()).Cookie(cookie).Id(id).Execute()
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(flow.RequestUrl)
	if err != nil {
		panic(err)
	}

	i := u.Query().Get("invite")
	if i == "" {
		c.Redirect(http.StatusSeeOther, c.GetHeader("referer"))
	}

	//todo: remove wellknown
	if i != "wellknown" && !invite.CheckInvite(i, Endpoint_registration) {
		c.Redirect(http.StatusSeeOther, util.UriWithQuery("error", c.GetHeader("referer"), "invalid_invite"))
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, vial []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req := c.Copy().Request
	req.RequestURI = ""

	deleteHopHeader(&req.Header)

	backend, err := url.Parse(util.UriWithQuery("flow", config.Config.Urls["registration_url_internal"], id))
	if err != nil {
		panic(err)
	}

	req.URL = backend
	req.Host = backend.Host

	log.Println(req)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	deleteHopHeader(&res.Header)

	for key, values := range res.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	_, err = io.Copy(c.Writer, res.Body)
	if err != nil {
		panic(err)
	}

	c.Status(res.StatusCode)

	//todo: remove wellknown
	if res.StatusCode == 303 && i != "wellknown" {
		invite.InvalidateInvite(i, "")
	}
}

func deleteHopHeader(header *http.Header) {
	for _, key := range hopHeaders {
		header.Del(key)
	}
}
