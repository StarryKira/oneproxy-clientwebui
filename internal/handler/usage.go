package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"oneproxy-clientwebui/internal/config"
	"oneproxy-clientwebui/internal/model"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func QueryUsage(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" || !strings.HasPrefix(key, "sk-") {
		c.JSON(http.StatusBadRequest, gin.H{"code": false, "message": "invalid API key"})
		return
	}

	cfg := config.Get()
	url := fmt.Sprintf("%s/api/usage/token", strings.TrimRight(cfg.APIBaseURL, "/"))

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": false, "message": err.Error()})
		return
	}
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": false, "message": "upstream API unreachable: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": false, "message": "failed to read upstream response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"code": false, "message": fmt.Sprintf("upstream returned status %d", resp.StatusCode)})
		return
	}

	var extResp model.ExternalUsageResponse
	if err := json.Unmarshal(body, &extResp); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": false, "message": "failed to parse upstream response"})
		return
	}

	if !extResp.Code {
		c.JSON(http.StatusOK, gin.H{"code": false, "message": extResp.Message})
		return
	}

	rate := cfg.ExchangeRate
	enriched := model.EnrichedUsageResponse{
		Code:    true,
		Message: extResp.Message,
		Data: model.EnrichedUsageData{
			UsageData:    extResp.Data,
			USDAvailable: extResp.Data.TotalAvailable / rate,
			USDUsed:      extResp.Data.TotalUsed / rate,
			USDTotal:     extResp.Data.TotalGranted / rate,
		},
	}

	c.JSON(http.StatusOK, enriched)
}
