package handlers

import (
	"alina/alina"
	"encoding/json"
	"fmt"
	"strconv"
	"vgontakte/vgontakte"
)

type echoHandlerConfig struct {
	peerId int `json:"peer_id"`
}

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

func (c *echoHandlerCreator) ParseHandler(data *string, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	var cfg echoHandlerConfig

	err := json.Unmarshal([]byte(*data), &cfg)
	if err != nil {
		return nil, err
	}

	return &echoHandler{peerId: cfg.peerId}, nil
}

type echoHandler struct {
	peerId int
	alina  alina.Alina
}

func (h *echoHandler) Jsonize() (*string, error) {
	cfg := &echoHandlerConfig{peerId: h.peerId}

	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	result := string(data)
	return &result, nil
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
	h.alina.GetMessagesApi().SendMessageWithForward(strconv.Itoa(message.GetPeerId()), message.GetText(), []string{strconv.Itoa(message.GetId())})
}
