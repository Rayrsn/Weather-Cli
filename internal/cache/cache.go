package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Rayrsn/weather-Cli/internal/api"
)

type CacheEntry struct {
	Data      api.CombinedResponse `json:"data"`
	Timestamp time.Time           `json:"timestamp"`
}

type Cache struct {
	Dir string
}

func NewCache() (*Cache, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	cacheDir := filepath.Join(userCacheDir, "weather-cli")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	return &Cache{Dir: cacheDir}, nil
}

func (c *Cache) Get(key string, ttl time.Duration) (*api.CombinedResponse, bool) {
	path := filepath.Join(c.Dir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if time.Since(entry.Timestamp) > ttl {
		return nil, false
	}

	return &entry.Data, true
}

func (c *Cache) Set(key string, data api.CombinedResponse) error {
	path := filepath.Join(c.Dir, key+".json")
	entry := CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

func GenerateKey(cityName string, isImperial bool, showForecast bool) string {
	// Simple key generation. For more robustness, could use MD5/SHA.
	units := "metric"
	if isImperial {
		units = "imperial"
	}
	forecast := "no-forecast"
	if showForecast {
		forecast = "forecast"
	}
	return fmt.Sprintf("%s-%s-%s", cityName, units, forecast)
}
