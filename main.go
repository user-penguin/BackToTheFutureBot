package main

import (
	"BackToTheFutureBot/command"
	"BackToTheFutureBot/message"
	"BackToTheFutureBot/reader"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	run()
}

func run() {
	token, err := reader.GetTokenFromFile()
	if err != nil {
		panic(err)
	}
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	log.Println("Bot is start up!")

	for update := range updates {
		// empty message
		if update.Message == nil {
			continue
		}
		text := update.Message.Text
		chatId := update.Message.Chat.ID
		switch command.Command(text) {
		case command.Start:
			_, _ = bot.Send(tgbot.NewMessage(chatId, message.StartMessage()))
		}
	}

}
