package handlers

import (
	"fmt"
	"vgontakte/vgontakte"
)

var handlersMap = map[vgontakte.HandlerType]vgontakte.HandlerCreator{
	vgontakte.RateMessageHandler:  &messageRaterCreator{},
	vgontakte.PeerRegisterHandler: &peerRegistratorCreator{},
}

func GetHandlerCreator(t vgontakte.HandlerType) (vgontakte.HandlerCreator, error) {
	if val, ok := handlersMap[t]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("cant find any handlecreator for: %v", t)
}
