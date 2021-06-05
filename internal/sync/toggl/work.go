package toggl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/auth"
)

// TODO make some of this options configurable (toml, env or yaml)
const (
	ZeroDay     = "2021-01-01"
	baseUrl     = "https://api.track.toggl.com/"
	reportsPath = "reports/api/v2/"
	summary     = "summary"
	user_agent  = "cosmedaniel8@gmail.com"
	wid         = "3338214"
	project_ids = "160884301,159981833,169021695,170288296"
	// since yyyy-mm-dd
	// until yyyy-mm-dd
)

var ErrNoRecord = fmt.Errorf("Error procesing the logs range result")

type WorkProvider struct {
	Auth  *auth.AuthConfig
	Scope string
}

func (wc *WorkProvider) LogsRange(start, end string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	dateLayout := "2006-01-02"
	firstDate, err := time.Parse(dateLayout, start)
	if err != nil {
		return nil, err
	}
	lastDate, err := time.Parse(dateLayout, end)
	if err != nil {
		return nil, err
	}

	maxYear, maxMonth, maxDay := lastDate.Date()
	current := firstDate
	for {
		strDate := current.Format(dateLayout)
		currYear, currMonth, currDay := current.Date()

		log.Println("Getting Work log for", strDate)
		response, err := wc.DayLog(strDate)
		if err != nil {
			log.Println("ERR", err.Error())
			return nil, err
		}

		result[strDate] = response

		if currYear == maxYear && currMonth == maxMonth && currDay == maxDay {
			break
		}

		current = current.AddDate(0, 0, 1)
		log.Println("Sleeping...")
		time.Sleep(time.Second)
	}

	return result, nil
}

func (wc *WorkProvider) DayLog(date string) ([]byte, error) {
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

func (wc *WorkProvider) createRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	e := auth.EncodeCredentials(wc.Auth.ClientID, wc.Auth.ClientSecret)
	req.Header.Add("Authorization", "Basic "+e)
	return req, err
}
