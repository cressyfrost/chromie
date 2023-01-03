package discord

import (
	"cressyfrost/chromie/internal/worldevents"
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers Commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all Commands after shutdowning or not")
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "chromie",
			Description: "Main command for Chromie Bot",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "activity",
					Description: "Activity type / World Events type",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "feast",
							Value: "feast",
						},
						{
							Name:  "siege",
							Value: "siege",
						},
					},
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"chromie": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))
			var strmargs string
			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["activity"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				margs = append(margs, option.StringValue())
				strmargs = fmt.Sprintf("%v", margs)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: generateResponse(s, i.ChannelID, strmargs),
				},
			})
		},
	}
)

func generateResponse(s *discordgo.Session, channelID string, msg string) string {
	worldevents.PostWorldEventSchedule(s, channelID, msg)
	return "Yes, Chef!"
}
