package dispatcher

import (
	"alina/alina"
	"sort"
	"sync"
	"vgontakte/vgontakte"
)

func NewPrivateMessageDispatcher() vgontakte.PrivateMessageDispatcher {
	return &privateMessageDispatcher{}
}

type privateMessageDispatcher struct {
	m        sync.RWMutex
	handlers []vgontakte.MessageHandler
}

func (d *privateMessageDispatcher) DispatchMessage(message alina.PrivateMessage, e error) {
	go func() {
		handlers := d.SafelyGetHandlers()
		for _, handler := range handlers {
			if handler.Meet(message) {
				handler.Handle(message, e)
				return
			}
		}
	}()
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
