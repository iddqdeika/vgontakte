package handlers

import (
	"alina/alina"
	"vgontakte/vgontakte"
)

type messageRaterCreator struct {
}

func (c *messageRaterCreator) CreateHandler(params map[string]interface{}, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := messageRater{}

	return &handler, nil
}

func (c *messageRaterCreator) ParseHandler(data *string, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := messageRater{}

	return &handler, nil
}

type messageRater struct {
	peerId int
}

func (r *messageRater) Handle(message alina.PrivateMessage, err error) {
	if message.GetText() == "+" && message.GetFwdMessages() != nil {

		fwd := message.GetFwdMessages()
		if fwd != nil {

		}
	}
}

func (r *messageRater) Jsonize() (*string, error) {
	panic("implement me")
}

func (r *messageRater) Meet(message alina.PrivateMessage) bool {
	if message.GetText() == "+" && message.GetFwdMessages() != nil {
		return true
	}

	return false
}

func (r *messageRater) Order() int {
	panic("implement me")
}

func (r *messageRater) RateMessage() {

}
