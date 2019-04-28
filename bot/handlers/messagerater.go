package handlers

import (
	"alina/alina"
	"sort"
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
	}

	return &handler, nil
}

func (c *messageRaterCreator) ParseHandler(data *string, alina alina.Alina, bot vgontakte.Bot) (vgontakte.MessageHandler, error) {
	handler := messageRater{
		storage: bot.GetStrorage(),
		alina:   alina,
	}

	return &handler, nil
}

type messageRater struct {
	storage vgontakte.Storage
	alina   alina.Alina
	peerId  int
}

func (r *messageRater) Handle(message alina.PrivateMessage, err error) {
	if message.GetText() == "+" && message.GetFwdMessages() != nil {
		fwd := message.GetFwdMessages()
		if fwd != nil && len(fwd) == 1 {
			path := strconv.Itoa(message.GetPeerId()) + ".rates." + strconv.Itoa(fwd[0].GetFromID())
			ratepath := path + ".rate" + strconv.Itoa(fwd[0].GetDate())
			datapath := path + ".data" + strconv.Itoa(fwd[0].GetDate())

			rate := 1
			val, err := r.storage.Get(ratepath)
			if err == nil {
				oldrate, err := strconv.Atoi(string(val))
				if err == nil {
					rate += oldrate
				}
			}
			r.storage.Update(ratepath, strconv.Itoa(rate))
			r.storage.Update(datapath, fwd[0].GetText())
		}
	}
	if strings.Contains(strings.ToLower(message.GetText()), "топ перлов") {
		text := message.GetText()
		if (strings.Contains(text, "|@")) && strings.Contains(text, "[id") {
			id := string([]rune(message.GetText())[:strings.Index(message.GetText(), "|@")])
			id = string([]rune(id)[strings.Index(id, "[id")+3:])
			ratepath := strconv.Itoa(message.GetPeerId()) + ".rates." + id
			messages := make(map[string]int)
			r.storage.Iterate(ratepath, func(k, v []byte) error {
				if strings.Contains(string(k), "rate") && strings.Index(string(k), "rate") == 0 {
					i, err := strconv.Atoi(string(v))
					if err == nil {
						messages[string([]rune(string(k))[4:])] = i
					}
				}
				return nil
			})
			dates := make([]string, 0)
			for k, _ := range messages {
				dates = append(dates, k)
			}

			sort.Slice(dates, func(i, j int) bool {
				return messages[dates[i]] > messages[dates[j]]
			})

			results := make([]string, 0)
			for index := 0; index < len(dates) && index < 3; index++ {
				data, err := r.storage.Get(ratepath + ".data" + dates[index])
				if err == nil {
					results = append(results, string(data))
				}
			}

			response := ""
			for k, v := range results {
				if k > 0 {
					response += "\r\n\r\n"
				}
				response += "топ " + strconv.Itoa(k+1) + ":\r\n\"" + v + "\""
			}
			r.alina.GetMessagesApi().SendSimpleMessage(strconv.Itoa(message.GetPeerId()), response)

		}
	}

}

func test(k, v []byte) {
	kk := string(k)
	vv := string(v)
	println(kk + "-" + vv)
}

func (r *messageRater) Jsonize() (*string, error) {
	panic("implement me")
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
	panic("implement me")
}

func (r *messageRater) RateMessage() {

}
