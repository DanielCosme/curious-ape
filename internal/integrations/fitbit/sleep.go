package fitbit

import (
	"net/http"
	"time"
)

type SleepService struct {
	client Client
}

type SleepEnvelope struct {
	Sleep   []Sleep `json:"sleep"`
	Summary Summary `json:"summary"`
}
type Summary struct {
	Deep  Deep  `json:"deep"`
	Light Light `json:"light"`
	Rem   Rem   `json:"rem"`
	Wake  Wake  `json:"wake"`
}
type Sleep struct {
	DateOfSleep         string `json:"dateOfSleep"`
	Duration            int    `json:"duration"`
	Efficiency          int    `json:"efficiency"`
	EndTime             string `json:"endTime"`
	InfoCode            int    `json:"infoCode"`
	IsMainSleep         bool   `json:"isMainSleep"`
	Levels              Levels `json:"levels"`
	LogID               int64  `json:"logId"`
	MinutesAfterWakeup  int    `json:"minutesAfterWakeup"`
	MinutesAsleep       int    `json:"minutesAsleep"`
	MinutesAwake        int    `json:"minutesAwake"`
	MinutesToFallAsleep int    `json:"minutesToFallAsleep"`
	LogType             string `json:"logType"`
	StartTime           string `json:"startTime"`
	TimeInBed           int    `json:"timeInBed"`
	Type                string `json:"type"`
}
type Levels struct {
	Data      []Data      `json:"data"`
	ShortData []ShortData `json:"shortData"`
	Summary   Summary     `json:"summary"`
}
type Data struct {
	DateTime string `json:"dateTime"`
	Level    string `json:"level"`
	Seconds  int    `json:"seconds"`
}
type ShortData struct {
	DateTime string `json:"dateTime"`
	Level    string `json:"level"`
	Seconds  int    `json:"seconds"`
}
type Deep struct {
	Count               int `json:"count"`
	Minutes             int `json:"minutes"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
}
type Light struct {
	Count               int `json:"count"`
	Minutes             int `json:"minutes"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
}
type Rem struct {
	Count               int `json:"count"`
	Minutes             int `json:"minutes"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
}
type Wake struct {
	Count               int `json:"count"`
	Minutes             int `json:"minutes"`
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"`
}
type Stages struct {
	Deep  int `json:"deep"`
	Light int `json:"light"`
	Rem   int `json:"rem"`
	Wake  int `json:"wake"`
}

func (s *SleepService) GetLogByDate(time time.Time) (*SleepEnvelope, error) {
	err := s.client.Call(http.MethodGet, nil, nil, nil)
	return nil, err
}
