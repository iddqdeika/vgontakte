package vgontakte

import (
	"alina/alina"
)

type Bot interface {
	Init(alina alina.Alina, storage Storage, logger Logger) error
	Run()
	GetPrivateMessageDispatcher() (PrivateMessageDispatcher, error)
	GetStrorage() Storage
	GetLogger() Logger
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Config interface {
	GetInt(name string) (int, error)
	GetString(name string) (string, error)
}

type Storage interface {
	IncrementMessageRate(peerId int, fromId int, fwdDate int, messageText string) error
	GetMessageTop(peerId int, fromId int) (map[string]int, error)
	RegisterPeer(peerId int) error
	CheckPeerRegistration(peerId int) bool
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
