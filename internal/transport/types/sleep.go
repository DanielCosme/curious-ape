package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type SleepLogTransport struct {
	ID          int    `json:"id"`
	Date        string `json:"date,omitempty"`
	IsMainSleep bool   `json:"is_main_sleep"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	TimeAsleep  string `json:"time_asleep"`
	TimeAwake   string `json:"time_awake"`
}

func FromSleepLogToTransport(sl *entity.SleepLog) *SleepLogTransport {
	var slt *SleepLogTransport
	if sl != nil {
		slt = &SleepLogTransport{
			ID:          sl.ID,
			IsMainSleep: sl.IsMainSleep,
			StartTime:   sl.StartTime.Format(entity.Timestamp),
			EndTime:     sl.EndTime.Format(entity.Timestamp),
			TimeAsleep:  sl.MinutesAsleep.String(),
			TimeAwake:   sl.MinutesAwake.String(),
		}

		if sl.Day != nil {
			slt.Date = sl.Day.Date.Format(entity.ISO8601)
		}
	}
	return slt
}

func FromSleepLogToTransportSlice(sls []*entity.SleepLog) []*SleepLogTransport {
	slst := []*SleepLogTransport{}
	for _, s := range sls {
		slst = append(slst, FromSleepLogToTransport(s))
	}
	return slst
}
