package logger

import (
	"fmt"
	"net/http"
	"net/url"
)

type Logger interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

func SendToTelegram(botToken string, chatID string, message string) error {
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Prepare the message payload
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "MarkdownV2") // Optional, for MarkdownV2 formatting

	// Make the POST request to Telegram API
	resp, err := http.PostForm(telegramAPI, data)
	if err != nil {
		return fmt.Errorf("failed to send message to Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Telegram: %d", resp.StatusCode)
	}

	return nil
}
