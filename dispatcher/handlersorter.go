package dispatcher

import "vgontakte/vgontakte"

type MessageHandlersSorter []vgontakte.MessageHandler

func (a MessageHandlersSorter) Len() int           { return len(a) }
func (a MessageHandlersSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MessageHandlersSorter) Less(i, j int) bool { return a[i].Order() < a[j].Order() }
