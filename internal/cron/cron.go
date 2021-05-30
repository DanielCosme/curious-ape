package cron

import (
	"log"
	"time"

	"github.com/danielcosme/curious-ape/internal/sync"
	"github.com/go-co-op/gocron"
)

type Cron struct {
	Collector *sync.Collectors
}

func (cron *Cron) Start() {
	time.Sleep(1 * time.Second)
	t := time.Now().Location()
	s := gocron.NewScheduler(t)

	ping, err := s.Every(1).Day().Tag("ping").At("00:00").Do(func() { log.Println("PING") })
	job1, err := s.Every(1).Day().Tag("sleep").At("18:00").Do(cron.Collector.GetTodayLog)
	job2, err := s.Every(1).Day().Tag("habits").At("00:01").Do(cron.Collector.InitializeDayHabits)
	if err != nil {
		log.Println(err)
	}

	log.Println("CRON JOB", job1.Tags()[0], job1.ScheduledAtTime())
	log.Println("CRON JOB", job2.Tags()[0], job2.ScheduledAtTime())
	log.Println("CRON JOB", ping.Tags()[0], ping.ScheduledAtTime())

	// Just in case
	cron.Collector.InitializeDayHabits()
	s.StartAsync()
}
