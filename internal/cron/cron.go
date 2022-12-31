package cron

import (
	"cressyfrost/chromie/internal/worldevents"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

// Run registers and run the specified cron jobs
func Run(s *discordgo.Session) {
	log.Println("Initializing World Events cron.")
	c := cron.New()

	// Community Feast Scheduler
	for _, v := range worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_FEAST+"-raw"] {
		eventTime, err := time.Parse(time.UnixDate, v)
		if err != nil {
			return
		}
		// subtract time so we get notified before the actual events begins
		newTime := eventTime.Add(-time.Minute * worldevents.NotificationsValueBefore)

		day := strconv.Itoa(newTime.Day())
		month := strconv.Itoa(int(newTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour := strconv.Itoa(newTime.Hour())
		minute := strconv.Itoa(newTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		ct := minute + " " + hour + " " + day + " " + month + " *"

		c.AddFunc(ct, func() {
			log.Println("[FEAST] Cron executed that was scheduled for " + ct)
			worldevents.PostWorldEventReminder(s, worldevents.WORLD_EVENT_FEAST)
		})
	}

	// Siege on the Dragonbane Keep Scheduler
	for _, v := range worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_SIEGE_DRAGONBANE+"-raw"] {
		eventTime, err := time.Parse(time.UnixDate, v)
		if err != nil {
			return
		}
		// subtract time so we get notified before the actual events begins
		newTime := eventTime.Add(-time.Minute * worldevents.NotificationsValueBefore)

		day := strconv.Itoa(newTime.Day())
		month := strconv.Itoa(int(newTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour := strconv.Itoa(newTime.Hour())
		minute := strconv.Itoa(newTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		ct := minute + " " + hour + " " + day + " " + month + " *"

		c.AddFunc(ct, func() {
			log.Println("[SIEGE] Cron executed that was scheduled for " + ct)
			worldevents.PostWorldEventReminder(s, worldevents.WORLD_EVENT_SIEGE_DRAGONBANE)
		})
	}

	c.Start()

}
