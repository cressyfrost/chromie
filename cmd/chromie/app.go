package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cressyfrost/chromie/internal/discord"
	"cressyfrost/chromie/internal/worldevents"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	// parse token from discord.token file in root directory
	var err error
	Token, err = parseToken("../../discord.token")
	if err != nil {
		Token, err = parseToken("discord.token")
		if err != nil {
			log.Fatalf("error getting discord client token")
		}
	}
}

func main() {
	discordClient, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discordClient session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discordClient.AddHandler(discord.MessageCreate)

	// Register the slash commands handler
	discordClient.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := discord.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	discordClient.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Open a websocket connection to discordClient and begin listening.
	err = discordClient.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(discord.Commands))
	for i, v := range discord.Commands {
		cmd, err := discordClient.ApplicationCommandCreate(discordClient.State.User.ID, *discord.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// In this example, we only care about receiving message events.
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages

	worldevents.PostInitialWorldEventSchedule(discordClient, worldevents.SetNextEvents())

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Chromie is now running on this timeline.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if *discord.RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		registeredCommands, err := discordClient.ApplicationCommands(discordClient.State.User.ID, *discord.GuildID)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := discordClient.ApplicationCommandDelete(discordClient.State.User.ID, *discord.GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	// Cleanly close down the discordClient session.
	discordClient.Close()
}

func parseToken(filename string) (string, error) {
	var token string
	f, err := os.Open(filename)

	if err != nil {
		return token, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		token += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return token, err
	}

	return token, nil
}
