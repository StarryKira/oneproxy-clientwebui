package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"oneproxy-clientwebui/internal/config"
)

func AdminAuth(c *gin.Context) {
	password := c.GetHeader("X-Admin-Password")
	if !config.CheckPassword(password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": false, "message": "密码错误"})
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
	c.JSON(http.StatusOK, gin.H{"code": true, "message": "ok"})
}
