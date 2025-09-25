package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// PersistenceManager handles color persistence using Redis
type PersistenceManager struct {
	client *redis.Client
	ctx    context.Context
}

// ColorEntry represents a stored color with metadata
type ColorEntry struct {
	Color     RGB       `json:"color"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"` // "directory", "claude", "manual"
}

// NewPersistenceManager creates a new persistence manager
func NewPersistenceManager() *PersistenceManager {
	// Try to connect to Redis with common configurations
	addresses := []string{
		"localhost:6379", // Default Redis port
		"127.0.0.1:6379", // Alternative localhost
		"redis:6379",     // Docker container name
	}

	var client *redis.Client
	ctx := context.Background()

	for _, addr := range addresses {
		client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     "", // No password by default
			DB:           0,  // Default DB
			DialTimeout:  time.Second * 2,
			ReadTimeout:  time.Second * 2,
			WriteTimeout: time.Second * 2,
		})

		// Test connection
		_, err := client.Ping(ctx).Result()
		if err == nil {
			break
		}
		client.Close()
		client = nil
	}

	return &PersistenceManager{
		client: client,
		ctx:    ctx,
	}
}

// IsEnabled returns true if Redis connection is available
func (pm *PersistenceManager) IsEnabled() bool {
	return pm.client != nil
}

// GetDirectoryColor retrieves stored color for a directory
func (pm *PersistenceManager) GetDirectoryColor(directoryPath string) (RGB, bool) {
	if !pm.IsEnabled() {
		return RGB{}, false
	}

	key := fmt.Sprintf("color:directory:%s", directoryPath)
	data, err := pm.client.Get(pm.ctx, key).Result()
	if err == redis.Nil {
		return RGB{}, false // Key doesn't exist
	}
	if err != nil {
		log.Printf("Redis error getting directory color: %v", err)
		return RGB{}, false
	}

	var entry ColorEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		log.Printf("Error unmarshaling color entry: %v", err)
		return RGB{}, false
	}

	return entry.Color, true
}

// SetDirectoryColor stores color for a directory
func (pm *PersistenceManager) SetDirectoryColor(directoryPath string, color RGB) error {
	if !pm.IsEnabled() {
		return nil // Fail silently if Redis unavailable
	}

	entry := ColorEntry{
		Color:     color,
		Timestamp: time.Now(),
		Source:    "directory",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error marshaling color entry: %w", err)
	}

	key := fmt.Sprintf("color:directory:%s", directoryPath)
	// Set with 30 day expiration to prevent infinite growth
	err = pm.client.Set(pm.ctx, key, data, time.Hour*24*30).Err()
	if err != nil {
		log.Printf("Redis error setting directory color: %v", err)
	}

	return err
}

// GetLastClaudeColor retrieves the last used Claude theme color
func (pm *PersistenceManager) GetLastClaudeColor() (RGB, bool) {
	if !pm.IsEnabled() {
		return RGB{}, false
	}

	key := "color:claude:last"
	data, err := pm.client.Get(pm.ctx, key).Result()
	if err == redis.Nil {
		return RGB{}, false
	}
	if err != nil {
		log.Printf("Redis error getting Claude color: %v", err)
		return RGB{}, false
	}

	var entry ColorEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		log.Printf("Error unmarshaling Claude color entry: %v", err)
		return RGB{}, false
	}

	// Only return if it's recent (within 24 hours)
	if time.Since(entry.Timestamp) > time.Hour*24 {
		return RGB{}, false
	}

	return entry.Color, true
}

// SetLastClaudeColor stores the last used Claude theme color
func (pm *PersistenceManager) SetLastClaudeColor(color RGB) error {
	if !pm.IsEnabled() {
		return nil
	}

	entry := ColorEntry{
		Color:     color,
		Timestamp: time.Now(),
		Source:    "claude",
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error marshaling Claude color entry: %w", err)
	}

	key := "color:claude:last"
	// Claude colors expire after 24 hours to allow theme variation
	err = pm.client.Set(pm.ctx, key, data, time.Hour*24).Err()
	if err != nil {
		log.Printf("Redis error setting Claude color: %v", err)
	}

	return err
}

// GetColorHistory retrieves recent color history
func (pm *PersistenceManager) GetColorHistory(limit int) ([]ColorEntry, error) {
	if !pm.IsEnabled() {
		return []ColorEntry{}, nil
	}

	// Get all directory color keys
	keys, err := pm.client.Keys(pm.ctx, "color:directory:*").Result()
	if err != nil {
		return []ColorEntry{}, err
	}

	var entries []ColorEntry
	for _, key := range keys {
		data, err := pm.client.Get(pm.ctx, key).Result()
		if err != nil {
			continue
		}

		var entry ColorEntry
		if err := json.Unmarshal([]byte(data), &entry); err != nil {
			continue
		}

		entries = append(entries, entry)
		if len(entries) >= limit {
			break
		}
	}

	return entries, nil
}

// ClearColorCache removes all stored colors
func (pm *PersistenceManager) ClearColorCache() error {
	if !pm.IsEnabled() {
		return nil
	}

	keys, err := pm.client.Keys(pm.ctx, "color:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return pm.client.Del(pm.ctx, keys...).Err()
	}

	return nil
}

// GetConnectionStatus returns Redis connection status
func (pm *PersistenceManager) GetConnectionStatus() string {
	if !pm.IsEnabled() {
		return "❌ Redis unavailable (colors won't persist)"
	}

	_, err := pm.client.Ping(pm.ctx).Result()
	if err != nil {
		return "⚠️ Redis connection issues"
	}

	return "✅ Redis connected (colors will persist)"
}

// Close closes the Redis connection
func (pm *PersistenceManager) Close() error {
	if pm.client != nil {
		return pm.client.Close()
	}
	return nil
}