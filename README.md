# Discord Tools Bot

[![Go report](http://goreportcard.com/badge/FlameInTheDark/dtbot)](http://goreportcard.com/report/FlameInTheDark/dtbot)

## Used APIs

* OpenWeatherMap
* Yandex Translate
* Newsapi.org
* Geonames.org
* cbr-xml-daily.ru

## How to use

Bot commands

Command | Description
------- | -----------
`!r join` | Add bot into you voice channel
`!r leave` | Remove bot from voice channel
`!r play [radio_station]` | Plays specified network radio station `!r play http://air2.radiorecord.ru:9003/rr_320`
`!r stop` | Stops radio
`!w [place]` | Shows the weather in a specified location `!w New York`
`!n [category]` | Displays news in the specified category `!n technology`
`!t [target_lang] [text]` | Translator `!t ru Hello world`
`!c` | Shows currencies (default from config)
`!c list` | Shows list of available currencies
`!c [currency]` | Shows specified currency `!c USD EUR`

## Build for docker

Easy way to build docker image for Ubuntu:

Clone reposytory and move inside app directory. Сompile app with command:

`go build`

Make (or use my) Dockerfile:

```Dockerfile
FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y ca-certificates ffmpeg
COPY . . 
ENTRYPOINT ["./dtbot"]
```

Build docker image

`docker build -t dtbot .`

Add environment variable `BOT_TOKEN` with token of discord bot.
And run container:

`docker run -d --rm -e BOT_TOKEN=$BOT_TOKEN --name dtbot dtbot:latest`