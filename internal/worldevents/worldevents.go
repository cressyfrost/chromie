// Golang program to show how
// to take input from the user
package worldevents

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// channelID   = "529930903556849698" // testing channel
	channelID   = "790534093799555103" // testing channel
	maxSchedule = 6                    // max posting sessions per day

	// Iskaara Tuskarr Community Feast
	WORLD_EVENT_FEAST = "feast"
	feastHour         = 3
	feastMinute       = 30

	// Siege on the Dragonbane Keep
	WORLD_EVENT_SIEGE_DRAGONBANE = "siege"
	siegeHour                    = 2
	siegeMinute                  = 0
)

// SetNextEvents set the next time for feast and siege
func SetNextEvents() map[string]string {
	m := make(map[string]string)

	// Manually set next Feast used for caluclations
	fmt.Println("Specify the next Community Feast in HH:MM format: ")
	var feast string
	fmt.Scanln(&feast)

	// Manually set next Siege used for caluclations
	fmt.Println("Specify the next Siege of the Dragonbane Keep in HH:MM format: ")
	var siege string
	fmt.Scanln(&siege)

	m[WORLD_EVENT_FEAST] = feast
	m[WORLD_EVENT_SIEGE_DRAGONBANE] = siege

	return m
}

// PostWorldEventSchedule posts the next x world events schedule
func PostWorldEventSchedule(s *discordgo.Session, schedule map[string]string) {
	re := regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)
	ft := re.MatchString(schedule[WORLD_EVENT_FEAST])
	st := re.MatchString(schedule[WORLD_EVENT_SIEGE_DRAGONBANE])
	if !ft {
		postErrorSchedule(s, errors.New("invalid manually set next Feast time ("+schedule[WORLD_EVENT_FEAST]+")"))
		return
	}

	if !st {
		postErrorSchedule(s, errors.New("invalid manually set next Siege time("+schedule[WORLD_EVENT_SIEGE_DRAGONBANE]+")"))
		return
	}

	// Construct Community Feast Schedule
	f := strings.Split(schedule[WORLD_EVENT_FEAST], ":")
	fs, err := constructSchedule(WORLD_EVENT_FEAST, f[0], f[1])
	if err != nil {
		postErrorSchedule(s, err)
	}

	// Construct Siege on the Dragonbane Keep Schedule
	sg := strings.Split(schedule[WORLD_EVENT_SIEGE_DRAGONBANE], ":")
	sgs, err := constructSchedule(WORLD_EVENT_SIEGE_DRAGONBANE, sg[0], sg[1])
	if err != nil {
		postErrorSchedule(s, err)
	}

	var feasts, sieges string
	feasts = "\n```"
	for _, v := range fs {
		feasts += v + "\n"
	}
	feasts += "```"

	sieges = "\n```"
	for _, v := range sgs {
		sieges += v + "\n"
	}
	sieges += "```"

	footers := "\n~~React on this message with  :cook: and :european_castle: _emote_ to get notified `10 Minutes Before` every session begins.~~ soon tm"

	s.ChannelMessageSend(channelID, "Hello, Champions. Here is the **World Events** schedule for today!")
	s.ChannelMessageSend(channelID, ":cook: **Community Feast** :cook: "+feasts+"\n")
	s.ChannelMessageSend(channelID, ":european_castle: **Siege on the Dragonbane Keep** :european_castle: "+sieges)
	s.ChannelMessageSend(channelID, footers)

}

func constructSchedule(worldEventType string, hours string, minutes string) (schedule []string, err error) {
	h, err := strconv.ParseInt(hours, 10, 64)
	if err != nil {
		return nil, err
	}

	m, err := strconv.ParseInt(minutes, 10, 64)
	if err != nil {
		return nil, err
	}

	timein := time.Now().Local()
	for i := 0; i < maxSchedule; i++ {

		if i == 0 {
			// log.Println(timein.Format("15:04 WIB"))
			timein = timein.Add(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
			// log.Println(timein.Format("15:04 WIB"))
		} else {
			if worldEventType == WORLD_EVENT_FEAST {
				timein = timein.Add(time.Hour*time.Duration(feastHour) + time.Minute*time.Duration(feastMinute))
			} else if worldEventType == WORLD_EVENT_SIEGE_DRAGONBANE {
				timein = timein.Add(time.Hour*time.Duration(siegeHour) + time.Minute*time.Duration(siegeMinute))
			}
		}

		schedule = append(schedule, timein.Format("15:04 WIB"))
		// log.Println(timein.Format("15:04 WIB"))
	}
	// log.Println(schedule)
	return
}

func postErrorSchedule(s *discordgo.Session, err error) {
	s.ChannelMessageSend(channelID, "Uh oh, i was supposed to post today's schedule, but looks like something broke in this timeline :face_with_spiral_eyes: ")
	s.ChannelMessageSend(channelID, "Please notify my Master **Nozdormu** with this error message to get this fixed! :envelope_with_arrow: ")
	s.ChannelMessageSend(channelID, "```"+err.Error()+"```")
}
