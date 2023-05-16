package config

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Urls     map[string]string `yaml:"urls,omitempty"`
	Hydra    *Hydra            `yaml:"hydra,omitempty"`
	Kratos   *Kratos           `yaml:"kratos,omitempty"`
	Keto     *Keto             `yaml:"keto,omitempty"`
	Server   *Server           `yaml:"server,omitempty"`
	Messages *Messages         `yaml:"messages,omitempty"`
	MasterInvite string
}

type Messages struct {
	NoInviteLinkErrorMessage      string `yaml:",omitempty"`
	InvalidInviteLinkErrorMessage string `yaml:",omitempty"`
}

type Hydra struct {
	Session     *Session `yaml:"session,omitempty"`
	AdminApiUrl string   `yaml:"admin_api_url,omitempty"`
}

type Kratos struct {
	AdminApiUrl  string `yaml:"admin_api_url"`
	PublicApiUrl string `yaml:"public_api_url"`
}

type Session struct {
	RememberFor int64 `yaml:"remember_for,omitempty"`
}

type Server struct {
	RunMode  string `yaml:"run_mode,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
}

type Keto struct {
	WriteApiUrl string `yaml:"write_api_url,omitempty"`
	ReadApiUrl  string `yaml:"read_api_url,omitempty"`
}

var Config *Configuration
var requiredUrls = []string{
	"hydra_login_url",
	"hydra_login_url",
	"welcome_url",
	"settings_url",
	"login_url",
	"invite_url",
	"registration_url_internal",
	"registration_url",
}
var listeners = []func() error{}

func defaultConfig() *Configuration {
	return &Configuration{
		Urls: map[string]string{
			"hydra_login_url":           "https://localhost:4000/flow/login",
			"welcome_url":               "https://localhost:4000/",
			"settings_url":              "https://localhost:4000/self-service/settings/browser",
			"login_url":                 "https://localhost:4000/self-service/login/browser",
			"invite_url":                "https://localhost:4000/invite",
			"registration_url_internal": "http://kratos:4433/self-service/registration",
			"registration_url":          "https://localhost:4000/self-service/registration/browser",
		},
		Hydra: &Hydra{
			Session: &Session{
				RememberFor: 3600,
			},
			AdminApiUrl: "http://hydra:4445",
		},
		Kratos: &Kratos{
			AdminApiUrl:  "http://kratos:4434",
			PublicApiUrl: "http://kratos:4433",
		},
		Keto: &Keto{
			ReadApiUrl:  "http://keto:4466",
			WriteApiUrl: "http://keto:4467",
		},
		Server: &Server{
			RunMode:  "debug",
			Endpoint: ":4455",
		},
		Messages: &Messages{
			NoInviteLinkErrorMessage:      "Registration without an invite link is forbidden.",
			InvalidInviteLinkErrorMessage: "Invite link is invalid.",
		},
		MasterInvite: "",
	}
}

func checkConfig(c *Configuration) error {
	for _, u := range requiredUrls {
		val := c.Urls[u]
		if _, err := url.ParseRequestURI(val); err != nil {
			return errors.New("Urls." + u + " is invalid.")
		}
	}

	if _, err := url.ParseRequestURI(c.Hydra.AdminApiUrl); err != nil {
		return errors.New("Hydra.AdminApiUrl is invalid.")
	}

	if _, err := url.ParseRequestURI(c.Kratos.AdminApiUrl); err != nil {
		return errors.New("Kratos.AdminApiUrl is invalid.")
	}

	if _, err := url.ParseRequestURI(c.Kratos.AdminApiUrl); err != nil {
		return errors.New("Kratos.PublicApiUrl is invalid.")
	}

	if _, err := url.ParseRequestURI(c.Keto.ReadApiUrl); err != nil {
		return errors.New("Keto.ReadApiUrl is invalid.")
	}

	if _, err := url.ParseRequestURI(c.Keto.WriteApiUrl); err != nil {
		return errors.New("Keto.PublicApiUrl is invalid.")
	}

	if c.MasterInvite == "" {
		return errors.New("MasterInvite must not be empty.")
	}

	return nil
}

func loadConfig(path string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := defaultConfig()

	if len(bytes) > 0 {
		if err := yaml.Unmarshal(bytes, c); err != nil {
			log.Println("file Decode")
			return nil, err
		}
	}

	c.MasterInvite = os.Getenv("MASTER_INVITE")

	if err := checkConfig(c); err != nil {
		return nil, err
	}

	return c, nil
}

func Listen(listener func() error) {
	listeners = append(listeners, listener)
}

func notify() error {
	for _, listener := range listeners {
		err := listener()
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadConfig(path string) error {
	log.Println("Loading config")

	backup := Config

	c, err := loadConfig(path)
	if err != nil {
		return err
	}

	Config = c
	if err := notify(); err != nil {
		log.Println(err)
		log.Println("rolling back config")

		Config = backup

		if err := notify(); err != nil {
			log.Fatalln("rollback failed")
		}
	}

	return nil
}

func WatchConfig(path string) error {
	if err := LoadConfig(path); err != nil {
		log.Panicln(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file has changed")

					if err := LoadConfig(path); err != nil {
						log.Println(err)
						log.Println("Reloading config failed")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
