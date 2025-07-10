# Telegram Anime Bot: DeusFlow

This is a simple Telegram bot for anime content, built with Go. I made it for fun and learning, and I hope you find it useful or inspiring!

---

## How to Get Started

### Run Locally

1.  **Clone the Repository:**
    ```bash
    git clone [https://github.com/deusflow/tganimeflow.git](https://github.com/deusflow/tganimeflow.git)
    cd tganimeflow
    ```

2.  **Set Up Environment Variables:**
    Create a `.env` file in the root of the project, or set these variables directly in your environment:

    * `TELEGRAM_TOKEN`: This is your unique token from BotFather on Telegram.
    * **If you're using a database (optional):**
        * `DB_HOST`
        * `DB_PORT`
        * `DB_NAME`
        * `DB_USER`
        * `DB_PASSWORD`

3.  **Run the Bot:**
    ```bash
    go run main.go
    ```

---

## Deploying to Railway

Railway makes deployment super easy!

1.  **Log In to Railway:**
    Head over to [railway.app](https://railway.app/) and sign in using your GitHub account.

2.  **Start a New Project:**
    Click "New Project" and then select "Deploy from GitHub repo."

3.  **Select Your Repository:**
    Find and choose the `deusflow/tganimeflow` repository from your list.

4.  **Configure Environment Variables:**
    In the Railway dashboard, add the environment variables listed in the "How to Get Started" section above (especially your `TELEGRAM_TOKEN` and any database credentials).

5.  **You're Done!**
    Railway will automatically build and launch your bot.

---

## About This Project

* **Core Logic:** You'll find the main bot code within the `internal/bot/` directory.
* **Dependencies:** This project leverages the popular [`go-telegram-bot-api`](https://github.com/go-telegram-bot-api/telegram-bot-api) library for Telegram bot interactions.
* **Purpose:** This bot is a personal project aimed at learning and having fun with Go and Telegram bot development. Your contributions and pull requests are very welcome!
