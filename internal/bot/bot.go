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
	"time"
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_top_popular"], "action_top_popular"),
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_top_season"], "action_top_season"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages[lang]["btn_top_year"], "action_top_year"),
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

// Универсальная функция для получения топ аниме

func getTopAnimeWithFirst(url, messageKey, lang string) TopAnimeResult {
	var result JikanResponse
	if err := fetchAndUnmarshal(url, &result); err != nil {
		logRequest("getTopAnimeList", err)
		return TopAnimeResult{
			Text:    messages[lang]["api_error"],
			HasData: false,
		}
	}

	if len(result.Data) == 0 {
		return TopAnimeResult{
			Text:    messages[lang]["not_found"],
			HasData: false,
		}
	}

	// Формируем текст списка
	topAnime := messages[lang][messageKey] + "\n\n"
	for i, anime := range result.Data {
		topAnime += fmt.Sprintf("%d. %s - ⭐ %.1f\n", i+1, anime.Title, anime.Score)
	}

	return TopAnimeResult{
		Text:       topAnime,
		FirstAnime: result.Data[0], // первое аниме для картинки
		HasData:    true,
	}
}

func getTopAnime(lang string) TopAnimeResult {
	url := fmt.Sprintf("%s/top/anime?limit=5", jikanBaseURL)
	return getTopAnimeWithFirst(url, "top_anime", lang)
}

func getTopPopularAnime(lang string) TopAnimeResult {
	url := fmt.Sprintf("%s/top/anime?filter=bypopularity&limit=5", jikanBaseURL)
	return getTopAnimeWithFirst(url, "top_popular", lang)
}

func getTopSeasonAnime(lang string) TopAnimeResult {
	// Определяем текущий год и месяц
	year := time.Now().Year()
	month := time.Now().Month()

	// Определяем сезон по месяцу
	season := "winter"
	switch month {
	case 1, 2, 12:
		season = "winter"
	case 3, 4, 5:
		season = "spring"
	case 6, 7, 8:
		season = "summer"
	case 9, 10, 11:
		season = "fall"
	}

	// Формируем URL и используем универсальную функцию
	url := fmt.Sprintf("%s/seasons/%d/%s?limit=5", jikanBaseURL, year, season)
	return getTopAnimeWithFirst(url, "top_season", lang)
}

func getTopYearAnime(lang string) TopAnimeResult {
	year := time.Now().Year()
	url := fmt.Sprintf("%s/anime?start_date=%d-01-01&end_date=%d-12-31&order_by=score&sort=desc&limit=5", jikanBaseURL, year, year)
	return getTopAnimeWithFirst(url, "top_year", lang)
}

func formatAnimeDetails(anime AnimeData, lang string) string {
	// Форматирую жанры в ст��оку
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

	//описание ограничи��аем длину
	synopsis := anime.Synopsis
	if len(synopsis) > 200 {
		synopsis = synopsis[:200] + "..."
	}

	// Определяем текст для серий в зависимости от языка
	episodesLabel := "episodes"
	switch lang {
	case "ua":
		episodesLabel = "серій"
	case "en":
		episodesLabel = "episodes"
	case "da":
		episodesLabel = "episoder"
	default:
		episodesLabel = "episodes"
	}

	return fmt.Sprintf(
		"🎌 %s\n⭐ %.1f\n📺 %s %s\n📊 %s\n🎭 %s\n\n📝 %s",
		anime.Title,
		anime.Score,
		episodesText,
		episodesLabel,
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

	if !knownUsers[userID] {
		knownUsers[userID] = true
		botAnalytics.TotalUsers++
		fmt.Printf("📊 Новый пользователь! Всего: %d\n", botAnalytics.TotalUsers)
	}

	// Считаем использование команд
	botAnalytics.CommandsUsed[action]++

	// Счи��аем использование языков
	botAnalytics.LanguagesUsed[lang]++

	fmt.Printf("📊 Действие: %s, Язык: %s, Всего команд '%s': %d\n",
		action, lang, action, botAnalytics.CommandsUsed[action])
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
	userLangs := make(map[int64]string) // userID -> выбр���нный язык

	for update := range updates {
		if update.Message != nil {
			// Проверяем, если это команда для смены языка
			fmt.Println("Message Received:", update.Message.Text)

			// Получаем общие данн����е ОДИН РАЗ для всех команд
			userID := update.Message.From.ID
			chatID := update.Message.Chat.ID
			lang := userLangs[userID]
			if lang == "" {
				lang = "ua" // По умолчанию українс���ий
			}

			var responseText string
			var keyboard *tgbotapi.InlineKeyboardMarkup

			if update.Message.IsCommand() && update.Message.Command() == cmdStart {
				logUserAction(userID, "start", lang)
				responseText = messages[lang]["start"]

				// Если язы�� уже выбран, показываем кнопки д��йствий, иначе выбор языка
				if userLangs[userID] != "" {
					quickKeyboard := createQuickActionsKeyboard(lang)
					keyboard = &quickKeyboard
				} else {
					languageKeyboard := createLanguageKeyboard()
					keyboard = &languageKeyboard
				}

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
				topResult := getTopAnime(lang)
				if topResult.HasData {
					// Сначала отправляем текст списка
					msg := tgbotapi.NewMessage(chatID, topResult.Text)
					bot.Send(msg)
					// Потом отправляем первое аниме с картинкой
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, topResult.FirstAnime, lang, &quickKeyboard)
					continue
				} else {
					responseText = topResult.Text // сообщение об ошибке
					quickKeyboard := createQuickActionsKeyboard(lang)
					keyboard = &quickKeyboard
				}

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
			var keyboard *tgbotapi.InlineKeyboardMarkup

			switch update.CallbackQuery.Data {
			case "lang_ua":
				logUserAction(userID, "lang_change", "ua")
				userLangs[userID] = "ua"
				responseText = messages["ua"]["lang_changed"] + "\n" + messages["ua"]["start"]
				quickKeyboard := createQuickActionsKeyboard("ua")
				keyboard = &quickKeyboard

			case "lang_en":
				logUserAction(userID, "lang_change", "en")
				userLangs[userID] = "en"
				responseText = messages["en"]["lang_changed"] + "\n" + messages["en"]["start"]
				quickKeyboard := createQuickActionsKeyboard("en")
				keyboard = &quickKeyboard

			case "lang_da":
				logUserAction(userID, "lang_change", "da")
				userLangs[userID] = "da"
				responseText = messages["da"]["lang_changed"] + "\n" + messages["da"]["start"]
				quickKeyboard := createQuickActionsKeyboard("da")
				keyboard = &quickKeyboard

			case "action_random":
				logUserAction(userID, "random", lang)
				anime := getRandomAnime(lang)
				quickKeyboard := createQuickActionsKeyboard(lang)
				sendAnimeWithPhoto(bot, chatID, anime, lang, &quickKeyboard)
				continue

			case "action_top":
				logUserAction(userID, "top", lang)
				topResult := getTopAnime(lang)
				if topResult.HasData {
					// Сначала отправляем текст списка
					msg := tgbotapi.NewMessage(chatID, topResult.Text)
					bot.Send(msg)
					// Потом отправляем первое аниме с картинкой
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, topResult.FirstAnime, lang, &quickKeyboard)
					continue
				} else {
					responseText = topResult.Text // сообщение об ошибке
				}
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			case "action_top_popular":
				logUserAction(userID, "top_popular", lang)
				topResult := getTopPopularAnime(lang)
				if topResult.HasData {
					// Сначала отправляем текст списка
					msg := tgbotapi.NewMessage(chatID, topResult.Text)
					bot.Send(msg)
					// Потом отправляем первое аниме с картинкой
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, topResult.FirstAnime, lang, &quickKeyboard)
					continue
				} else {
					responseText = topResult.Text // сообщение об ошибке
				}
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			case "action_top_year":
				logUserAction(userID, "top_year", lang)
				topResult := getTopYearAnime(lang)
				if topResult.HasData {
					// Сначала отправляем текст списка
					msg := tgbotapi.NewMessage(chatID, topResult.Text)
					bot.Send(msg)
					// Потом отправляем первое аниме с картинкой
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, topResult.FirstAnime, lang, &quickKeyboard)
					continue
				} else {
					responseText = topResult.Text // сообщение об ошибке
				}
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			case "action_top_season":
				logUserAction(userID, "top_season", lang)
				topResult := getTopSeasonAnime(lang)
				if topResult.HasData {
					// Сначала отправляем текст списка
					msg := tgbotapi.NewMessage(chatID, topResult.Text)
					bot.Send(msg)
					// Потом отправляем первое аниме с картинкой
					quickKeyboard := createQuickActionsKeyboard(lang)
					sendAnimeWithPhoto(bot, chatID, topResult.FirstAnime, lang, &quickKeyboard)
					continue
				} else {
					responseText = topResult.Text // сообщение об ошибке
				}
				quickKeyboard := createQuickActionsKeyboard(lang)
				keyboard = &quickKeyboard

			case "donate_thanks":
				logUserAction(userID, "donate_thanks", lang)
				responseText = messages[lang]["donate_thanks"]
			}

			msg := tgbotapi.NewMessage(chatID, responseText)

			if keyboard != nil {
				msg.ReplyMarkup = *keyboard
			}

			bot.Send(msg)
		}
	}
}
