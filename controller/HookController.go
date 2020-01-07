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
	GitlabProject2JiraProjectMap map[string]string
	GitlabLabel2JiraIssueTypeMap map[string]string
}

func New(cfg config.Config) *HookController {
	return &HookController{
		JiraService: service.New(cfg),
		GitlabProject2JiraProjectMap: map[string]string{
			"finup-decision-saas":     "DECISION-2",
			"finup-decision-saas/web": "DECISION-4",
		},
		GitlabLabel2JiraIssueTypeMap: map[string]string{
			"BUG": "故障",
		},x
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
		projectKey := ""
		if str, ok := s.GitlabProject2JiraProjectMap[pathWithNamespace]; ok {
			projectKey = str
		} else if str2, ok := s.GitlabProject2JiraProjectMap[namespace]; ok {
			projectKey = str2
		} else {
			return
		}
		url := jsonMap.ObjectAttributes["url"].(string)
		title := jsonMap.ObjectAttributes["title"].(string)
		description := jsonMap.ObjectAttributes["description"].(string)

		labels := jsonMap.Labels
		issueType := ""
		for _, item := range labels {
			issueType = s.GitlabLabel2JiraIssueTypeMap[item["title"].(string)]
		}
		if issueType == "" {
			return
		}
		s.JiraService.CreateIssue(projectKey, issueType, title, description, url)
	}
}

type HookIssue struct {
	ObjectAttributes map[string]interface{}   `json:"object_attributes"`
	Project          map[string]interface{}   `json:"project"`
	Labels           []map[string]interface{} `json:"labels"`
}
