package kratosui

import (
	"context"
	"net/http"

	"github.com/byteplow/idd4/internal/container"
	"github.com/gin-gonic/gin"
)

func GetError(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		panic(ok)
	}

	flow, _, err := container.KratosPublicClient.GetSelfServiceError(context.Background()).Id(id).Execute()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "error.tmpl", flow.GetError())
}
