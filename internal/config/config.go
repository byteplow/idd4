package config

import (
	"errors"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Urls     map[string]string `json:"urls,omitempty"`
	Hydra    *Hydra            `json:"hydra,omitempty"`
	Kratos   *Kratos           `json:"kratos,omitempty"`
	Keto     *Keto             `json:"keto,omitempty"`
	Server   *Server           `json:"server,omitempty"`
	Messages *Messages         `json:"messages,omitempty"`
}

type Messages struct {
	NoInviteLinkErrorMessage      string `json:",omitempty"`
	InvalidInviteLinkErrorMessage string `json:",omitempty"`
}

type Hydra struct {
	Session     *Session `json:"session,omitempty"`
	AdminApiUrl string   `json:"admin_api_url,omitempty"`
}

type Kratos struct {
	AdminApiUrl  string `json:"admin_api_url,omitempty"`
	PublicApiUrl string `json:"public_api_url,omitempty"`
}

type Session struct {
	RememberFor int64 `json:"remember_for,omitempty"`
}

type Server struct {
	RunMode  string `json:"run_mode,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type Keto struct {
	WriteApiUrl string `json:"write_api_url,omitempty"`
	ReadApiUrl  string `json:"read_api_url,omitempty"`
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
			AdminApiUrl: "http://hydra:4444",
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

	return nil
}

func loadConfig(r io.Reader) (*Configuration, error) {
	c := defaultConfig()
	if err := yaml.NewDecoder(r).Decode(c); err != nil {
		log.Println("file Decode")
		return nil, err
	}

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
	file, err := os.Open(path)
	if err != nil {
		log.Println("file open failed")
		return err
	}
	defer file.Close()

	backup := Config

	c, err := loadConfig(file)
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
