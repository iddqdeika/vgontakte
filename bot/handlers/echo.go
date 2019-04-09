package handlers

import (
	"alina/alina"
	"fmt"
	"strconv"
	"vgontakte/vgontakte"
)

type echoHandlerCreator struct {
}

func (c *echoHandlerCreator) CreateHandler(params map[string]interface{}, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := echoHandler{}

	if val, ok := params["peer_id"]; ok {
		switch k := val.(type) {
		case int:
			handler.peerId = val.(int)
		default:
			return nil, fmt.Errorf("cant parse peerid: %T type got instead of int", k)
		}
	}
	handler.alina = alina
	return &handler, nil
}

type echoHandler struct {
	peerId int
	alina  alina.Alina
}

func (h *echoHandler) Order() int {
	return 10
}

func (h *echoHandler) Meet(message alina.PrivateMessage) bool {

	if h.peerId == 0 || message.GetPeerId() == h.peerId {
		return true
	}
	return false
}

func (h *echoHandler) Handle(message alina.PrivateMessage, err error) {
	h.alina.GetMessagesApi().SendSimpleMessage(strconv.Itoa(message.GetPeerId()), message.GetText())
}
