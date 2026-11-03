package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WebhookPayload struct {
	Event       string `json:"event"` // "milestone_reached" or "max_clicks_reached"
	ShortCode   string `json:"short_code"`
	Title       string `json:"title"`
	OriginalURL string `json:"original_url"`
	TotalClicks int64  `json:"total_clicks"`
	Timestamp   string `json:"timestamp"`
}

// Discord Webhook Payload
type DiscordWebhookPayload struct {
	Username string         `json:"username"`
	Embeds   []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Fields      []DiscordField `json:"fields"`
	Timestamp   string         `json:"timestamp"`
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Slack Webhook Payload
type SlackWebhookPayload struct {
	Text   string       `json:"text"`
	Blocks []SlackBlock `json:"blocks,omitempty"`
}

type SlackBlock struct {
	Type string     `json:"type"`
	Text *SlackText `json:"text,omitempty"`
}

type SlackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Telegram Webhook Payload
type TelegramSendMessagePayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// DispatchWebhook sends an asynchronous HTTP notification to a user's webhook URL
func DispatchWebhook(webhookURL string, payload WebhookPayload) {
	if strings.TrimSpace(webhookURL) == "" {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var bodyBytes []byte
		var err error
		targetURL := webhookURL

		// Format title
		displayTitle := payload.Title
		if displayTitle == "" {
			displayTitle = fmt.Sprintf("Shortlink /%s", payload.ShortCode)
		}

		icon := "📈"
		if payload.Event == "max_clicks_reached" {
			icon = "🚫"
		}

		// 1. DISCORD WEBHOOK
		if strings.Contains(webhookURL, "discord.com/api/webhooks") {
			embedColor := 0x38bdf8 // Sky blue
			if payload.Event == "max_clicks_reached" {
				embedColor = 0xf87171 // Red
			}

			discordMsg := DiscordWebhookPayload{
				Username: "Shortlink Notifier",
				Embeds: []DiscordEmbed{
					{
						Title:       fmt.Sprintf("%s Event: %s (/%s)", icon, strings.ReplaceAll(payload.Event, "_", " "), payload.ShortCode),
						Description: fmt.Sprintf("Shortlink **%s** has triggered a telemetry milestone notification.", displayTitle),
						Color:       embedColor,
						Timestamp:   payload.Timestamp,
						Fields: []DiscordField{
							{Name: "Short Code", Value: fmt.Sprintf("`/%s`", payload.ShortCode), Inline: true},
							{Name: "Total Clicks", Value: fmt.Sprintf("**%d**", payload.TotalClicks), Inline: true},
							{Name: "Target URL", Value: payload.OriginalURL, Inline: false},
						},
					},
				},
			}
			bodyBytes, err = json.Marshal(discordMsg)

		// 2. SLACK WEBHOOK
		} else if strings.Contains(webhookURL, "hooks.slack.com/services") {
			slackMsg := SlackWebhookPayload{
				Text: fmt.Sprintf("%s *Event: %s* on `/%s` - Clicks: %d - Target: %s",
					icon, strings.ReplaceAll(payload.Event, "_", " "), payload.ShortCode, payload.TotalClicks, payload.OriginalURL),
				Blocks: []SlackBlock{
					{
						Type: "section",
						Text: &SlackText{
							Type: "mrkdwn",
							Text: fmt.Sprintf("%s *Shortlink Milestone Alert*\n*Code:* `/%s`\n*Clicks:* *%d*\n*Target:* %s",
								icon, payload.ShortCode, payload.TotalClicks, payload.OriginalURL),
						},
					},
				},
			}
			bodyBytes, err = json.Marshal(slackMsg)

		// 3. TELEGRAM BOT
		} else if strings.Contains(webhookURL, "api.telegram.org/bot") {
			parsed, parseErr := url.Parse(webhookURL)
			chatID := ""
			if parseErr == nil {
				chatID = parsed.Query().Get("chat_id")
			}
			if chatID == "" {
				chatID = "@channel"
			}

			tgMsg := TelegramSendMessagePayload{
				ChatID: chatID,
				Text: fmt.Sprintf("%s *Shortlink Alert*\n*Event:* %s\n*Code:* `/%s`\n*Total Clicks:* *%d*\n*Target:* %s",
					icon, strings.ReplaceAll(payload.Event, "_", " "), payload.ShortCode, payload.TotalClicks, payload.OriginalURL),
				ParseMode: "Markdown",
			}
			bodyBytes, err = json.Marshal(tgMsg)

		// 4. GENERIC / CUSTOM REST WEBHOOK
		} else {
			bodyBytes, err = json.Marshal(payload)
		}

		if err != nil {
			log.Printf("[Webhook] Error serializing payload: %v", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Printf("[Webhook] Error creating request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "UltraFast-Redirector-Bot/1.0")
		req.Header.Set("X-Event-Type", payload.Event)
		req.Header.Set("X-Short-Code", payload.ShortCode)
		req.Header.Set("X-Timestamp", payload.Timestamp)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Webhook] Warning: failed to send webhook to %s: %v", webhookURL, err)
			return
		}
		defer resp.Body.Close()

		log.Printf("[Webhook] Notification sent for /%s (status %d)", payload.ShortCode, resp.StatusCode)
	}()
}
