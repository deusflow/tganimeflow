package bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

// тут будет запрос к API аниме и манги
func searchAnime(query string) string {
	url := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=1", query)
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching data from Jikan API:", err)
		return messages["ua"]["api_error"]
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return messages["ua"]["read_error"]
	}
	var result JikanResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return messages["ua"]["json_error"]
	}
	if len(result.Data) == 0 {
		return messages["ua"]["not_found"]
	}

	anime := result.Data[0]
	return fmt.Sprintf("Название: %s\nРейтинг: %.1f\n", anime.Title, anime.Score)
}

// Структура для разбора ответа от Jikan API
type AnimeData struct {
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

type JikanResponse struct {
	Data []AnimeData `json:"data"`
}

// Тексты на разных языках
var messages = map[string]map[string]string{
	"ua": {
		"start":         "Хмм... Кто тут такой смелый, чтобы потревожить DeusAnimeFlow бота? 💀\n\nНу ладно... Я — *Anime Finder Bot*, твой личный тёмный проводник в мире аниме. Напиши название, и я найду его быстрее, чем ты скажешь 'Sugoi'.\n\nНо учти… если это скучное аниме — я фыркну. 😏\n\nПогнали искать, rebel-чан!",
		"help":          "🌀 Ти активував СТЕНД *ANIME FINDER*! 🌀\n\nЭтот бот создан для тех, кто ищет своё аниме-предназначение. Напиши название аниме или манги — и Я, твой персональный стенд, выдам тебе:\n🎯 Название\n📊 Рейтинг\n💥 (в будущем — жанр и описание)\n\n💬 Команды, достойные Джостара:\n/start — *Призови стенд!*\n/help — *Вызови силу мудрости!*",
		"empty_message": "А щож тут так пусто, трясця богу? Розширь свої володіння, напиши назву аніме і я його знайду! Не будь таким ледащим, rebel-чан!",
		"api_error":     "Сталася помилка при пошуку аніме. Спробуй пізніше, rebel-чан.",
		"read_error":    "Помилка читання відповіді від API. Може, сервер втомився? Чи це Kuromi знову шалить?",
		"json_error":    "Помилка розбору JSON відповіді від API. Може, сервер вирішив поговорити на своєму таємному діалекті?",
		"not_found":     "Аніме не знайдено. Спробуй іншу назву, може щось більш EPIC?",
	},
}

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
			fmt.Println("Message Received:", update.Message.Text)
		}
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["ua"]["start"])
			bot.Send(msg)

			// Отправляем приветственное сообщение и берем его из messages map
		} else if update.Message.IsCommand() && update.Message.Command() == "help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["ua"]["help"])
			bot.Send(msg)

		} else if !update.Message.IsCommand() {
			if update.Message.Text == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["ua"]["empty_message"])
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, searchAnime(update.Message.Text))
				bot.Send(msg)
			}
		}
	}
}
