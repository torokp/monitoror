package repository

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/monitorables/http/config"
)

type (
	httpRepository struct {
		httpClient *http.Client
	}
)

func NewHTTPRepository(config *config.HTTP) api.Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	return &httpRepository{client}
}

func (r *httpRepository) Get(url string) (response *models.Response, err error) {
	req, reqerr := http.NewRequest("GET", url, nil)
	if reqerr != nil {
		return
	}

	req.Header.Set("User-Agent", "monitoror/"+version.Version)

	// resp, err := r.httpClient.Get(url)
	resp, err := r.httpClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response = &models.Response{
		StatusCode: resp.StatusCode,
		Body:       bytes,
	}

	return
}
