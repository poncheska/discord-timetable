package utils

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const SSMan = `__***Инструкция к боту "Тайный Санта":***__
*Каждый участник должен писать боту в ЛС. Бот идентифицирует всех по id чата(то есть права организатора будут только в том чате, в котором создана игра).*
***"!ss new"*** - создать новую игру, будет выслан GAME_ID (организатор не является участником игры автоматически)
***"!ss info <GAME_ID>"*** - получить список участников тайного санты по GAME_ID
***"!ss join <GAME_ID> <YOUR_NAME>"*** - вступить в созданную игру по GAME_ID, YOUR_NAME - твоё имя (если вызвать будучи уже зарегистрированным имя обновится)
***"!ss start <GAME_ID>"*** - начать игру (доступно только организатору), после старт все данные об игре удаляются.
***"!ss delete <GAME_ID>"*** - удалить игру (доступно только организатору)
`

type SecretSanta struct {
	Games map[string]*Game
	mx    *sync.Mutex
}

type Game struct {
	Members map[string]string // map[member_id]member_name
	mx      *sync.Mutex
	OrgId   string
}

func NewSecretSanta() *SecretSanta {
	return &SecretSanta{
		Games: make(map[string]*Game),
		mx:    &sync.Mutex{},
	}
}

func (s *SecretSanta) NewGame(orgId string) string {
	s.mx.Lock()
	defer s.mx.Unlock()
	var id string
	ok := true
	for ok {
		id = GenerateRandomId(6)
		_, ok = s.Games[id]
	}
	s.Games[id] = &Game{
		Members: make(map[string]string),
		mx:      &sync.Mutex{},
		OrgId:   orgId,
	}
	return id
}

func (s *SecretSanta) DeleteGame(callerId string, gameId string) ([]string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, ok := s.Games[gameId]; !ok {
		return nil, errors.New("Тайного санты с таким ID не существует.")
	} else if s.Games[gameId].OrgId != callerId {
		return nil, errors.New("Вы не являетесь организатором данного тайного санты.")
	}
	notificationList := make([]string, 0, len(s.Games[gameId].Members))
	for k := range s.Games[gameId].Members {
		notificationList = append(notificationList, k)
	}
	delete(s.Games, gameId)
	return notificationList, nil
}

func (s *SecretSanta) StartGame(callerId string, gameId string) (map[string]string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, ok := s.Games[gameId]; !ok {
		return nil, errors.New("Тайного санты с таким ID не существует.")
	} else if s.Games[gameId].OrgId != callerId {
		return nil, errors.New("Вы не являетесь организатором данного тайного санты.")
	}
	res := Shuffle(s.Games[gameId].Members)
	delete(s.Games, gameId)
	return res, nil
}

func (g *Game) NewMember(memId string, memName string) error {
	g.mx.Lock()
	defer g.mx.Unlock()
	if oldName, ok := g.Members[memId]; ok {
		if oldName == memName {
			return errors.New("Ты уже участвуешь в этом тайном санте.")
		} else {
			g.Members[memId] = memName
			return errors.New(fmt.Sprintf("Твоё имя в этом тайном санте было заменено на %v", memName))
		}
	}
	g.Members[memId] = memName
	return nil
}

func (g *Game) GetMembersStr(gameId string) string {
	g.mx.Lock()
	defer g.mx.Unlock()
	sb := bytes.NewBuffer([]byte{})
	sb.WriteString(fmt.Sprintf("Список участвников тайного санты %v:\n", gameId))
	i := 1
	for _, v := range g.Members {
		sb.WriteString(fmt.Sprintf("**%v.** %v\n", i, v))
		i++
	}
	return sb.String()
}

const idSymbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = idSymbols[rand.Intn(len(idSymbols))]
	}
	return string(b)
}

func Shuffle(in map[string]string) map[string]string {
	arr := make([]string, len(in))
	arrShfl := make([]string, len(in))
	i := 0
	for k := range in {
		arr[i] = k
		arrShfl[i] = k
		i++
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arrShfl), func(i, j int) {
		if (arr[i] == arrShfl[i]) || (arr[j] == arrShfl[j]) {
			arrShfl[i], arrShfl[j] = arrShfl[j], arrShfl[i]
		}
	})

	if arr[len(in)-1] == arrShfl[len(in)-1] {
		arrShfl[0], arrShfl[len(in)-1] = arrShfl[len(in)-1], arrShfl[0]
	}
	for i := 0; i < len(in)-1; i++ {
		if arr[i] == arrShfl[i] {
			arrShfl[i], arrShfl[i+1] = arrShfl[i+1], arrShfl[i]
		}
	}

	res := make(map[string]string)
	for i := 0; i < len(arr); i++ {
		res[arr[i]] = in[arrShfl[i]]
	}
	return res
}
