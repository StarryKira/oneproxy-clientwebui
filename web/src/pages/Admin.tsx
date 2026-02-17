import { useState } from 'react'
import { Button, Form, Input, InputNumber, message, Spin, Tooltip, Typography } from 'antd'
import { LockOutlined, QuestionCircleOutlined } from '@ant-design/icons'
import { adminLogin, getConfig, setAdminPassword, updateConfig, type AppConfig } from '../api'

const { Title } = Typography

export default function Admin() {
  const [authed, setAuthed] = useState(false)
  const [loginLoading, setLoginLoading] = useState(false)
  const [form] = Form.useForm<AppConfig>()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  const onLogin = async (values: { password: string }) => {
    setLoginLoading(true)
    try {
      const res = await adminLogin(values.password)
      if (res.data.code) {
        setAdminPassword(values.password)
        setAuthed(true)
        setLoading(true)
        getConfig()
          .then(res => form.setFieldsValue(res.data))
          .catch(() => message.error('加载配置失败'))
          .finally(() => setLoading(false))
      } else {
        message.error(res.data.message)
      }
    } catch {
      message.error('登录失败')
    } finally {
      setLoginLoading(false)
    }
  }

  const onFinish = async (values: AppConfig) => {
    setSaving(true)
    try {
      const res = await updateConfig(values)
      if (res.data.code) {
        message.success('配置已保存')
      } else {
        message.error(res.data.message)
      }
    } catch {
      message.error('保存失败')
    } finally {
      setSaving(false)
    }
  }

  if (!authed) {
    return (
      <div style={{ maxWidth: 400, margin: '80px auto' }}>
        <Title level={3}>管理员登录</Title>
        <Form onFinish={onLogin}>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="管理员密码" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loginLoading} block size="large">
              登录
            </Button>
          </Form.Item>
        </Form>
      </div>
    )
  }

  return (
    <div style={{ maxWidth: 500, margin: '0 auto' }}>
      <Title level={3}>系统配置</Title>
      <Spin spinning={loading}>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="api_base_url"
            label="API 基础地址"
            rules={[{ required: true, message: '请输入 API 地址' }]}
          >
            <Input placeholder="https://api.example.com" />
          </Form.Item>
          <Form.Item
            name="exchange_rate"
            label={
              <span>
                换算比例&nbsp;
                <Tooltip title="多少额度单位等于 1 美金">
                  <QuestionCircleOutlined />
                </Tooltip>
              </span>
            }
            rules={[{ required: true, message: '请输入换算比例' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} placeholder="500000" />
          </Form.Item>
          <Form.Item
            name="admin_password"
            label="修改管理员密码"
          >
            <Input.Password placeholder="留空则不修改" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving} block>
              保存配置
            </Button>
          </Form.Item>
        </Form>
      </Spin>
    </div>
  )
}
