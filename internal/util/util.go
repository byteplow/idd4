package util

import (
	"encoding/json"
	"html/template"
	"log"
	"net/url"

	kratos "github.com/ory/kratos-client-go"
)

func UiToNodes(uiContainer kratos.UiContainer) *map[string]interface{} {
	ui := map[string]interface{}{}

	str, _ := uiContainer.MarshalJSON()

	log.Println(string(str))

	err := json.Unmarshal(str, &ui)
	if err != nil {
		panic(err)
	}

	ui["groups"] = generateUiGroups(&ui)

	return &ui
}

func generateUiGroups(ui *map[string]interface{}) *map[string]map[string]interface{} {
	nodes, ok := (*ui)["nodes"].([]interface{})
	if !ok {
		panic(ok)
	}

	// group nodes by "group"
	grouped := map[string][]interface{}{}
	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		group, ok := nodeMap["group"].(string)
		if !ok {
			continue
		}

		grouped[group] = append(grouped[group], node)
	}

	// copy nodes from "default" to all groups
	for key := range grouped {
		for _, node := range grouped["default"] {
			grouped[key] = append(grouped[key], node)
		}
	}

	// delete default group
	delete(grouped, "default")

	// create ui group object
	groups := map[string]map[string]interface{}{}
	for key, group := range grouped {
		groups[key] = map[string]interface{}{
			"nodes":  group,
			"title":  key,
			"method": (*ui)["method"],
			"action": (*ui)["action"],
		}
	}

	return &groups
}

func UrlWithReturnTo(originalUrl string, to string) string {
	return UriWithQuery("return_to", originalUrl, to)
}

func UriWithQuery(key string, originalUrl string, value string) string {
	newUrl, err := url.Parse(originalUrl)
	if err != nil {
		panic(err)
	}

	q := newUrl.Query()
	q.Add(key, value)

	newUrl.RawQuery = q.Encode()

	return newUrl.String()
}

func ToTemplateUrl(str string) template.URL {
	return template.URL(str)
}

func ToTemplateJs(str string) template.JS {
	return template.JS(str)
}
