import axios from 'axios'

const http = axios.create({ baseURL: '/' })

export function setAdminPassword(password: string) {
  http.defaults.headers.common['X-Admin-Password'] = password
}

export interface UsageData {
  name: string
  object: string
  total_available: number
  total_granted: number
  total_used: number
  unlimited_quota: boolean
  expires_at: number
  model_limits_enabled: boolean
  model_limits: Record<string, unknown>
  usd_available: number
  usd_used: number
  usd_total: number
}

export interface UsageResult {
  code: boolean
  message: string
  data: UsageData
}

export interface AppConfig {
  api_base_url: string
  exchange_rate: number
  admin_password: string
}

export const queryUsage = (key: string) =>
  http.get<UsageResult>('/api/usage/query', { params: { key } })

export const adminLogin = (password: string) =>
  http.post<{ code: boolean; message: string }>('/api/admin/login', { password })

export const getConfig = () =>
  http.get<AppConfig>('/api/config')

export const updateConfig = (cfg: AppConfig) =>
  http.post<{ code: boolean; message: string }>('/api/config', cfg)
