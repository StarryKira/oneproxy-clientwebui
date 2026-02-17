package model

type AppConfig struct {
	APIBaseURL    string  `json:"api_base_url"`
	ExchangeRate  float64 `json:"exchange_rate"`
	AdminPassword string  `json:"admin_password"`
}

type ExternalUsageResponse struct {
	Code    bool      `json:"code"`
	Data    UsageData `json:"data"`
	Message string    `json:"message"`
}

type UsageData struct {
	ExpiresAt          int64                  `json:"expires_at"`
	ModelLimits        map[string]interface{} `json:"model_limits"`
	ModelLimitsEnabled bool                   `json:"model_limits_enabled"`
	Name               string                 `json:"name"`
	Object             string                 `json:"object"`
	TotalAvailable     float64                `json:"total_available"`
	TotalGranted       float64                `json:"total_granted"`
	TotalUsed          float64                `json:"total_used"`
	UnlimitedQuota     bool                   `json:"unlimited_quota"`
}

type EnrichedUsageResponse struct {
	Code    bool              `json:"code"`
	Data    EnrichedUsageData `json:"data"`
	Message string            `json:"message"`
}

type EnrichedUsageData struct {
	UsageData
	USDAvailable float64 `json:"usd_available"`
	USDUsed      float64 `json:"usd_used"`
	USDTotal     float64 `json:"usd_total"`
}
