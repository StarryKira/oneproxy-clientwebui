package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"oneproxy-clientwebui/internal/config"
)

type session struct {
	createdAt time.Time
}

var (
	sessions   = make(map[string]session)
	sessionsMu sync.RWMutex
	sessionTTL = 2 * time.Hour
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func AdminAuth(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": false, "message": "未登录"})
		c.Abort()
		return
	}

	sessionsMu.RLock()
	s, ok := sessions[token]
	sessionsMu.RUnlock()

	if !ok || time.Since(s.createdAt) > sessionTTL {
		if ok {
			sessionsMu.Lock()
			delete(sessions, token)
			sessionsMu.Unlock()
		}
		c.JSON(http.StatusUnauthorized, gin.H{"code": false, "message": "会话已过期，请重新登录"})
		c.Abort()
		return
	}
	c.Next()
}

func GetConfig(c *gin.Context) {
	cfg := config.Get()
	c.JSON(http.StatusOK, gin.H{
		"api_base_url":  cfg.APIBaseURL,
		"exchange_rate": cfg.ExchangeRate,
	})
}

func UpdateConfig(c *gin.Context) {
	var body struct {
		APIBaseURL    string  `json:"api_base_url"`
		ExchangeRate  float64 `json:"exchange_rate"`
		AdminPassword string  `json:"admin_password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": false, "message": "invalid request body"})
		return
	}

	cfg := config.Get()
	cfg.APIBaseURL = body.APIBaseURL
	cfg.ExchangeRate = body.ExchangeRate

	if err := config.Save(cfg, body.AdminPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": true, "message": "config updated"})
}

func AdminLogin(c *gin.Context) {
	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": false, "message": "invalid request"})
		return
	}
	if !config.CheckPassword(body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": false, "message": "密码错误"})
		return
	}

	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": false, "message": "生成会话失败"})
		return
	}

	sessionsMu.Lock()
	sessions[token] = session{createdAt: time.Now()}
	sessionsMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"code": true, "message": "ok", "token": token})
}
