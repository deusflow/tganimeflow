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
	return fmt.Sprintf(messages["ua"]["anime_found"], anime.Title, anime.Score)
}

// Структура для разбора ответа от Jikan API
type AnimeData struct {
	Title string  `json:"title"`
	Score float64 `json:"score"`
}

type JikanResponse struct {
	Data []AnimeData `json:"data"`
}

// Создает кнопки выбора языка
func createLanguageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇺🇦 Українська", "lang_ua"),
			tgbotapi.NewInlineKeyboardButtonData("🇺🇸 English", "lang_en"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇩🇰 Dansk", "lang_da"),
		),
	)
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
		"anime_found":   "🎌 Назва: %s\n⭐ Рейтинг: %.1f",
		"lang_changed":  "🌍 Мову змінено на солов'їна! Пощебечемо разом, rebel-чан!",
	},
	"en": {
		"start":         "Hmm... Who dares to disturb the DeusAnimeFlow bot? 💀\n\nAlright... I'm *Anime Finder Bot*, your personal dark guide to the anime world. Write a title, and I'll find it faster than you can say 'Sugoi'.\n\nBut remember... if it's boring anime — I'll snort. 😏\n\nLet's go searching, rebel-chan!",
		"help":          "🌀 You activated STAND *ANIME FINDER*! 🌀\n\nThis bot is created for those who seek their anime destiny. Write anime or manga title — and I, your personal stand, will give you:\n🎯 Title\n📊 Rating\n💥 (in future — genre and description)\n\n💬 Commands worthy of Joestar:\n/start — *Summon the stand!*\n/help — *Call the power of wisdom!*",
		"empty_message": "What's so empty here, for crying out loud? Expand your domain, write anime title and I'll find it! Don't be so lazy, rebel-chan!",
		"api_error":     "Error occurred while searching anime. Try later, rebel-chan.",
		"read_error":    "Error reading API response. Maybe server got tired? Or is Kuromi messing around again?",
		"json_error":    "Error parsing JSON response from API. Maybe server decided to speak its secret dialect?",
		"not_found":     "Anime not found. Try another title, maybe something more EPIC?",
		"anime_found":   "🎌 Title: %s\n⭐ Rating: %.1f",
		"lang_changed":  "🌍 Language changed to English! Now I'll speak with you in English, rebel-chan!",
	},
	"da": {
		"start":         "Hvem tør forstyrre DeusAnimeFlow-botten? 💀\n\nOkay da... Jeg er *Anime Finder Bot*, din personlige mørke guide til anime-verdenen. Skriv en titel, og jeg finder det hurtigere, end du kan sige 'Sugoi'.\n\nMen husk... hvis det er kedelig anime — så fnyster jeg. 😏\n\nLad os søge, rebel-chan!",
		"help":          "🌀 Du har aktiveret STANDEN *ANIME FINDER*! 🌀\n\nDenne bot er skabt til dem, der søger deres anime-skæbne. Skriv titlen på en anime eller manga — og jeg, din personlige stand, vil give dig:\n🎯 Titel\n📊 Bedømmelse\n💥 (senere — genre og beskrivelse)\n\n💬 Kommandoer værdige en Joestar:\n/start — *Påkald standen!*\n/help — *Tilkald visdommens kraft!*",
		"empty_message": "Hvad er så tomt her, altså? Udvid dit domæne og skriv en anime-titel! Vær nu ikke doven, rebel-chan!",
		"api_error":     "Der opstod en fejl under søgning. Prøv igen senere, rebel-chan.",
		"read_error":    "Fejl ved læsning af API-svar. Måske blev serveren træt? Eller leger Kuromi igen?",
		"json_error":    "Fejl ved fortolkning af JSON-svar fra API. Taler serveren sit hemmelige sprog?",
		"not_found":     "Anime ikke fundet. Prøv en anden titel — måske noget mere EPISK?",
		"anime_found":   "🎌 Titel: %s\n⭐ Bedømmelse: %.1f",
		"lang_changed":  "🌍 Sproget er nu ændret til dansk! Klar til at snakke med mig, rebel-chan? Rødgrød med fløde, huh?! 😏🇩🇰",
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
		// Обработка нажатий на inline-кнопки
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "lang_ua" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages["ua"]["lang_changed"])
				bot.Send(msg)
			} else if update.CallbackQuery.Data == "lang_en" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages["en"]["lang_changed"])
				bot.Send(msg)
			} else if update.CallbackQuery.Data == "lang_da" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages["da"]["lang_changed"])
				bot.Send(msg)
			}
		}
		// Обработка текстовых сообщений
		if update.Message != nil {
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Оберіть мову / Choose language / Vælg sprog:")
				msg.ReplyMarkup = createLanguageKeyboard()
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
}
