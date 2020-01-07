package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedVan/gitlab-2-jira-hook/service"
	"github.com/SpeedVan/go-common/app/web"
)

// HookController todo
type HookController struct {
	web.Controller
	JiraService *service.JiraService
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
	// jsonMap := make(map[string]interface{})
	// json.Unmarshal(bs, &jsonMap)
}
