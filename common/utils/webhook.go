package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebhookPayload struct {
	Event       string `json:"event"`        // "milestone_reached" or "max_clicks_reached"
	ShortCode   string `json:"short_code"`
	Title       string `json:"title"`
	OriginalURL string `json:"original_url"`
	TotalClicks int64  `json:"total_clicks"`
	Timestamp   string `json:"timestamp"`
}

// DiscordEmbed structure for native Discord webhook integration
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

		// Automatically format Discord rich embeds if target is Discord
		if strings.Contains(webhookURL, "discord.com/api/webhooks") {
			embedColor := 0x38bdf8 // Sky blue
			icon := "📈"
			if payload.Event == "max_clicks_reached" {
				embedColor = 0xf87171 // Red
				icon = "🚫"
			}

			discordMsg := DiscordWebhookPayload{
				Username: "Shortlink Notifier",
				Embeds: []DiscordEmbed{
					{
						Title:       fmt.Sprintf("%s Event: %s (/%s)", icon, strings.ReplaceAll(payload.Event, "_", " "), payload.ShortCode),
						Description: fmt.Sprintf("Shortlink **%s** has hit an event notification.", payload.Title),
						Color:       embedColor,
						Timestamp:   payload.Timestamp,
						Fields: []DiscordField{
							{Name: "Short Code", Value: fmt.Sprintf("`/%s`", payload.ShortCode), Inline: true},
							{Name: "Total Clicks", Value: fmt.Sprintf("**%d**", payload.TotalClicks), Inline: true},
							{Name: "Destination", Value: payload.OriginalURL, Inline: false},
						},
					},
				},
			}
			bodyBytes, err = json.Marshal(discordMsg)
		} else {
			bodyBytes, err = json.Marshal(payload)
		}

		if err != nil {
			log.Printf("[Webhook] Error serializing payload: %v", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Printf("[Webhook] Error creating request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "UltraFast-Redirector-Bot/1.0")

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
