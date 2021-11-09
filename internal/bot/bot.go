package bot

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/poncheska/discord-timetable/internal/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	dg       *discordgo.Session
	ttLink   string
	ttSpamID string
	ss       *utils.SecretSanta
}

func NewBot(dg *discordgo.Session, ttLink string, ttSpamID string) *Bot {
	return &Bot{
		dg:       dg,
		ttLink:   ttLink,
		ttSpamID: ttSpamID,
		ss:       utils.NewSecretSanta(),
	}
}

func (b *Bot) ConfigureAndOpen() error {
	b.dg.AddHandler(b.TTHandler)

	if b.ttSpamID != "" {
		c := cron.New()
		_, err := c.AddFunc("0 13 * * SUN", func() {
			sendTT(b.dg, b.ttLink, b.ttSpamID)
			logrus.Info("cron timetable spam")
		})
		if err != nil {
			logrus.Error(err)
		} else {
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

	if cmd[0] == "!qr" {
		if len(cmd) < 2 {
			_, err := s.ChannelMessageSend(m.ChannelID, `"empty"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		file, err := qrcode.Encode(strings.Join(cmd[1:], " "), qrcode.Medium, 512)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, `"encoding error"`)
			logrus.Error(err)
		}
		s.ChannelFileSend(m.ChannelID, "qr.png", bytes.NewBuffer(file))
	}

	if cmd[0] == "!ss" {
		if len(cmd) < 2 {
			_, err := s.ChannelMessageSend(m.ChannelID, `Не верная команда. Используй: "!ss info"`)
			if err != nil {
				logrus.Error(err)
			}
			return
		}
		switch cmd[1] {
		case "info":
			if len(cmd) == 3 {
				gameId := cmd[2]
				if _, ok := b.ss.Games[gameId]; !ok {
					_, err := s.ChannelMessageSend(m.ChannelID, "Игры с таким id не существует")
					if err != nil {
						logrus.Error(err)
					}
					return
				}
				gInfo := b.ss.Games[gameId].GetMembersStr(gameId)
				_, err := s.ChannelMessageSend(m.ChannelID, gInfo)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			_, err := s.ChannelMessageSend(m.ChannelID, utils.SSMan)
			if err != nil {
				logrus.Error(err)
			}
			return
		case "new":
			gameId := b.ss.NewGame(m.ChannelID)
			_, err := s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("Новый тайный санта создан с идентификатором: %v", gameId))
			if err != nil {
				logrus.Error(err)
			}
			return
		case "join":
			if len(cmd) < 4 {
				_, err := s.ChannelMessageSend(m.ChannelID, `Не верная команда. Используй: "!ss info"`)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			gameId := cmd[2]
			memName := strings.Join(cmd[3:], " ")
			if _, ok := b.ss.Games[gameId]; !ok {
				_, err := s.ChannelMessageSend(m.ChannelID, "Игры с таким id не существует")
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			err := b.ss.Games[gameId].NewMember(m.ChannelID, memName)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ты был зарегистрирован в тайном санте %v, как: %v",
				gameId, memName))
			if err != nil {
				logrus.Error(err)
			}
		case "start":
			if len(cmd) != 3 {
				_, err := s.ChannelMessageSend(m.ChannelID, `Не верная команда. Используй: "!ss info"`)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			gameId := cmd[2]
			if _, ok := b.ss.Games[gameId]; !ok {
				_, err := s.ChannelMessageSend(m.ChannelID, "Игры с таким id не существует")
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			gInfo := b.ss.Games[gameId].GetMembersStr(gameId)
			ssMap, err := b.ss.StartGame(m.ChannelID, gameId)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			for chId, name := range ssMap {
				_, err := s.ChannelMessageSend(chId, gInfo)
				if err != nil {
					logrus.Error(err)
				}
				_, err = s.ChannelMessageSend(chId, fmt.Sprintf("Ты тайный санта у ***%v***\n",
					strings.ToUpper(name)))
				if err != nil {
					logrus.Error(err)
				}
			}
		case "delete":
			if len(cmd) != 3 {
				_, err := s.ChannelMessageSend(m.ChannelID, `Не верная команда. Используй: "!ss info"`)
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			gameId := cmd[2]
			notifyMap, err := b.ss.DeleteGame(m.ChannelID, gameId)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					logrus.Error(err)
				}
				return
			}
			msg := fmt.Sprintf("Игра %v удалена организатором", gameId)
			for _, chId := range notifyMap {
				_, err := s.ChannelMessageSend(chId, msg)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
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
