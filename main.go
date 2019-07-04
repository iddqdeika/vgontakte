package main

import (
	"alina/alina"
	"alina/alinafactory"
	"alina/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"vgontakte/bot"
	"vgontakte/defaultlogger"
	"vgontakte/storage/localstorage"
)

const (
	cfgName = "config.cfg"
)

func main() {

	storage := localstorage.GetLocalStorage("testdb")
	alina := initAlina()
	vkbot := bot.NewBot()

	err := vkbot.Init(alina, storage, defaultlogger.Logger)
	if err != nil {
		panic(err)
	}

	vkbot.Run()

}

func initAlina() alina.Alina {
	cfg := &struct {
		AccessToken string `json:"access_token"`
		GroupId     string `json:"group_id"`
		LongPollInt int    `json:"long_poll_interval"`
	}{}

	cfgData, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(cfgData, cfg)
	if err != nil || len(cfg.AccessToken) == 0 {
		panic(err)
	}

	logger.InitDefaultLogger()
	logger := logger.DefaultLogger
	al, err := alinafactory.New(cfg.AccessToken, "5.85", cfg.GroupId, logger, cfg.LongPollInt)
	if err != nil {
		logger.Error(fmt.Sprintf("fatal error during Alina initialization: ", err))
		panic("cant create alina")
	}

	err = al.Init()
	if err != nil {
		logger.Error(fmt.Sprintf("fatal error during Alina initialization: %v", err))
		panic("cant create alina")
	}

	return al
}
