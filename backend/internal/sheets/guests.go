package sheets

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Guest struct {
	GuestID   string `json:"guest_id"`
	Slug      string `json:"slug"`
	GuestName string `json:"guest_name"`
	CreatedAt string `json:"-"`
}

var (
	nonAlnum     = regexp.MustCompile(`[^a-z0-9]+`)
	guestCache   []Guest
	guestCacheAt time.Time
	guestCacheMu sync.Mutex
)

const guestCacheTTL = 60 * time.Second

func Slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlnum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func NewGuestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "g_" + hex.EncodeToString(b)
}

// CachedGuests returns a warm in-memory list, refreshing every 60s.
func (c *Client) CachedGuests(ctx context.Context, tab string) ([]Guest, error) {
	guestCacheMu.Lock()
	if guestCache != nil && time.Since(guestCacheAt) < guestCacheTTL {
		cached := guestCache
		guestCacheMu.Unlock()
		return cached, nil
	}
	guestCacheMu.Unlock()

	guests, err := c.SyncAndLoad(ctx, tab)
	if err != nil {
		return nil, err
	}

	guestCacheMu.Lock()
	guestCache = guests
	guestCacheAt = time.Now()
	guestCacheMu.Unlock()
	return guests, nil
}

// SyncAndLoad reads the guest sheet, generates missing technical metadata,
// writes it back idempotently, and returns the final guest list.
func (c *Client) SyncAndLoad(ctx context.Context, tab string) ([]Guest, error) {
	rng := fmt.Sprintf("%s!A2:D", tab)
	rows, err := c.Read(ctx, rng)
	if err != nil {
		return nil, err
	}

	// Pre-collect existing slugs so freshly-generated ones don't collide.
	existingSlugs := map[string]bool{}
	for _, row := range rows {
		if s := cell(row, 2); s != "" {
			existingSlugs[s] = true
		}
	}

	usedIDs := map[string]bool{}
	usedSlugs := map[string]bool{}
	guests := make([]Guest, 0, len(rows))
	outRows := make([][]interface{}, len(rows))
	needsWrite := false
	now := time.Now().UTC().Format(time.RFC3339)

	for i, row := range rows {
		name := cell(row, 0)
		guestID := cell(row, 1)
		slug := cell(row, 2)
		createdAt := cell(row, 3)

		if name == "" {
			outRows[i] = []interface{}{"", "", "", ""}
			continue
		}

		// --- guest_id ---
		if guestID == "" || usedIDs[guestID] {
			guestID = NewGuestID()
			for usedIDs[guestID] {
				guestID = NewGuestID()
			}
			needsWrite = true
		}
		usedIDs[guestID] = true

		// --- slug ---
		if slug == "" || usedSlugs[slug] {
			base := Slugify(name)
			if base == "" {
				base = "guest"
			}
			candidate := base
			n := 2
			for usedSlugs[candidate] || existingSlugs[candidate] {
				candidate = fmt.Sprintf("%s-%d", base, n)
				n++
			}
			slug = candidate
			needsWrite = true
		}
		usedSlugs[slug] = true

		// --- created_at ---
		if createdAt == "" {
			createdAt = now
			needsWrite = true
		}

		outRows[i] = []interface{}{name, guestID, slug, createdAt}
		guests = append(guests, Guest{
			GuestID:   guestID,
			Slug:      slug,
			GuestName: name,
			CreatedAt: createdAt,
		})
	}

	if needsWrite && len(outRows) > 0 {
		writeRng := fmt.Sprintf("%s!A2:D%d", tab, 1+len(outRows))
		if err := c.Write(ctx, writeRng, outRows); err != nil {
			return nil, err
		}
	}

	return guests, nil
}

func FindBySlug(guests []Guest, slug string) *Guest {
	for i := range guests {
		if guests[i].Slug == slug {
			return &guests[i]
		}
	}
	return nil
}

func cell(row []interface{}, i int) string {
	if i >= len(row) {
		return ""
	}
	s, _ := row[i].(string)
	return strings.TrimSpace(s)
}
