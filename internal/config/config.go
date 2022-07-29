package config

type Configuration struct {
	Urls   map[string]string
	Hydra  *Hydra
	Kratos *Kratos
	Server *Server
}

type Hydra struct {
	Session     *Session
	AdminApiUrl string
}

type Kratos struct {
	AdminApiUrl  string
	PublicApiUrl string
}

type Session struct {
	RememberFor int64
}

type Server struct {
	RunMode  string
	Endpoint string
}

var Config *Configuration

func Setup() error {
	Config = &Configuration{
		Urls: map[string]string{
			"hydra_login_url": "https://localhost:4000/flow/login",
			"welcome_url":     "https://localhost:4000/",
			"settings_url":    "https://localhost:4000/self-service/settings/browser",
			"login_url":       "https://localhost:4000/self-service/login/browser",
			"invite_url":      "https://localhost:4000/invite",
		},
		Hydra: &Hydra{
			Session: &Session{
				RememberFor: 3600,
			},
			AdminApiUrl: "http://hydra:4444",
		},
		Kratos: &Kratos{
			AdminApiUrl:  "http://kratos:4434",
			PublicApiUrl: "http://kratos:4433",
		},
		Server: &Server{
			RunMode:  "debug",
			Endpoint: ":4455",
		},
	}

	return nil
}
