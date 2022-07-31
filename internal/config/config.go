package config

type Configuration struct {
	Urls     map[string]string
	Hydra    *Hydra
	Kratos   *Kratos
	Server   *Server
	Messages *Messages
}

type Messages struct {
	NoInviteLinkErrorMessage      string
	InvalidInviteLinkErrorMessage string
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
			"hydra_login_url":           "https://localhost:4000/flow/login",
			"welcome_url":               "https://localhost:4000/",
			"settings_url":              "https://localhost:4000/self-service/settings/browser",
			"login_url":                 "https://localhost:4000/self-service/login/browser",
			"invite_url":                "https://localhost:4000/invite",
			"registration_url_internal": "http://kratos:4433/self-service/registration",
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
		Messages: &Messages{
			NoInviteLinkErrorMessage:      "Registration without an invite link is forbidden.",
			InvalidInviteLinkErrorMessage: "Invite link is invalid.",
		},
	}

	return nil
}
