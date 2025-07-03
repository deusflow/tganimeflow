package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Loading error .env file")
	}
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_TOKEN not found in .env")
	}
	fmt.Println("Bot started")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Bot authorized How:", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			fmt.Println("Massage Received:", update.Message.Text)
		}
	}

}
