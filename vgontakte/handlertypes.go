package vgontakte

type HandlerType string

const (
	PeerRegisterHandler HandlerType = "peer_register"
	RateMessageHandler  HandlerType = "rate_message"
)
