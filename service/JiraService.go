package service

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/SpeedVan/go-common/config"
	"github.com/SpeedVan/go-common/client/httpclient"
)

// JiraService todo
type JiraService struct {
	// RestAPIPrefix string
	HTTPClient *http.Client
}

// New todo
func New(cfg config.Config) *JiraService {
	// prefix := ""

	httpClient, _ := httpclient.New(cfg)
	return &JiraService{
		// RestAPIPrefix: prefix,
		HTTPClient: httpClient,
	}
}

type Issue struct {
	fields *Fields
}

type Fields struct {
	project     map[string]string
	summary     string
	description string
	issuetype   map[string]string
}

// CreateIssue response return
// {
//    "id": "176643",
//    "key": "DECISION-6",
//    "self": "http://jira.renmaitech.com/rest/api/2/issue/176643"
// }
func (s *JiraService) CreateIssue(projectKey, issueType, summary, description, remoteLink string) error {
	fullURL := fmt.Sprintf("http://jira.renmaitech.com/rest/api/2/issue/")
	issue := &Issue{
		fields: &Fields{
			project: map[string]string{
				"key": projectKey,
			},
			summary:     summary,
			description: description,
			issuetype: map[string]string{
				"name": issueType,
			},
		},
	}
	bs, _ := json.Marshal(issue)
	req, _ := http.NewRequest("POST", fullURL, bytes.NewReader(bs))
	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	result := make(map[string]string)
	rbs, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(rbs, &result)
	return s.CreateRemoteLink(result["id"], remoteLink, "任买Gitlab issue")
}

// CreateRemoteLink todo
func (s *JiraService) CreateRemoteLink(issueIDOrKey, url, title string) error {
	fullURL := fmt.Sprintf("http://jira.renmaitech.com/rest/api/2/issue/%v/remotelink", issueIDOrKey)
	jsonMap := map[string]map[string]string{
		"object": map[string]string{
			"url":   url,
			"title": title,
		},
	}

	bs, _ := json.Marshal(jsonMap)
	req, _ := http.NewRequest("POST", fullURL, bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	_, err := s.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	return nil
}
