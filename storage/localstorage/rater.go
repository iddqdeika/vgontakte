package localstorage

func getNewPeerMessageRates(peerId int) peerMessageRates {
	res := peerMessageRates{}
	res.PeerId = peerId
	res.Users = make(map[int]*userMessageRates)
	return res
}

func getNewUserMessageRates(fromId int) *userMessageRates {
	res := userMessageRates{}
	res.FromId = fromId
	res.Messages = make(map[int]*message)
	return &res
}

type peerMessageRates struct {
	PeerId int                       `json:"peer_id"`
	Users  map[int]*userMessageRates `json:"users"`
}

func (p *peerMessageRates) getUserRate(fromId int) *userMessageRates {
	if val, ok := p.Users[fromId]; ok {
		return val
	}
	p.Users[fromId] = getNewUserMessageRates(fromId)
	return p.Users[fromId]
}

func (p *peerMessageRates) incrementRate(fromId int, date int, text string) {
	p.getUserRate(fromId).getMessage(date, text).Rate++

}

type userMessageRates struct {
	FromId   int              `json:"from_id"`
	Messages map[int]*message `json:"messages"`
}

func (u *userMessageRates) getMessage(date int, text string) *message {
	if val, ok := u.Messages[date]; ok {
		return val
	}
	u.Messages[date] = &message{Text: text}
	return u.Messages[date]
}

type message struct {
	Text string `json:"text"`
	Rate int    `json:"rate"`
}
