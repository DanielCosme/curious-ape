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
	time.Sleep(1 * time.Millisecond)
	t := time.Now().Location()
	s := gocron.NewScheduler(t)

	yesterday := time.Now().AddDate(0, 0, -1)

	ping, err := s.Every(1).Day().Tag("ping").At("00:00").Do(func() { log.Println("PING") })
	job1, err := s.Every(1).Day().Tag("sleep").At("18:00").Do(cron.Collector.Sleep.GetCurrentDayRecord)
	job2, err := s.Every(1).Day().Tag("habits").At("00:01").Do(cron.Collector.InitializeDayHabits)
	work, err := s.Every(1).Day().Tag("work").At("04:00").Do(cron.Collector.Work.GetRecord, yesterday)
	if err != nil {
		log.Println(err)
	}

	log.Println("CRON JOB", job1.Tags()[0], job1.ScheduledAtTime())
	log.Println("CRON JOB", job2.Tags()[0], job2.ScheduledAtTime())
	log.Println("CRON JOB", ping.Tags()[0], ping.ScheduledAtTime())
	log.Println("CRON JOB", work.Tags()[0], work.ScheduledAtTime())

	// Just in case
	cron.Collector.InitializeDayHabits()

	s.StartAsync()
}
