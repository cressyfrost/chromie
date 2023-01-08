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

	var eventList []string

	eventList = append(eventList, worldevents.WORLD_EVENT_FEAST)
	eventList = append(eventList, worldevents.WORLD_EVENT_SIEGE_DRAGONBANE)

	for _, worldEventType := range eventList {
		worldEventType := worldEventType
		eventTime := worldevents.GetFixedEventSchedule(worldEventType)
		// subtract time so we get notified before the actual events begins
		newTime := eventTime.Add(-time.Minute * worldevents.NotificationsValueBefore)

		// convert to CRON format
		day := strconv.Itoa(newTime.Day())
		month := strconv.Itoa(int(newTime.Month()))
		// year := strconv.Itoa(newTime.Year())
		hour := strconv.Itoa(newTime.Hour())
		minute := strconv.Itoa(newTime.Minute())
		// second := strconv.Itoa(newTime.Second())

		ct := minute + " " + hour + " " + day + " " + month + " *"

		// c.AddFunc(ct, postReminder(s))
		log.Println("[" + worldEventType + "] Scheduling events for " + ct)
		c.AddFunc(ct, func() {
			var dh, dt string
			if worldEventType == worldevents.WORLD_EVENT_FEAST {
				dh = strconv.Itoa(worldevents.WORLD_EVENT_FEAST_HOUR)
				dt = strconv.Itoa(worldevents.WORLD_EVENT_FEAST_MINUTE)
			} else if worldEventType == worldevents.WORLD_EVENT_SIEGE_DRAGONBANE {
				dh = strconv.Itoa(worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_HOUR)
				dt = strconv.Itoa(worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_MINUTE)
			}
			log.Println("[" + worldEventType + "] Master Cron executed " + eventTime.String())
			postReminder(s, worldEventType)
			c.AddFunc("@every "+dh+"h"+dt+"m", func() {
				log.Println("[" + worldEventType + "] Child Cron executed ")
				postReminder(s, worldEventType)
			})
			c.Start()
		})
		c.Start()
	}

	// daily notifications every morning 9AM WIB
	c.AddFunc("0 9 * * *", func() {
		log.Println("[Daily/Morning] Cron executed 0 9 * * *")
		worldevents.PostDailyWorldEventSchedule(s)
	})
	c.Start()

}

func postReminder(s *discordgo.Session, worldEventType string) {
	worldevents.PostWorldEventReminder(s, worldEventType)
}

// func postReminder(c *cron.Cron, s *discordgo.Session, worldEventType string) {
// 	// log.Println(worldevents.GetWorldEventSchedule())
// 	log.Println("[" + worldEventType + "] Cron executed")
// 	worldevents.PostWorldEventReminder(s, worldEventType)
// 	// log.Println(worldevents.GetWorldEventSchedule())

// 	strNextTime := worldevents.GetWorldEventSchedule()[worldEventType+"-raw"][len(worldevents.GetWorldEventSchedule()[worldEventType+"-raw"])-1]
// 	nextTime, err := time.Parse(time.UnixDate, strNextTime)
// 	if err != nil {
// 		return
// 	}

// 	var h, m time.Duration
// 	if worldEventType == worldevents.WORLD_EVENT_FEAST {
// 		h = worldevents.WORLD_EVENT_FEAST_HOUR
// 		m = worldevents.WORLD_EVENT_FEAST_MINUTE
// 	} else if worldEventType == worldevents.WORLD_EVENT_SIEGE_DRAGONBANE {
// 		h = worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_HOUR
// 		m = worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_MINUTE
// 	}

// 	nextTime = nextTime.Add(time.Hour * h)
// 	nextTime = nextTime.Add(time.Minute * m)

// 	worldevents.UpdateSchedule(worldEventType, nextTime)

// 	// also subtract new time so we get notified before the actual events begins
// 	newNextTime := nextTime.Add(-time.Minute * worldevents.NotificationsValueBefore)
// 	day := strconv.Itoa(newNextTime.Day())
// 	month := strconv.Itoa(int(newNextTime.Month()))
// 	// year := strconv.Itoa(newTime.Year())
// 	hour := strconv.Itoa(newNextTime.Hour())
// 	minute := strconv.Itoa(newNextTime.Minute())
// 	// second := strconv.Itoa(newTime.Second())

// 	nt := minute + " " + hour + " " + day + " " + month + " *"

// 	// register the next cron
// 	c.AddFunc(nt, func() {
// 		postReminder(c, s, worldEventType)
// 	})
// 	c.Start()
// }
