package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/poncheska/discord-timetable/internal/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	dg       *discordgo.Session
	ttLink   string
	ttSpamID string
}

func NewBot(dg *discordgo.Session, ttLink string, ttSpamID string) *Bot {
	return &Bot{
		dg:       dg,
		ttLink:   ttLink,
		ttSpamID: ttSpamID,
	}
}

func (b *Bot) ConfigureAndOpen() error {
	b.dg.AddHandler(b.TTHandler)

	if b.ttSpamID != "" {
		c := cron.New()
		_,err := c.AddFunc("0 13 * * SUN",func() {
			sendTT(b.dg, b.ttLink, b.ttSpamID)
			logrus.Info("cron timetable spam")
		})
		if err != nil{
			logrus.Error(err)
		}else{
			c.Start()
			logrus.Info("cron started: " + time.Now().String())
		}
	}

	return b.dg.Open()
}

func (b *Bot) Close() error {
	return b.dg.Close()
}

func (b *Bot) TTHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	logrus.Info(m.ChannelID)

	cmd := strings.Split(strings.TrimSpace(m.Content), " ")
	if len(cmd) == 0 {
		return
	}

	if cmd[0] == "!tt" {
		if len(cmd) > 1 {
			sendTT(s, cmd[1], m.ChannelID)
			return
		}
		sendTT(s, b.ttLink, m.ChannelID)
		return
	}

	if cmd[0] == "!r" {
		if len(cmd) != 3 {
			_, err := s.ChannelMessageSend(m.ChannelID, `"/r min max"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		min, err := strconv.Atoi(cmd[1])
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, `"/r min max"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		max, err := strconv.Atoi(cmd[2])
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, `"/r min max"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		if min > max {
			_, err := s.ChannelMessageSend(m.ChannelID, `"/r min max"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		res := min + r.Intn(max+1-min)

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("your random number is %v", res))
		if err != nil {
			logrus.Error(err)
		}
		return
	}
}

func sendTTByDays(s *discordgo.Session, tt *utils.Timetable, ChID string) {
	for _, d := range tt.Days {
		ss := d.GetString()
		for _, str := range ss {
			_, err := s.ChannelMessageSend(ChID, str)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func sendTT(s *discordgo.Session, ttLink string, spamID string) {
	tt, err := utils.ParseTimetable(ttLink)
	if err != nil {
		_, err = s.ChannelMessageSend(spamID, "error")
		if err != nil {
			logrus.Error(err)
		}
	} else {
		sendTTByDays(s, tt, spamID)
	}
}
