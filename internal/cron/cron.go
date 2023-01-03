package cron

import (
	"cressyfrost/chromie/internal/worldevents"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func testCron(c *cron.Cron, s *discordgo.Session, worldEventType string) {
	// log.Println(worldevents.GetWorldEventSchedule())
	// log.Println("[FEAST] Cron executed that was scheduled for " + ct)
	worldevents.PostWorldEventReminder(s, worldevents.WORLD_EVENT_FEAST)
	worldevents.UpdateSchedule(worldevents.WORLD_EVENT_FEAST, nextTime)
	// log.Println(worldevents.GetWorldEventSchedule())

	// register the next cron
	c.AddFunc(nt, func() {
		// log.Println(worldevents.GetWorldEventSchedule())
		// log.Println("[FEAST] Cron executed that was scheduled for " + ct)
		worldevents.PostWorldEventReminder(s, worldevents.WORLD_EVENT_FEAST)
		worldevents.UpdateSchedule(worldevents.WORLD_EVENT_FEAST, nextTime)
		// log.Println(worldevents.GetWorldEventSchedule())
	})
	c.Start()
}

// Run registers and run the specified cron jobs
// todo: refactor into something maintainable
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
		strNextTime := worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_FEAST+"-raw"][len(worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_FEAST+"-raw"])-1]
		nextTime, err := time.Parse(time.UnixDate, strNextTime)
		if err != nil {
			return
		}
		nextTime = nextTime.Add(time.Hour * worldevents.WORLD_EVENT_FEAST_HOUR)
		nextTime = nextTime.Add(time.Minute * worldevents.WORLD_EVENT_FEAST_MINUTE)

		// also subtract new time so we get notified before the actual events begins
		newNextTime := nextTime.Add(-time.Minute * worldevents.NotificationsValueBefore)

		// convert to CRON format
		day := strconv.Itoa(newTime.Day())
		month := strconv.Itoa(int(newTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour := strconv.Itoa(newTime.Hour())
		minute := strconv.Itoa(newTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		ct := minute + " " + hour + " " + day + " " + month + " *"

		day = strconv.Itoa(newNextTime.Day())
		month = strconv.Itoa(int(newNextTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour = strconv.Itoa(newNextTime.Hour())
		minute = strconv.Itoa(newNextTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		nt := minute + " " + hour + " " + day + " " + month + " *"

		// c.AddFunc(ct, testCron(s))
		c.AddFunc(ct, func() {
			testCron(c, s, worldevents.WORLD_EVENT_FEAST)
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
		strNextTime := worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_SIEGE_DRAGONBANE+"-raw"][len(worldevents.GetWorldEventSchedule()[worldevents.WORLD_EVENT_SIEGE_DRAGONBANE+"-raw"])-1]
		nextTime, err := time.Parse(time.UnixDate, strNextTime)
		if err != nil {
			return
		}
		nextTime = nextTime.Add(time.Hour * worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_HOUR)
		nextTime = nextTime.Add(time.Minute * worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_MINUTE)

		// also subtract new time so we get notified before the actual events begins
		newNextTime := nextTime.Add(-time.Minute * worldevents.NotificationsValueBefore)

		// convert to CRON format

		day := strconv.Itoa(newTime.Day())
		month := strconv.Itoa(int(newTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour := strconv.Itoa(newTime.Hour())
		minute := strconv.Itoa(newTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		ct := minute + " " + hour + " " + day + " " + month + " *"

		day = strconv.Itoa(newNextTime.Day())
		month = strconv.Itoa(int(newNextTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour = strconv.Itoa(newNextTime.Hour())
		minute = strconv.Itoa(newNextTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		nt := minute + " " + hour + " " + day + " " + month + " *"

		c.AddFunc(ct, func() {
			// log.Println(worldevents.GetWorldEventSchedule())
			// log.Println("[SIEGE] Cron executed that was scheduled for " + ct)
			worldevents.PostWorldEventReminder(s, worldevents.WORLD_EVENT_SIEGE_DRAGONBANE)
			worldevents.UpdateSchedule(worldevents.WORLD_EVENT_SIEGE_DRAGONBANE, nextTime)
			// log.Println(worldevents.GetWorldEventSchedule())
		})
	}

	c.Start()

}
