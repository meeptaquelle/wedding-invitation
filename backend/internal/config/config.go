package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceAccountEmail string
	PrivateKey          string
	SheetID             string
	GuestListTab        string
	MessagesTab         string
	PageSize            int
	MaxMessageLength    int
	RateLimitSeconds    int
}

func Load() *Config {
	return &Config{
		ServiceAccountEmail: os.Getenv("GOOGLE_SERVICE_ACCOUNT_EMAIL"),
		PrivateKey:          normalizeKey(os.Getenv("GOOGLE_PRIVATE_KEY")),
		SheetID:             os.Getenv("GOOGLE_SHEET_ID"),
		GuestListTab:        getEnv("GUEST_LIST_TAB", "Guests"),
		MessagesTab:         getEnv("MESSAGES_TAB", "Messages"),
		PageSize:            getEnvInt("PAGE_SIZE", 20),
		MaxMessageLength:    getEnvInt("MAX_MESSAGE_LENGTH", 500),
		RateLimitSeconds:    getEnvInt("RATE_LIMIT_SECONDS", 30),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getEnvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func normalizeKey(k string) string {
	// Vercel often stores private keys with literal \n escapes
	return strings.ReplaceAll(k, `\n`, "\n")
}
