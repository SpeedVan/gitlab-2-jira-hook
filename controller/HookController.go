package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedVan/gitlab-2-jira-hook/service"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config"
)

// HookController todo
type HookController struct {
	web.Controller
	JiraService                  *service.JiraService
	GitlabProject2JiraEpicMap    map[string]string
	GitlabLabel2JiraIssueTypeMap map[string]string
}

func New(cfg config.Config) *HookController {
	return &HookController{
		JiraService: service.New(cfg),
		GitlabProject2JiraEpicMap: map[string]string{
			"finup-decision-saas":     "DECISION-1",
			"finup-decision-saas/web": "DECISION-2",
		},
		GitlabLabel2JiraIssueTypeMap: map[string]string{
			"jira:bug":   "11001",
			"jira:task":  "10000",
			"jira:story": "11000",
		},
	}
}

// GetRoute todo
func (s *HookController) GetRoute() web.RouteMap {
	items := []*web.RouteItem{
		&web.RouteItem{Path: "/issue/hook", HandleFunc: s.Hook},
	}

	return web.NewRouteMap(items...)
}

// Hook todo
func (s *HookController) Hook(w http.ResponseWriter, r *http.Request) {
	bs, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(bs))
	jsonMap := &HookIssue{}
	json.Unmarshal(bs, jsonMap)

	action := jsonMap.ObjectAttributes["action"].(string)

	if action == "open" {
		namespace := jsonMap.Project["namespace"].(string)
		pathWithNamespace := jsonMap.Project["path_with_namespace"].(string)
		epic := ""
		if str, ok := s.GitlabProject2JiraEpicMap[pathWithNamespace]; ok {
			epic = str
		} else if str2, ok := s.GitlabProject2JiraEpicMap[namespace]; ok {
			epic = str2
		} else {
			fmt.Println("not currect epic")
			return
		}
		fmt.Println("epic:" + epic)
		url := jsonMap.ObjectAttributes["url"].(string)
		title := jsonMap.ObjectAttributes["title"].(string)
		description := jsonMap.ObjectAttributes["description"].(string)

		labels := jsonMap.Labels
		issueType := ""
		for _, item := range labels {
			issueType = s.GitlabLabel2JiraIssueTypeMap[item["title"].(string)]
		}
		if issueType == "" {
			fmt.Println("not currect issueType")
			return
		}
		fmt.Println("issueType:" + issueType)

		reporter := jsonMap.User["username"]
		if len(jsonMap.Assignees) > 0 {
			reporter = jsonMap.Assignees[0]["username"]
		}
		err := s.JiraService.CreateIssue("DECISION", epic, issueType, title, description, url, reporter)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

type HookIssue struct {
	ObjectAttributes map[string]interface{}   `json:"object_attributes"`
	Project          map[string]interface{}   `json:"project"`
	Labels           []map[string]interface{} `json:"labels"`
	// {
	// 	"name": "alex",
	// 	"username": "username",
	// 	"avatar_url": "http://www.gravatar.com/avatar/f9d154aeef64cf3c3d0150ec31096ddb?s=80&d=identicon"
	// }
	Assignees []map[string]string `json:"assignees"`
	User      map[string]string   `json:"user"`
}
