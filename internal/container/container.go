package container

import (
	"github.com/byteplow/idd4/internal/config"
	hydra "github.com/ory/hydra-client-go"

	kratos "github.com/ory/kratos-client-go"
)

var HydraAdminClient *hydra.APIClient
var KratosPublicClient kratos.V0alpha2Api
var KratosAdminClient kratos.V0alpha2Api

func Setup() error {
	hydraAdminConf := hydra.NewConfiguration()
	hydraAdminConf.Servers = []hydra.ServerConfiguration{
		{
			URL: config.Config.Hydra.AdminApiUrl,
		},
	}

	HydraAdminClient = hydra.NewAPIClient(hydraAdminConf)

	kratosPublicConf := kratos.NewConfiguration()
	kratosPublicConf.Servers = []kratos.ServerConfiguration{
		{
			URL: config.Config.Kratos.PublicApiUrl,
		},
	}

	KratosPublicClient = kratos.NewAPIClient(kratosPublicConf).V0alpha2Api

	kratosAdminConf := kratos.NewConfiguration()
	kratosAdminConf.Servers = []kratos.ServerConfiguration{
		{
			URL: config.Config.Kratos.AdminApiUrl,
		},
	}

	KratosAdminClient = kratos.NewAPIClient(kratosAdminConf).V0alpha2Api

	return nil
}
