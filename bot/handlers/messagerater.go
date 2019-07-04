package handlers

import (
	"alina/alina"
	"strconv"
	"strings"
	"vgontakte/vgontakte"
)

type messageRaterCreator struct {
}

func (c *messageRaterCreator) CreateHandler(params map[string]interface{}, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := messageRater{
		storage: bot.GetStrorage(),
		alina:   alina,
		logger:  bot.GetLogger(),
	}

	return &handler, nil
}

func (c *messageRaterCreator) ParseHandler(data *string, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := messageRater{
		storage: bot.GetStrorage(),
		alina:   alina,
		logger:  bot.GetLogger(),
	}

	return &handler, nil
}

type messageRater struct {
	storage vgontakte.Storage
	alina   alina.Alina
	logger  vgontakte.Logger
	peerId  int
}

func (r *messageRater) Handle(message alina.PrivateMessage, err error) {
	if message.GetText() == "+" && message.GetFwdMessages() != nil {
		r.rateMessage(message)
	}
	if strings.Contains(strings.ToLower(message.GetText()), "топ перлов") {
		r.returnMessageTop(message)
	}

}

//return top of messages for mocked user
func (r *messageRater) returnMessageTop(message alina.PrivateMessage) {

	text := message.GetText()
	if (strings.Contains(text, "|@") || strings.Contains(text, "|")) && (strings.Contains(text, "[id") || strings.Contains(text, "[club")) {
		var postfix string
		var prefix string
		if strings.Contains(text, "|@") {
			postfix = "|@"
		} else {
			postfix = "|"
		}
		if strings.Contains(text, "[id") {
			prefix = "[id"
		} else {
			prefix = "[club"
		}
		id := string([]rune(message.GetText())[:strings.Index(message.GetText(), postfix)])
		id = string([]rune(id)[strings.Index(id, prefix)+len(prefix):])
		if prefix == "[club" {
			id = "-" + id
		}
		fromId, err := strconv.Atoi(id)
		if err != nil {
			r.logger.Error("cannot convert fromid to int: " + id)
		}

		rate, err := r.storage.GetMessageTop(message.GetPeerId(), fromId)
		if err != nil {
			r.logger.Error("cannot get messages top: " + err.Error())
		}

		top := getTop(5, rate)

		response := ""
		for k, v := range top {
			if k > 0 {
				response += "\r\n\r\n"
			}
			response += "топ " + strconv.Itoa(k+1) + ":\r\n\"" + v + "\""
		}
		if top == nil || len(top) == 0 {
			response = "нет топа по этому пользаку"
		}
		r.alina.GetMessagesApi().SendSimpleMessage(strconv.Itoa(message.GetPeerId()), response)

	}
}

//rate one forwarded message
func (r *messageRater) rateMessage(message alina.PrivateMessage) {
	fwd := message.GetFwdMessages()
	if fwd != nil && len(fwd) == 1 {
		err := r.storage.IncrementMessageRate(message.GetPeerId(), fwd[0].GetFromID(), fwd[0].GetDate(), fwd[0].GetText())
		if err != nil {
			r.logger.Error(err)
		}
	}
}

func (r *messageRater) Jsonize() (*string, error) {
	result := "{}"
	return &result, nil
}

func (r *messageRater) Meet(message alina.PrivateMessage) bool {
	if message.GetText() == "+" && message.GetFwdMessages() != nil {
		return true
	}
	if strings.Contains(strings.ToLower(message.GetText()), "топ перлов") {
		return true
	}

	return false
}

func (r *messageRater) Order() int {
	return 10
}

func (r *messageRater) RateMessage() {

}

func getTop(i int, rate map[string]int) []string {
	top := make([]string, 0)
	temp := make(map[string]struct{})

	for len(top) < i || len(top) < len(rate) {
		tr := -1
		tt := ""
		for text, rate := range rate {
			if _, ok := temp[text]; !ok && rate > tr {
				tr = rate
				tt = text
			}
		}
		if tr == -1 {
			return top
		}
		top = append(top, tt)
		temp[tt] = struct{}{}
	}
	return top
}
