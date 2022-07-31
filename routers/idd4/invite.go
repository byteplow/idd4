package idd4

import (
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/invite"
	"github.com/byteplow/idd4/internal/util"
	"github.com/byteplow/idd4/routers/kratosui"
	"github.com/gin-gonic/gin"
	kratos "github.com/ory/kratos-client-go"
)

func PostInvite(c *gin.Context) {
	session := util.RequireActiveKratosSession(c)
	if session == nil {
		c.Redirect(http.StatusSeeOther, util.UrlWithReturnTo(config.Config.Urls["login_url"], config.Config.Urls["invite_url"]))
		return
	}

	_, ok := c.GetPostForm("create_invite")
	if ok {
		i, err := invite.CreateInvite(session.Identity.Id, kratosui.Endpoint_registration)
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusSeeOther, util.UriWithQuery("invite", config.Config.Urls["invite_url"], i))
	}

	i, ok := c.GetPostForm("remove_invite")
	if ok {
		err := invite.InvalidateInvite(i, "")
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusSeeOther, config.Config.Urls["invite_url"])
	}
}

func GetInvite(c *gin.Context) {
	session := util.RequireActiveKratosSession(c)
	if session == nil {
		c.Redirect(http.StatusSeeOther, util.UrlWithReturnTo(config.Config.Urls["login_url"], config.Config.Urls["invite_url"]))
		return
	}

	ui := kratos.UiContainer{
		Action: config.Config.Urls["invite_flow_url"],
		Method: "POST",
		Nodes:  []kratos.UiNode{},
	}

	invites, err := invite.ListInvites(session.Identity.Id)
	if err != nil {
		panic(err)
	}

	for _, invite := range invites {
		ui.Nodes = append(ui.Nodes, kratos.UiNode{
			Attributes: kratos.UiNodeAttributes{
				UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
					Type:     "submit",
					Value:    invite,
					Disabled: false,
					Name:     "remove_invite",
				},
			},
			Type:  "input",
			Group: "invite",
			Meta: kratos.UiNodeMeta{
				Label: &kratos.UiText{
					Text: "Remove " + invite,
				},
			},
		})
	}

	ui.Nodes = append(ui.Nodes, kratos.UiNode{
		Attributes: kratos.UiNodeAttributes{
			UiNodeInputAttributes: &kratos.UiNodeInputAttributes{
				Type:     "submit",
				Disabled: false,
				Name:     "create_invite",
				NodeType: "input",
				Value:    "create_invite",
			},
		},
		Type:  "input",
		Group: "invite",
		Meta: kratos.UiNodeMeta{
			Label: &kratos.UiText{
				Text: "Create Invite",
			},
		},
	})

	invite := c.Query("invite")
	if invite != "" {
		ui.Messages = []kratos.UiText{{
			Type: "info",
			Text: util.UriWithQuery("invite", config.Config.Urls["registration_url"], invite),
		}}
	}

	nodes := util.UiToNodes(ui)
	c.HTML(http.StatusOK, "invite.tmpl", gin.H{
		"ui":   nodes,
		"urls": config.Config.Urls,
	})
}
