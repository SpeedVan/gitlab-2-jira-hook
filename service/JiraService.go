package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SpeedVan/go-common/client/httpclient"
	"github.com/SpeedVan/go-common/config"
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
	Fields *Fields `json:"fields"`
}

type Fields struct {
	Project     map[string]string `json:"project"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Issuetype   map[string]string `json:"issuetype"`
	Epic        string            `json:"customfield_10006"`
	Reporter    map[string]string `json:"reporter"`
}

// CreateIssue response return
// {
//    "id": "176643",
//    "key": "DECISION-6",
//    "self": "http://<hostname>/rest/api/2/issue/176643"
// }
func (s *JiraService) CreateIssue(projectKey, parentIssue, issueTypeID, summary, description, remoteLink, reporter string) error {
	fullURL := fmt.Sprintf("http://<hostname>/rest/api/2/issue/")
	issue := &Issue{
		Fields: &Fields{
			Project: map[string]string{
				"key": projectKey,
			},
			Summary:     summary,
			Description: description,
			Issuetype: map[string]string{
				"id": issueTypeID,
			},
			Epic: parentIssue,
			Reporter: map[string]string{
				"name": reporter,
			},
		},
	}
	bs, _ := json.Marshal(issue)
	fmt.Println(string(bs))
	req, _ := http.NewRequest("POST", fullURL, bytes.NewReader(bs))
	req.SetBasicAuth("username", "password")
	req.Header.Set("Content-Type", "application/json")
	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	result := make(map[string]string)
	rbs, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(rbs))
	json.Unmarshal(rbs, &result)
	return s.CreateRemoteLink(result["id"], remoteLink, "Gitlab issue")
}

// CreateRemoteLink todo
func (s *JiraService) CreateRemoteLink(issueIDOrKey, url, title string) error {
	fullURL := fmt.Sprintf("http://<hostname>/rest/api/2/issue/%v/remotelink", issueIDOrKey)
	jsonMap := map[string]map[string]string{
		"object": map[string]string{
			"url":   url,
			"title": title,
		},
	}

	bs, _ := json.Marshal(jsonMap)
	req, _ := http.NewRequest("POST", fullURL, bytes.NewReader(bs))
	req.SetBasicAuth("username", "password")
	req.Header.Set("Content-Type", "application/json")
	_, err := s.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	return nil
}
