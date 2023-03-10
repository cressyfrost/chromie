module cressyfrost/chromie

go 1.18

require github.com/bwmarrin/discordgo v0.26.1

replace internal/worldevents => ././internal/worldevents

replace internal/discord => ././internal/discord

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/robfig/cron/v3 v3.0.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
