package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"test/model"
	"time"
)

const (
	apiHost       = "en.wikipedia.org"
	apiProtocol   = "https"
	apiReqTimeout = 15
)

var notFoundErr = errors.New("article not found")

type Service interface {
	GetWikiShortDesc(query string) (string, error)
}

type service struct {
	client *http.Client
}

func New() Service {
	return service{&http.Client{
		Timeout: time.Duration(apiReqTimeout) * time.Second,
	}}
}
func (s service) GetWikiShortDesc(query string) (string, error) {
	var res model.Response
	actionURL := getQueryActionURL(query)
	body, err := s.get(actionURL)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return "", err
	}

	if res.Query.Pages[0].Missing { //article wasn't found
		return "", notFoundErr
	}
	return getShortDescFromResponse(res.Query.Pages[0].Revisions[0].Content)
}

func (s service) get(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest: failed to create request")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {

		return "", errors.New(getRespBodyAsString(resp))
	} else {
		return getRespBodyAsString(resp), nil
	}
}
func getQueryActionURL(query string) string {
	qs := url.Values{}

	qs.Set("titles", query)
	qs.Set("prop", "revisions")
	qs.Set("action", "query")
	qs.Set("rvlimit", "1")
	qs.Set("formatversion", "2")
	qs.Set("rvprop", "content")
	qs.Set("format", "json")

	url := url.URL{
		Host:     apiHost,
		Scheme:   apiProtocol,
		Path:     path.Join("w", "api.php"),
		RawQuery: qs.Encode(),
	}

	return url.String()
}

func getRespBodyAsString(resp *http.Response) string {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR Reading HTTP Response Body (error)"
	}
	return string(bodyBytes)
}

func getShortDescFromResponse(content string) (string, error) {
	r, err := regexp.Compile("((?i)Short description)+(\\|)([\\w\\s\\d\\.×]*)")
	if err != nil {
		return "", errors.Wrap(err, "regexp error")
	}

	//the short-desc should be in the 3rd capture group
	return r.FindAllStringSubmatch(content, -1)[0][3], nil
}
