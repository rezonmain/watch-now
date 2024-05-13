package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
)

const AuthURL = "https://id.twitch.tv/oauth2/token"

type Token struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func getTwitchToken() *Token {
	// 1. get twitch authorization token
	clientId := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	authReqBody := []byte(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", clientId, clientSecret))
	bodyReader := bytes.NewReader(authReqBody)
	req, err := http.NewRequest("POST", AuthURL, bodyReader)
	if err != nil {
		fmt.Printf("Unable to create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	var token Token
	if err := json.Unmarshal(body, &token); err != nil { // Parse []byte to the go struct pointer
		panic(err)
	}
	return &token
}

func subscribeToWebhooks(token *Token, channelNames []string) {
	return
}

func main() {
	accessToken := getTwitchToken()
	telegram, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}
	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	go subscribeToWebhooks(accessToken, []string{"hasanabi"})

	router := http.NewServeMux()
	router.HandleFunc("/twitch/callback", func(w http.ResponseWriter, r *http.Request) {
		msg := tgbotapi.NewMessage(chatId, "hi from go")
		if _, err = telegram.Send(msg); err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	// 4. Server
	Port := fmt.Sprintf(":%s",os.Getenv("WATCH_NOW_PORT"))
	s := http.Server{
		Addr:    Port,
		Handler: router,
	}

	fmt.Printf("Listening in port %s\n", Port)

	s.ListenAndServe()
}
