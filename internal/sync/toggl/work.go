package toggl

import (
	"fmt"
	"io"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/auth"
)

// TODO make some of this options configurable (toml, env or yaml)
const (
	ZeroDay     = "2020-02-01"
	baseUrl     = "https://api.track.toggl.com/"
	reportsPath = "reports/api/v2/"
	summary     = "summary"
	user_agent  = "cosmedaniel8@gmail.com"
	wid         = "3338214"
	project_ids = "160884301,159981833,169021695"
	// since yyyy-mm-dd
	// until yyyy-mm-dd
)

var ErrNoRecord = fmt.Errorf("Error procesing the logs range result")

type WorkCollector struct {
	Auth  *auth.AuthConfig
	Scope string
}

func (wc *WorkCollector) LogsRange(start, end string) (map[string][]byte, error) {
	return nil, nil
}

func (wc *WorkCollector) DayLog(date string) ([]byte, error) {
	url := makeURL(date)
	req, err := wc.createRequest(url)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	return body, nil
}

func makeURL(date string) string {
	url := fmt.Sprintf("%s%s%s", baseUrl, reportsPath, summary)
	urlParams := map[string]string{
		"user_agent":   user_agent,
		"workspace_id": wid,
		"project_ids":  project_ids,
		"since":        date,
		"until":        date,
	}
	urlEncoded := url + "?" + auth.UrlEncode(urlParams)
	return urlEncoded
}

func (wc *WorkCollector) createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	e := auth.EncodeCredentials(wc.Auth.ClientID, wc.Auth.ClientSecret)
	req.Header.Add("Authorization", "Basic "+e)
	return req, err
}
