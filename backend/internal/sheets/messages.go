package sheets

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Message struct {
	MessageID   string `json:"message_id"`
	GuestID     string `json:"guest_id,omitempty"`
	DisplayName string `json:"display_name"`
	Message     string `json:"message"`
	CreatedAt   string `json:"created_at"`
}

func NewMessageID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "m_" + hex.EncodeToString(b)
}

func HashIP(ip string) string {
	h := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(h[:8])
}

func (c *Client) LoadMessages(ctx context.Context, tab string) ([]Message, error) {
	rng := fmt.Sprintf("%s!A2:F", tab)
	rows, err := c.Read(ctx, rng)
	if err != nil {
		return nil, err
	}

	msgs := make([]Message, 0, len(rows))
	for _, row := range rows {
		m := Message{
			MessageID:   cell(row, 0),
			GuestID:     cell(row, 1),
			DisplayName: cell(row, 2),
			Message:     cell(row, 3),
			CreatedAt:   cell(row, 4),
		}
		if m.MessageID == "" || m.Message == "" {
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (c *Client) AppendMessage(ctx context.Context, tab string, m Message, ipHash string) error {
	row := []interface{}{
		m.MessageID,
		m.GuestID,
		m.DisplayName,
		m.Message,
		m.CreatedAt,
		ipHash,
	}
	return c.Append(ctx, fmt.Sprintf("%s!A:F", tab), [][]interface{}{row})
}

// RecentIPSubmission returns true if this IP has posted within the window.
func (c *Client) RecentIPSubmission(ctx context.Context, tab, ipHash string, window time.Duration) (bool, error) {
	rng := fmt.Sprintf("%s!A2:F", tab)
	rows, err := c.Read(ctx, rng)
	if err != nil {
		return false, err
	}
	cutoff := time.Now().Add(-window).UTC()
	for i := len(rows) - 1; i >= 0 && i >= len(rows)-30; i-- {
		row := rows[i]
		if cell(row, 5) != ipHash {
			continue
		}
		t, err := time.Parse(time.RFC3339, cell(row, 4))
		if err != nil {
			continue
		}
		if t.After(cutoff) {
			return true, nil
		}
	}
	return false, nil
}
