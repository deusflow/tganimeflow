package bot

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

type TopAnimeResult struct {
	Text       string    // текст списка
	FirstAnime AnimeData // первое аниме для картинки
	HasData    bool      // есть ли данные
}
