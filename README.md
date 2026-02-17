# OneProxy Client WebUI

API 额度查询 + AI 接口反向代理，支持 OpenAI / Anthropic 等兼容接口。

## 功能

- 输入 API Key 查询余额，自动换算为美金
- 反向代理 AI 接口（`/v1/*`），隐藏真实后端地址，支持 SSE 流式和 MCP
- 管理员配置页面（密码 bcrypt 加密存储）

## 快速开始

### Docker Compose

```bash
docker compose up -d
```

访问 `http://localhost:8080`

### 本地开发

```bash
# 前端
cd web && npm install && npm run build && cd ..

# 启动
go run .
```

## 管理员配置

浏览器访问 `/admin`（导航栏不可见，需手动输入），默认密码 `admin123`。

可配置项：
- API 基础地址（上游 AI 服务地址）
- 换算比例（额度单位 → 美金）
- 管理员密码

## 反向代理

将客户端 base URL 设为 `http://localhost:8080`，所有 `/v1/*` 请求会被透传到配置的上游地址。

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-xxx" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"hello"}]}'
```

## 项目结构

```
├── main.go                  # Gin 入口
├── internal/
│   ├── model/model.go       # 数据模型
│   ├── config/config.go     # 配置管理（bcrypt + RWMutex）
│   └── handler/
│       ├── usage.go         # 额度查询
│       ├── config.go        # 配置 API + 鉴权
│       └── proxy.go         # 反向代理
├── web/                     # React + Vite + Ant Design
├── Dockerfile
└── docker-compose.yml
```
