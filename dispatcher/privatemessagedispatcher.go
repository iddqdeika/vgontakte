package dispatcher

import (
	"alina/alina"
	"sort"
	"sync"
	"vgontakte/config"
	"vgontakte/vgontakte"
)

func NewPrivateMessageDispatcher(storage vgontakte.Storage) vgontakte.PrivateMessageDispatcher {
	return &privateMessageDispatcher{storage: storage}
}

type privateMessageDispatcher struct {
	storage  vgontakte.Storage
	m        sync.RWMutex
	handlers []vgontakte.MessageHandler
}

func (d *privateMessageDispatcher) DispatchMessage(message alina.PrivateMessage, e error) {
	go func() {
		d.dispatchMessage(message, e)
	}()
}

func (d *privateMessageDispatcher) dispatchMessage(message alina.PrivateMessage, e error) {
	if !d.filterMessage(message.GetPeerId()) {
		return
	}
	handlers := d.SafelyGetHandlers()
	for _, handler := range handlers {
		if handler.Meet(message) {
			handler.Handle(message, e)
			return
		}
	}
}

func (d *privateMessageDispatcher) SafelyGetHandlers() []vgontakte.MessageHandler {
	var handlers = make([]vgontakte.MessageHandler, 0)
	d.m.RLock()
	defer d.m.RUnlock()
	for _, handler := range d.handlers {
		handlers = append(handlers, handler)
	}
	return handlers
}

func (d *privateMessageDispatcher) SafelyAddHandler(handler vgontakte.MessageHandler) {
	d.m.Lock()
	defer d.m.Unlock()
	d.handlers = append(d.handlers, handler)
	sort.Sort(MessageHandlersSorter(d.handlers))
}

func (d *privateMessageDispatcher) filterMessage(peerId int) bool {
	owner, err := config.GetCommandLineArgsConfig().GetInt("ownerpeer")
	if err == nil && owner == peerId {
		return true
	}
	return d.storage.CheckPeerRegistration(peerId)
}
