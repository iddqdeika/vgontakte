package bot

import (
	"alina/alina"
	"fmt"
	"vgontakte/bot/handlers"
	"vgontakte/defaultlogger"
	"vgontakte/dispatcher"
	"vgontakte/storage/localstorage"
	"vgontakte/vgontakte"
)

func NewBot() vgontakte.Bot {
	return &vbot{}
}

type vbot struct {
	alina             alina.Alina
	storage           vgontakte.Storage
	logger            vgontakte.Logger
	messageDispatcher vgontakte.PrivateMessageDispatcher
}

func (b *vbot) GetLogger() vgontakte.Logger {
	return b.logger
}

func (b *vbot) GetStrorage() vgontakte.Storage {
	return b.storage
}

func (b *vbot) GetPrivateMessageDispatcher() (vgontakte.PrivateMessageDispatcher, error) {
	if b.messageDispatcher != nil {
		return b.messageDispatcher, nil
	}
	return nil, fmt.Errorf("bot has no messagedispatcher yet")
}

func (b *vbot) Init(alina alina.Alina, storage vgontakte.Storage, logger vgontakte.Logger) error {
	if alina == nil {
		return fmt.Errorf("alina is nil, we need alina!")
	}
	b.alina = alina

	if storage == nil {
		storage = localstorage.GetLocalStorage("default")
	}
	b.storage = storage

	if logger == nil {
		logger = defaultlogger.Logger
	}
	b.logger = logger

	b.messageDispatcher = dispatcher.NewPrivateMessageDispatcher(b.storage)

	mh := b.initializeMessageHandlers()
	for _, h := range mh {
		b.messageDispatcher.SafelyAddHandler(h)
	}

	alina.AddMessageHandler(b.messageDispatcher.DispatchMessage)
	return nil
}

func (b *vbot) Run() {
	b.alina.Run()
}

func (b *vbot) initializeMessageHandlers() []vgontakte.MessageHandler {
	result := make([]vgontakte.MessageHandler, 0)

	handler, err := b.getHandlerBuCreatorName(vgontakte.RateMessageHandler)
	if err != nil {
		b.logger.Error(err)
	} else {
		result = append(result, handler)
	}

	handler, err = b.getHandlerBuCreatorName(vgontakte.PeerRegisterHandler)
	if err != nil {
		b.logger.Error(err)
	} else {
		result = append(result, handler)
	}

	return result
}

func (b *vbot) getHandlerBuCreatorName(creator vgontakte.HandlerType) (vgontakte.MessageHandler, error) {
	creater, err := handlers.GetHandlerCreator(creator)
	if err != nil {
		return nil, err
	}
	var i int
	params := map[string]interface{}{
		"peer_id": i,
	}
	handler, err := creater.CreateHandler(params, b.alina, b)
	if err != nil {
		panic(err)
	}
	return handler, nil
}
