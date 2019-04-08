package privatemessage

import (
	"alina/definitions"
	"sort"
	"sync"
)

type MessageHandler interface {
	Order() int
	Meet(message definitions.PrivateMessage) bool
	Handle(definitions.PrivateMessage, error)
}

type MessageHandlers []MessageHandler

func (a MessageHandlers) Len() int           { return len(a) }
func (a MessageHandlers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MessageHandlers) Less(i, j int) bool { return a[i].Order() < a[j].Order() }

type PrivateMessageDispatcher struct {
	m        sync.RWMutex
	handlers []MessageHandler
}

func (d *PrivateMessageDispatcher) Dispatch(message definitions.PrivateMessage, e error) {
	d.m.RLock()
	for _, handler := range d.handlers {
		if handler.Meet(message) {
			handler.Handle(message, e)
			return
		}
	}
	defer d.m.RUnlock()
}

func (d *PrivateMessageDispatcher) AddHandler(handler MessageHandler) {
	d.m.Lock()
	defer d.m.Unlock()
	d.handlers = append(d.handlers, handler)
	sort.Sort(MessageHandlers(d.handlers))
}
