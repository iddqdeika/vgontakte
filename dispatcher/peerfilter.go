package dispatcher

import (
	"alina/alina"
	"vgontakte/vgontakte"
)

type peerFilter struct {
	storage vgontakte.Storage
	alina   alina.Alina
	logger  vgontakte.Logger
}
