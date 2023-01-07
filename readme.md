# Chromie
## _World of Warcraft Dragonflight World Events Discord Bot Reminder_



[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

Chromie is a World Events Discord Bot Reminder for:

- Iskaara Tuskarr Community Feast
- Siege on the Dragonbane Keep

## Features

- Reminds and ping subscribed users whenever feast or siege is starting soon.
- Simply subscribe by reacting and unsubscribe by unreacting.
- No database needed.


## Installation

Chromie requires [Go](https://go.dev/) to build. And you'll need to host it somewhere yourself.
You can either install Go+Git and build it directly, or build it anywhere and run the executables on your server.

You'll need to create `discord.token` file on the same directory as the exectuables and write your bot token here.
```sh
nano discord.token 
paste your discord bot token here
ctrl+o
ctrl+x
```

To run it

```sh
cd cmd/chromie
./chromie
type the time until next feasts in hh:mm format, e.g: 01:45
type the time until next siege in hh:mm format, e.g: 00:15
```

## Limitations due to Laziness

- The notifications won't get sent if you start the bot with the time until next feast/siege LESS THAN 10 minutes.
- This is due to incomplete initial cron logic during initialization.
- User notifications data is stored into memory, so if you reload/restart the bot it will be LOST.
- Apart from that, once you start the bot it should notify you indefinitely 
- (however if blizzard modifies the timers, or it's somehow shifted due to maintenance you'll need to restart the bot and provide the new next time until feast/siege).

