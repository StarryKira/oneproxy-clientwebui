package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"oneproxy-clientwebui/internal/model"
)

const configPath = "config.json"

var (
	mu  sync.RWMutex
	cfg model.AppConfig
)

func init() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	cfg = model.AppConfig{
		APIBaseURL:    "https://api.example.com",
		ExchangeRate:  500000,
		AdminPassword: string(hash),
	}
}

func Load() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return writeFile(cfg)
		}
		return err
	}
	return json.Unmarshal(data, &cfg)
}

func Get() model.AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

func CheckPassword(password string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return bcrypt.CompareHashAndPassword([]byte(cfg.AdminPassword), []byte(password)) == nil
}

func Save(newCfg model.AppConfig, rawPassword string) error {
	if newCfg.ExchangeRate <= 0 {
		return errors.New("exchange_rate must be greater than 0")
	}
	if newCfg.APIBaseURL == "" {
		return errors.New("api_base_url must not be empty")
	}

	if rawPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		newCfg.AdminPassword = string(hash)
	} else {
		mu.RLock()
		newCfg.AdminPassword = cfg.AdminPassword
		mu.RUnlock()
	}

	mu.Lock()
	defer mu.Unlock()
	cfg = newCfg
	return writeFile(cfg)
}

func writeFile(c model.AppConfig) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
