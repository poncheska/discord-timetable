package app

import (
	"github.com/bwmarrin/discordgo"
	bot "github.com/poncheska/discord-timetable/internal/bot"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Start(){
	token := os.Getenv("BOT_TOKEN")
	if token == ""{
		logrus.Fatal("no bot token configured")
	}
	ttLink := os.Getenv("TT_LINK")
	if ttLink == ""{
		logrus.Fatal("no timetable link configured")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.Fatal(err)
	}

	b := bot.NewBot(dg, ttLink)

	err = b.ConfigureAndOpen()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("bot started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = b.Close()
	if err != nil {
		logrus.Fatal(err)
	}
}