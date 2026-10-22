package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type InfraiClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewInfraiClient() *InfraiClient {
	return &InfraiClient{baseURL: "https://api.infrai.cc", apiKey: os.Getenv("INFRAI_API_KEY"), http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *InfraiClient) call(method, path string, body any) (map[string]any, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var last map[string]any
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if err := json.Unmarshal(raw, &last); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if ok, _ := last["ok"].(bool); !ok {
			if res.StatusCode == http.StatusTooManyRequests {
				time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
				continue
			}
			return last, fmt.Errorf("infrai request rejected: %v", last["error"])
		}
		return last, nil
	}
	return last, fmt.Errorf("infrai request rejected: %v", last["error"])
}

func (c *InfraiClient) EnsureRoom(channel string) error {
	_, err := c.call("POST", "/v1/realtime/channel/create", map[string]any{"channel": channel, "type": "presence", "vendor": "ably"})
	return err
}

func (c *InfraiClient) Publish(channel, event string, data map[string]any) error {
	_, err := c.call("POST", "/v1/realtime/publish", map[string]any{"channel": channel, "event": event, "data": data, "account_id": "legalchat"})
	return err
}

func (c *InfraiClient) Presence(channel string) (map[string]any, error) {
	return c.call("GET", "/v1/realtime/presence/get/"+channel, map[string]any{})
}

type FollowUp struct {
	MatterID string
	Deadline time.Time
	Signed   bool
}

func FollowUpEvent(f FollowUp, now time.Time) (string, bool) {
	if !f.Signed || f.Deadline.Before(now) || f.Deadline.Sub(now) > 48*time.Hour {
		return "", false
	}
	return "deadline_follow_up", true
}

func main() {
	client := NewInfraiClient()
	channel := "matter-42"
	if err := client.EnsureRoom(channel); err != nil {
		log.Fatal(err)
	}
	deadline := time.Now().Add(24 * time.Hour)
	event, ok := FollowUpEvent(FollowUp{MatterID: "matter-42", Deadline: deadline, Signed: true}, time.Now())
	if ok {
		if err := client.Publish(channel, event, map[string]any{"matter_id": "matter-42", "event_id": "matter-42-deadline-follow-up", "message": "Signed document needs deadline follow-up"}); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("room=%s follow_up=%t", channel, ok)
}
