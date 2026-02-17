import { useState } from 'react'
import { Input, message, Spin, Typography } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { queryUsage, type UsageData } from '../api'
import UsageCard from '../components/UsageCard'

const { Title } = Typography

export default function Home() {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<UsageData | null>(null)

  const onSearch = async (key: string) => {
    const trimmed = key.trim()
    if (!trimmed || !trimmed.startsWith('sk-')) {
      message.warning('请输入有效的 API Key（以 sk- 开头）')
      return
    }
    setLoading(true)
    setData(null)
    try {
      const res = await queryUsage(trimmed)
      if (res.data.code) {
        setData(res.data.data)
      } else {
        message.error(res.data.message || '查询失败')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '网络错误'
      message.error(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 700, margin: '0 auto' }}>
      <Title level={3}>API 额度查询</Title>
      <Input.Search
        placeholder="输入 API Key（sk-xxxxxx）"
        enterButton={<><SearchOutlined /> 查询</>}
        size="large"
        onSearch={onSearch}
        loading={loading}
        allowClear
      />
      <Spin spinning={loading}>
        {data && <UsageCard data={data} />}
      </Spin>
    </div>
  )
}
