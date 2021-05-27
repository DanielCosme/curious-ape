package sync

import "github.com/danielcosme/curious-ape/internal/sync/fitbit"

type Collectors struct {
	Sleep *fitbit.SleepCollector
}

func NewCollectors() *Collectors {
	return &Collectors{
		Sleep: fitbit.Fitbit,
	}
}
