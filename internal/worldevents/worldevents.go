package worldevents

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var WorldEventNextSchedule = make(map[string]string)
var WorldEventSchedule = make(map[string][]string)
var m sync.RWMutex

const (
	channelID = "529930903556849698" // testing channel
	// channelID   = "790534093799555103" // actual channel
	maxSchedule              = 6 // max posting sessions per day
	NotificationsValueBefore = 10
	NotificationsUnitBefore  = "Minutes"

	// Iskaara Tuskarr Community Feast
	WORLD_EVENT_FEAST               = "feast"
	WORLD_EVENT_FEAST_FULLNAME      = "Iskaara Tuskarr Community Feast"
	WORLD_EVENT_FEAST_EMOTE         = ":cook:"
	WORLD_EVENT_FEAST_EMOTE_UNICODE = "🧑‍🍳"
	WORLD_EVENT_FEAST_HOUR          = 3
	WORLD_EVENT_FEAST_MINUTE        = 30

	// Siege on the Dragonbane Keep
	WORLD_EVENT_SIEGE_DRAGONBANE               = "siege"
	WORLD_EVENT_SIEGE_DRAGONBANE_FULLNAME      = "Siege on the Dragonbane Keep"
	WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE         = ":european_castle:"
	WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE_UNICODE = "🏰"
	WORLD_EVENT_SIEGE_DRAGONBANE_HOUR          = 2
	WORLD_EVENT_SIEGE_DRAGONBANE_MINUTE        = 0
)

// GetNextEvents set the next time for feast and siege
func GetNextEvents() map[string]string {
	m.RLock()
	defer m.RUnlock()
	return WorldEventNextSchedule
}

// SetNextEvents set the next time for feast and siege
func SetNextEvents(s *discordgo.Session) error {
	m.Lock()
	// Manually set next Feast used for calculations
	fmt.Println("Specify the next Community Feast in HH:MM format: ")
	var feast string
	fmt.Scanln(&feast)

	// Manually set next Siege used for calculations
	fmt.Println("Specify the next Siege of the Dragonbane Keep in HH:MM format: ")
	var siege string
	fmt.Scanln(&siege)

	re := regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)
	ft := re.MatchString(feast)
	st := re.MatchString(siege)
	if !ft {
		err := errors.New("invalid manually set next Feast time (" + feast + ")")
		postErrorSchedule(s, err)
		return err
	}

	if !st {
		err := errors.New("invalid manually set next Siege time(" + siege + ")")
		postErrorSchedule(s, err)
		return err
	}

	WorldEventNextSchedule[WORLD_EVENT_FEAST] = feast
	WorldEventNextSchedule[WORLD_EVENT_SIEGE_DRAGONBANE] = siege
	m.Unlock()
	return nil
}

// GetNextEvents set the next time for feast and siege
func GetWorldEventSchedule() map[string][]string {
	m.RLock()
	defer m.RUnlock()
	return WorldEventSchedule
}

// SetNextEvents set the next time for feast and siege
func SetWorldEventSchedule(worldEventType string, clean []string, raw []string) {
	m.Lock()
	WorldEventSchedule[worldEventType] = clean
	WorldEventSchedule[worldEventType+"-raw"] = raw
	m.Unlock()
}

// PostInitialWorldEventSchedule posts the next x world events schedule
func PostInitialWorldEventSchedule(s *discordgo.Session) {
	var feasts, sieges string
	feasts = "\n```"
	for _, v := range GetWorldEventSchedule()[WORLD_EVENT_FEAST] {
		feasts += v + "\n"
	}
	feasts += "```"

	sieges = "\n```"
	for _, v := range GetWorldEventSchedule()[WORLD_EVENT_SIEGE_DRAGONBANE] {
		sieges += v + "\n"
	}
	sieges += "```"

	footers := "\nReact on this message with  " + WORLD_EVENT_FEAST_EMOTE + " and " + WORLD_EVENT_SIEGE_DRAGONBANE + " _emote_ to get notified `10 Minutes Before` every session begins."

	s.ChannelMessageSend(channelID, "Hello, Champions. Here is the **World Events** schedule for today!")
	s.ChannelMessageSend(channelID, "To display this later, type the `/chromie` command")
	s.ChannelMessageSend(channelID, WORLD_EVENT_FEAST_EMOTE+" **Community Feast** "+WORLD_EVENT_FEAST_EMOTE+" "+feasts+"\n")
	s.ChannelMessageSend(channelID, WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE+" **Siege on the Dragonbane Keep** "+WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE+" "+sieges)
	s.ChannelMessageSend(channelID, footers)

}

// PostWorldEventSchedule posts the next x world events schedule
func PostWorldEventSchedule(s *discordgo.Session, originChannelID string, worldEventType string) {
	schedule := GetNextEvents()
	if len(schedule[WORLD_EVENT_FEAST]) < 1 || len(schedule[WORLD_EVENT_SIEGE_DRAGONBANE]) < 1 {
		return
	}

	var feasts, sieges string
	feasts = "\n```"
	for _, v := range WorldEventSchedule[WORLD_EVENT_FEAST] {
		feasts += v + "\n"
	}
	feasts += "```"

	sieges = "\n```"
	for _, v := range WorldEventSchedule[WORLD_EVENT_SIEGE_DRAGONBANE] {
		sieges += v + "\n"
	}
	sieges += "```"

	var emote string
	if strings.Contains(worldEventType, WORLD_EVENT_FEAST) {
		emote = WORLD_EVENT_FEAST_EMOTE
	} else if strings.Contains(worldEventType, WORLD_EVENT_SIEGE_DRAGONBANE) {
		emote = WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE
	}

	footers := "\nReact on this message with " + emote + " _emote_ to get notified `10 Minutes Before` every session begins."

	s.ChannelMessageSend(originChannelID, "Here is your requested schedule for today!")
	if strings.Contains(worldEventType, WORLD_EVENT_FEAST) {
		s.ChannelMessageSend(originChannelID, WORLD_EVENT_FEAST_EMOTE+" **Community Feast** "+WORLD_EVENT_FEAST_EMOTE+" "+feasts+"\n")
	} else if strings.Contains(worldEventType, WORLD_EVENT_SIEGE_DRAGONBANE) {
		s.ChannelMessageSend(originChannelID, WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE+" **Siege on the Dragonbane Keep** "+WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE+" "+sieges)
	}
	s.ChannelMessageSend(originChannelID, footers)

}

// PostWorldEventReminder posts the next x world events schedule to each subscribers
func PostWorldEventReminder(s *discordgo.Session, worldEventType string) {
	var headers, emote, footers string
	if worldEventType == WORLD_EVENT_FEAST {
		headers = WORLD_EVENT_FEAST_FULLNAME
		emote = WORLD_EVENT_FEAST_EMOTE
		subscribers := GetWorldEventSubscribers()
		for _, v := range subscribers[WORLD_EVENT_FEAST] {
			footers += "<@" + v + "> "
		}

	} else if worldEventType == WORLD_EVENT_SIEGE_DRAGONBANE {
		headers = WORLD_EVENT_SIEGE_DRAGONBANE_FULLNAME
		emote = WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE
		subscribers := GetWorldEventSubscribers()
		for _, v := range subscribers[WORLD_EVENT_SIEGE_DRAGONBANE] {
			footers += "<@" + v + "> "
		}
	}

	s.ChannelMessageSend(channelID, emote+" **"+headers+"** is starting in **"+strconv.Itoa(NotificationsValueBefore)+" "+NotificationsUnitBefore+"**! "+emote)
	s.ChannelMessageSend(channelID, footers)

}

func ConstructSchedule(worldEventType string) (clean []string, raw []string, err error) {
	sch := GetNextEvents()[worldEventType]
	f := strings.Split(sch, ":")

	h, err := strconv.ParseInt(f[0], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	m, err := strconv.ParseInt(f[1], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	timein := time.Now().Local()
	for i := 0; i < maxSchedule; i++ {

		if i == 0 {
			timein = timein.Add(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
		} else {
			if worldEventType == WORLD_EVENT_FEAST {
				timein = timein.Add(time.Hour*time.Duration(WORLD_EVENT_FEAST_HOUR) + time.Minute*time.Duration(WORLD_EVENT_FEAST_MINUTE))
			} else if worldEventType == WORLD_EVENT_SIEGE_DRAGONBANE {
				timein = timein.Add(time.Hour*time.Duration(WORLD_EVENT_SIEGE_DRAGONBANE_HOUR) + time.Minute*time.Duration(WORLD_EVENT_SIEGE_DRAGONBANE_MINUTE))
			}
		}

		clean = append(clean, timein.Format("15:04 WIB"))
		raw = append(raw, timein.Format(time.UnixDate))
	}

	SetWorldEventSchedule(worldEventType, clean, raw)
	return
}

func postErrorSchedule(s *discordgo.Session, err error) {
	s.ChannelMessageSend(channelID, "Uh oh, i was supposed to post today's schedule, but looks like something broke in this timeline :face_with_spiral_eyes: ")
	s.ChannelMessageSend(channelID, "Please notify my Master **Nozdormu** with this error message to get this fixed! :envelope_with_arrow: ")
	s.ChannelMessageSend(channelID, "```"+err.Error()+"```")
}
