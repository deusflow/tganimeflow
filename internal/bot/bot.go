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

const jikanBaseURL = "https://api.jikan.moe/v4"

// Константы для команд бота
const (
	cmdStart  = "start"
	cmdHelp   = "help"
	cmdRandom = "random"
	cmdTop    = "top"
	cmdDonate = "donate"
	cmdStats  = "stats"
)

// Универсальная функция
func fetchAndUnmarshal(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	return json.Unmarshal(body, target)
}

// Централизованная обработка ошибок API
func handleAPIError(lang, errType string) AnimeData {
	return AnimeData{Title: messages[lang][errType]}
}

// Централизованное логирование запросов
func logRequest(operation string, err error) {
	if err != nil {
		log.Printf("Error in %s: %v", operation, err)
	}
}

// тут будет запрос к API аниме и манги
func searchAnime(query string, lang string) AnimeData {
	url := fmt.Sprintf("%s/anime?q=%s&limit=1", jikanBaseURL, query)

	var result JikanResponse
	if err := fetchAndUnmarshal(url, &result); err != nil {
		logRequest("searchAnime", err)
		return handleAPIError(lang, "api_error")
	}

	if len(result.Data) == 0 {
		return handleAPIError(lang, "not_found")
	}

	return result.Data[0]
}

func getRandomAnime(lang string) AnimeData {
	url := fmt.Sprintf("%s/random/anime", jikanBaseURL)

	var result RandomAnimeResponse
	if err := fetchAndUnmarshal(url, &result); err != nil {
		logRequest("getRandomAnime", err)
		return handleAPIError(lang, "api_error")
	}

	return result.Data
}

// AnimeData Структура для разбора ответа от Jikan API
type AnimeData struct {
	Title    string  `json:"title"`
	Score    float64 `json:"score"`
	Synopsis string  `json:"synopsis"`
	Episodes int     `json:"episodes"`
	Status   string  `json:"status"`
	Genres   []Genre `json:"genres"`
	Images   Images  `json:"images"`
}

// Структура для хранения аналитики
type Analytics struct {
	TotalUsers    int            `json:"total_users"`
	CommandsUsed  map[string]int `json:"commands_used"`
	LanguagesUsed map[string]int `json:"languages_used"`
}

type Genre struct {
	Name string `json:"name"`
}

type Images struct {
	JPG ImageData `json:"jpg"`
}

type ImageData struct {
	LargeImageURL string `json:"large_image_url"`
}
type JikanResponse struct {
	Data []AnimeData `json:"data"`
}

type RandomAnimeResponse struct {
	Data AnimeData `json:"data"`
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

// Создает кнопки быстрых действий на нужном языке
func createQuickActionsKeyboard(lang string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_random"], "action_random"),
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_top"], "action_top"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_search"], "action_search"),
		),
	)
}

// Создает кнопки для донатов
func createDonateKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("💳 PayPal", "https://paypal.me/deusflowro"),
			tgbotapi.NewInlineKeyboardButtonURL("📱 MobilePay", "https://qr.mobilepay.dk/box/d017b43a-052e-4884-8fd6-851349b234a2/pay-in"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❤️ Дякую!", "donate_thanks"),
		),
	)
}

func getTopAnime(lang string) string {
	url := fmt.Sprintf("%s/top/anime?limit=5", jikanBaseURL)

	var result JikanResponse
	if err := fetchAndUnmarshal(url, &result); err != nil {
		logRequest("getTopAnime", err)
		return messages[lang]["api_error"]
	}

	if len(result.Data) == 0 {
		return messages[lang]["not_found"]
	}

	topAnime := messages[lang]["top_anime"] + "\n\n"
	for i, anime := range result.Data {
		topAnime += fmt.Sprintf("%d. %s - ⭐ %.1f\n", i+1, anime.Title, anime.Score)
	}
	return topAnime
}

func formatAnimeDetails(anime AnimeData, lang string) string {
	// Форматирую жанры в строку
	genresText := ""
	for i, genre := range anime.Genres {
		if i > 0 {
			genresText += ", "
		}
		genresText += genre.Name
	}

	//кол-во серий
	episodesText := "?" // если серий нет, то будет "?"
	if anime.Episodes > 0 {
		episodesText = fmt.Sprintf("%d", anime.Episodes)
	}

	//описание ограничиваем длину
	synopsis := anime.Synopsis
	if len(synopsis) > 200 {
		synopsis = synopsis[:200] + "..."
	}

	return fmt.Sprintf(
		"🎌 %s\n⭐ %.1f\n📺 %s серій\n📊 %s\n🎭 %s\n\n📝 %s",
		anime.Title,
		anime.Score,
		episodesText,
		anime.Status,
		genresText,
		synopsis,
	)
}

// Отправляет аниме с картинкой
func sendAnimeWithPhoto(bot *tgbotapi.BotAPI, chatID int64, anime AnimeData, lang string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	caption := formatAnimeDetails(anime, lang)

	if anime.Images.JPG.LargeImageURL != "" {
		// Отправляем фото с описанием
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(anime.Images.JPG.LargeImageURL))
		photo.Caption = caption
		if keyboard != nil {
			photo.ReplyMarkup = *keyboard
		}
		bot.Send(photo)
	} else {
		// Если нет картинки, отправляем обычное сообщение
		msg := tgbotapi.NewMessage(chatID, caption)
		if keyboard != nil {
			msg.ReplyMarkup = *keyboard
		}
		bot.Send(msg)
	}
}

// Логирует действие пользователя для аналитики
func logUserAction(userID int64, action string, lang string) {
	// Проверяем, новый ли это пользователь
	if !knownUsers[userID] {
		knownUsers[userID] = true
		botAnalytics.TotalUsers++
		fmt.Printf("📊 Новый пользователь! Всего: %d\n", botAnalytics.TotalUsers)
	}

	// Считаем использование команд
	botAnalytics.CommandsUsed[action]++

	// Считаем использование языков
	botAnalytics.LanguagesUsed[lang]++

	fmt.Printf("📊 Действие: %s, Язык: %s, Всего команд '%s': %d\n",
		action, lang, action, botAnalytics.CommandsUsed[action])
}

// Тексты на разных языках
var messages = map[string]map[string]string{
	"ua": {
		"start":          "Хмм... Хіто тут такий сміливий, щоб відволікати могутнього DeusAnimeFlow бота? 💀\n\nНу доообре... Я — *Anime Finder Bot*, твій особистий таємний провідник у пітьму. Напиши назву, та я знайду його швидше, ніж ти вигукнеш 'Sugoi'.\n\n на нудні аніме - фиркаю 😏\n\n Погнали, rebel-чан!",
		"help":           "🌀 Ти активував СТЕНД *ANIME FINDER*! 🌀\n\nЭтот бот создан для тех, кто ищет своё аниме-предназначение. Напиши название аниме или манги — и Я, твой персональный стенд, выдам тебе:\n🎯 Название\n📊 Рейтинг\n💥 (в будущем — жанр и описание)\n\n💬 Команды, достойные Джостара:\n/start — *Призови стенд!*\n/help — *Вызови силу мудрости!*",
		"empty_message":  "А щож тут так пусто, трясця богу? Розширь свої володіння, напиши назву аніме і я його знайду! Не будь таким ледащим, rebel-чан!",
		"api_error":      "Сталася помилка при пошуку аніме. Спробуй пізніше, rebel-чан.",
		"read_error":     "Помилка читання відповіді від API. Може, сервер втомився? Чи це Kuromi знову шалить?",
		"json_error":     "Помилка розбору JSON відповіді від API. Може, сервер вирішив поговорити на своєму таємному діалекті?",
		"not_found":      "Аніме не знайдено. Спробуй іншу назву, може щось більш EPIC?",
		"anime_found":    "🎌 Назва: %s\n⭐ Рейтинг: %.1f",
		"lang_changed":   "🌍 Мову змінено на солов'їна! Пощебечемо разом, rebel-чан!",
		"random_anime":   "🎲 Видкусіньке аніме для тебе мій пупсику:",
		"top_anime":      "🏆 Топ аніме:",
		"donate_message": "💖 Подобається бот? Підтримай розробника!\n\n🌟 Твоя підтримка допомагає розвивати бота та додавати нові функції!\n\nОбери зручний спосіб:",
		"donate_thanks":  "💖 Дякую за підтримку, rebel-чан! Ти крутий! 🔥",
		"btn_random":     "🎲 Випадкове",
		"btn_top":        "🏆 Топ",
		"btn_search":     "🔄 Новий пошук",
		"stats":          "📊 Статистика",
	},
	"en": {
		"start":          "Hmm... Who dares to disturb the DeusAnimeFlow bot? 💀\n\nAlright... I'm *Anime Finder Bot*, your personal dark guide to the anime world. Write a title, and I'll find it faster than you can say 'Sugoi'.\n\nBut remember... if it's boring anime — I'll snort. 😏\n\nLet's go searching, rebel-chan!",
		"help":           "🌀 You activated STAND *ANIME FINDER*! 🌀\n\nThis bot is created for those who seek their anime destiny. Write anime or manga title — and I, your personal stand, will give you:\n🎯 Title\n📊 Rating\n💥 (in future — genre and description)\n\n💬 Commands worthy of Joestar:\n/start — *Summon the stand!*\n/help — *Call the power of wisdom!*",
		"empty_message":  "What's so empty here, for crying out loud? Expand your domain, write anime title and I'll find it! Don't be so lazy, rebel-chan!",
		"api_error":      "Error occurred while searching anime. Try later, rebel-chan.",
		"read_error":     "Error reading API response. Maybe server got tired? Or is Kuromi messing around again?",
		"json_error":     "Error parsing JSON response from API. Maybe server decided to speak its secret dialect?",
		"not_found":      "Anime not found. Try another title, maybe something more EPIC?",
		"anime_found":    "🎌 Title: %s\n⭐ Rating: %.1f",
		"lang_changed":   "🌍 Language changed to English! Now I'll speak with you in English, rebel-chan!",
		"random_anime":   "🎲 Random anime for you:",
		"top_anime":      "🏆 Top anime:",
		"donate_message": "💖 Like the bot? Support the creator!\n\n🌟 Your power-up helps us grow and unlock new features!\n\nChoose your favorite way to support:",
		"donate_thanks":  "💖 Arigato for your support, rebel-chan! You're awesome! 🔥",
		"btn_random":     "🎲 Random",
		"btn_top":        "🏆 Top",
		"btn_search":     "🔄 New search",
		"stats":          "📊 Statistics",
	},
	"da": {
		"start":          "Hvem tør forstyrre DeusAnimeFlow-botten? 💀\n\nOkay da... Jeg er *Anime Finder Bot*, din personlige mørke guide til anime-verdenen. Skriv en titel, og jeg finder det hurtigere, end du kan sige 'Sugoi'.\n\nMen husk... hvis det er kedelig anime — så fnyster jeg. 😏\n\nLad os søge, rebel-chan!",
		"help":           "🌀 Du har aktiveret STANDEN *ANIME FINDER*! 🌀\n\nDenne bot er skabt til dem, der søger deres anime-skæbne. Skriv titlen på en anime eller manga — og jeg, din personlige stand, vil give dig:\n🎯 Titel\n📊 Bedømmelse\n💥 (senere — genre og beskrivelse)\n\n💬 Kommandoer værdige en Joestar:\n/start — *Påkald standen!*\n/help — *Tilkald visdommens kraft!*",
		"empty_message":  "Hvad er så tomt her, altså? Udvid dit domæne og skriv en anime-titel! Vær nu ikke doven, rebel-chan!",
		"api_error":      "Der opstod en fejl under søgning. Prøv igen senere, rebel-chan.",
		"read_error":     "Fejl ved læsning af API-svar. Måske blev serveren træt? Eller leger Kuromi igen?",
		"json_error":     "Fejl ved fortolkning af JSON-svar fra API. Taler serveren sit hemmelige sprog?",
		"not_found":      "Anime ikke fundet. Prøv en anden titel — måske noget mere EPISK?",
		"anime_found":    "🎌 Titel: %s\n⭐ Bedømmelse: %.1f",
		"lang_changed":   "🌍 Sproget er nu ændret til dansk! Klar til at snakke med mig, rebel-chan? Rødgrød med fløde, huh?! 😏🇩🇰",
		"random_anime":   "🎲 Tilfældig anime til dig:",
		"top_anime":      "🏆 Top anime:",
		"donate_message": "💖 Kan du lide botten? Støt skaberen!\n\n🌟 Din energi hjælper os med at vokse og få nye funktioner!\n\nVælg den måde, du vil støtte på:",
		"donate_thanks":  "💖 Tak for støtten, rebel-chan! Du er mega sej! 🔥",
		"btn_random":     "🎲 Tilfældig",
		"btn_top":        "🏆 Top",
		"btn_search":     "🔄 Ny søgning",
		"stats":          "📊 Statistik",
	},
}

// Глобальная аналитика
var botAnalytics = Analytics{
	TotalUsers:    0,
	CommandsUsed:  make(map[string]int),
	LanguagesUsed: make(map[string]int),
}

var knownUsers = make(map[int64]bool) // Для отслеживания уникальных пользователей

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
	userLangs := make(map[int64]string) // userID -> выбр��нный язык

	for update := range updates {
		if update.Message != nil {
			// Проверяем, если это команда для смены языка
			fmt.Println("Message Received:", update.Message.Text)

			// Получаем общие данные ОДИН РАЗ для всех команд
			userID := update.Message.From.ID
			chatID := update.Message.Chat.ID
			lang := userLangs[userID]
			if lang == "" {
				lang = "ua" // По умолчанию украинский
			}

			var responseText string
			var keyboard *tgbotapi.InlineKeyboardMarkup

			if update.Message.IsCommand() && update.Message.Command() == cmdStart {
				logUserAction(userID, "start", lang)
				responseText = messages[lang]["start"]
				languageKeyboard := createLanguageKeyboard()
				keyboard = &languageKeyboard

			} else if update.Message.IsCommand() && update.Message.Command() == cmdHelp {
				logUserAction(userID, "help", lang)
				responseText = messages[lang]["help"]
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			} else if update.Message.IsCommand() && update.Message.Command() == cmdRandom {
				logUserAction(userID, "random", lang)
				anime := getRandomAnime(lang)
				quickKeyboard := createQuickActionsKeyboard(lang)
				sendAnimeWithPhoto(bot, chatID, anime, lang, &quickKeyboard)
				continue

			} else if update.Message.IsCommand() && update.Message.Command() == cmdTop {
				logUserAction(userID, "top", lang)
				responseText = getTopAnime(lang)
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			} else if update.Message.IsCommand() && update.Message.Command() == cmdDonate {
				logUserAction(userID, "donate", lang)
				responseText = messages[lang]["donate_message"]
				donateKeyboard := createDonateKeyboard()
				keyboard = &donateKeyboard

			} else if update.Message.IsCommand() && update.Message.Command() == cmdStats {
				logUserAction(userID, "stats", lang)
				statsText := fmt.Sprintf("📊 СТАТИСТИКА БОТА:\n\n👥 Всего пользователей: %d\n\n📈 Популярные команды:\n", botAnalytics.TotalUsers)

				for command, count := range botAnalytics.CommandsUsed {
					statsText += fmt.Sprintf("• %s: %d раз\n", command, count)
				}

				statsText += "\n🌍 Языки:\n"
				for language, count := range botAnalytics.LanguagesUsed {
					statsText += fmt.Sprintf("• %s: %d раз\n", language, count)
				}

				responseText = statsText

			} else if !update.Message.IsCommand() {
				if update.Message.Text == "" {
					responseText = messages[lang]["empty_message"]
				} else {
					logUserAction(userID, "search", lang)
					anime := searchAnime(update.Message.Text, lang)
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, anime, lang, &quickKeyboard)
					continue
				}
			}

			// Единая точка отправки
			msg := tgbotapi.NewMessage(chatID, responseText)
			if keyboard != nil {
				msg.ReplyMarkup = *keyboard
			}
			bot.Send(msg)
		}

		// Обработка нажатий на inline-кнопки
		if update.CallbackQuery != nil {
			userID := update.CallbackQuery.From.ID
			chatID := update.CallbackQuery.Message.Chat.ID

			lang := userLangs[userID]
			if lang == "" {
				lang = "ua"
			}

			var responseText string
			var withKeyboard bool

			switch update.CallbackQuery.Data {
			case "lang_ua":
				logUserAction(userID, "lang_change", "ua")
				userLangs[userID] = "ua"
				responseText = messages["ua"]["lang_changed"] + "\n" + messages["ua"]["start"]

			case "lang_en":
				logUserAction(userID, "lang_change", "en")
				userLangs[userID] = "en"
				responseText = messages["en"]["lang_changed"] + "\n" + messages["en"]["start"]

			case "lang_da":
				logUserAction(userID, "lang_change", "da")
				userLangs[userID] = "da"
				responseText = messages["da"]["lang_changed"] + "\n" + messages["da"]["start"]

			case "action_random":
				logUserAction(userID, "random", lang)
				anime := getRandomAnime(lang)
				quickKeyboard := createQuickActionsKeyboard(lang)
				sendAnimeWithPhoto(bot, chatID, anime, lang, &quickKeyboard)
				continue

			case "action_top":
				logUserAction(userID, "top", lang)
				responseText = getTopAnime(lang)
				withKeyboard = true

			case "action_search":
				logUserAction(userID, "search_help", lang)
				responseText = messages[lang]["empty_message"]

			case "donate_thanks":
				logUserAction(userID, "donate_thanks", lang)
				responseText = messages[lang]["donate_thanks"]
			}

			msg := tgbotapi.NewMessage(chatID, responseText)

			if withKeyboard {
				msg.ReplyMarkup = createQuickActionsKeyboard(lang)
			}

			bot.Send(msg)
		}
	}
}
