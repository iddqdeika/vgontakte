package vgontakte

import (
	"alina/alina"
)

type Bot interface {
	Init(alina alina.Alina, storage Storage, logger Logger) error
	Run()
	GetPrivateMessageDispatcher() (PrivateMessageDispatcher, error)
	GetStrorage() Storage
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Storage interface {
	IncrementMessageRate(peerId int, fromId int, fwdDate int) error
	GetMessageTop(peerId int, fromId int) (map[string]int, error)
}

type MessageHandler interface {
	Order() int
	Meet(message alina.PrivateMessage) bool
	Handle(alina.PrivateMessage, error)
	Jsonize() (*string, error)
}

type HandlerCreator interface {
	CreateHandler(params map[string]interface{}, alina alina.Alina, bot Bot) (MessageHandler, error)
	ParseHandler(json *string, alina alina.Alina, bot Bot) (MessageHandler, error)
}

type PrivateMessageDispatcher interface {
	DispatchMessage(message alina.PrivateMessage, e error)
	SafelyGetHandlers() []MessageHandler
	SafelyAddHandler(handler MessageHandler)
}
