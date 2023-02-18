package discord

import (
	"cressyfrost/chromie/internal/worldevents"
	"log"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time reaction is added
func ReactionAdded(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if r.UserID == s.State.User.ID {
		return
	}

	if r.Emoji.Name == worldevents.WORLD_EVENT_FEAST_EMOTE_UNICODE {
		worldevents.SetWorldEventSubscribers(worldevents.WORLD_EVENT_FEAST, r.Member.User.ID)
	} else if r.Emoji.Name == worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE_UNICODE {
		worldevents.SetWorldEventSubscribers(worldevents.WORLD_EVENT_SIEGE_DRAGONBANE, r.Member.User.ID)
	}

	log.Println("Reaction Added Debug")
	log.Println("<@" + r.Member.User.ID + "> reacted with: " + r.Emoji.ID + " -> `" + r.Emoji.Name + "`")
	log.Println(worldevents.GetWorldEventSubscribers())

}

// This function will be called (due to AddHandler above) every time reaction is added
func ReactionRemoved(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if r.UserID == s.State.User.ID {
		return
	}

	if r.MessageReaction.Emoji.Name == worldevents.WORLD_EVENT_FEAST_EMOTE_UNICODE {
		worldevents.UnsetWorldEventSubscribers(worldevents.WORLD_EVENT_FEAST, r.MessageReaction.UserID)
	} else if r.MessageReaction.Emoji.Name == worldevents.WORLD_EVENT_SIEGE_DRAGONBANE_EMOTE_UNICODE {
		worldevents.UnsetWorldEventSubscribers(worldevents.WORLD_EVENT_SIEGE_DRAGONBANE, r.MessageReaction.UserID)
	}

	log.Println("Reaction Removed Debug")
	log.Println("<@" + r.MessageReaction.UserID + "> reacted with: " + "n/a" + " -> `" + r.MessageReaction.Emoji.Name + "`")
	log.Println(worldevents.GetWorldEventSubscribers())

}
