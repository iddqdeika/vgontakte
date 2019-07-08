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
		runes := []byte(message.GetText())
		index := strings.Index(message.GetText(), postfix)
		if runes != nil && index > 1 {

		}
		id := string([]byte(message.GetText())[:strings.Index(message.GetText(), postfix)])
		id = string([]byte(id)[strings.Index(id, prefix)+len(prefix):])
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

		if top == nil || len(top) == 0 {
			r.alina.GetMessagesApi().SendSimpleMessage(strconv.Itoa(message.GetPeerId()), "нет топа по пользователю")
		} else {
			for k, v := range top {
				response := "топ " + strconv.Itoa(k+1) + ":\r\nтекст\"" + v.GetText() + "\""
				r.alina.GetMessagesApi().SendMessageWithAttachment(strconv.Itoa(message.GetPeerId()), response, v.GetAttachmentTokensList())
			}
		}
	}
}

//rate one forwarded message
func (r *messageRater) rateMessage(message alina.PrivateMessage) {
	fwd := message.GetFwdMessages()
	if fwd != nil && len(fwd) == 1 {
		tokens := make([]string, 0)
		for _, v := range fwd[0].GetAttachments() {
			if v.IsMedia() {
				token, err := v.GetPrivateMessageToken()
				if err != nil {
					r.logger.Info(err)
				} else {
					tokens = append(tokens, token)
				}
			}
		}
		err := r.storage.IncrementMessageRate(message.GetPeerId(), fwd[0].GetFromID(), fwd[0].GetDate(), fwd[0].GetText(), tokens)
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

func getTop(i int, rate map[vgontakte.RaterMessage]int) []vgontakte.RaterMessage {
	top := make([]vgontakte.RaterMessage, 0)
	temp := make(map[string]struct{})

	for len(top) < i || len(top) < len(rate) {
		tr := -1
		var tt vgontakte.RaterMessage
		for msg, rate := range rate {
			if _, ok := temp[msg.GetDate()]; !ok && rate > tr {
				tr = rate
				tt = msg
			}
		}
		if tr == -1 {
			return top
		}
		top = append(top, tt)
		temp[tt.GetDate()] = struct{}{}
	}
	return top
}
