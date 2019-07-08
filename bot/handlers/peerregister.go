package handlers

import (
	"alina/alina"
	"strconv"
	"strings"
	"vgontakte/config"
	"vgontakte/vgontakte"
)

type peerRegistratorCreator struct{}

func (c *peerRegistratorCreator) CreateHandler(params map[string]interface{}, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := peerRegistrator{
		storage: bot.GetStrorage(),
		alina:   alina,
		logger:  bot.GetLogger(),
	}

	return &handler, nil
}

func (c *peerRegistratorCreator) ParseHandler(data *string, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := peerRegistrator{
		storage: bot.GetStrorage(),
		alina:   alina,
		logger:  bot.GetLogger(),
	}

	return &handler, nil
}

type peerRegistrator struct {
	storage vgontakte.Storage
	alina   alina.Alina
	logger  vgontakte.Logger
}

func (r *peerRegistrator) Order() int {
	return 2
}

func (r *peerRegistrator) Meet(message alina.PrivateMessage) bool {
	owner, err := config.GetCommandLineArgsConfig().GetInt("ownerpeer")
	return strings.Index(message.GetText(), "register ") == 0 && err == nil && (owner == message.GetPeerId() || owner == message.GetFromId())
}

func (r *peerRegistrator) Handle(message alina.PrivateMessage, err error) {
	temp := strings.Split(message.GetText(), " ")[1]
	if temp == "this" {
		r.storage.RegisterPeer(message.GetPeerId())
	}
	newPeer, err := strconv.Atoi(temp)
	if err != nil {
		r.logger.Info("cant parse peer " + temp)
		return
	}
	r.storage.RegisterPeer(newPeer)
}

func (r *peerRegistrator) Jsonize() (*string, error) {
	result := "{}"
	return &result, nil
}
