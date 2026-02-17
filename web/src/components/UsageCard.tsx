import { Card, Col, Descriptions, Progress, Row, Statistic, Tag } from 'antd'
import { DollarOutlined } from '@ant-design/icons'
import type { UsageData } from '../api'

interface Props {
  data: UsageData
}

export default function UsageCard({ data }: Props) {
  const usedPercent = data.total_granted > 0
    ? (data.total_used / data.total_granted) * 100
    : 0

  return (
    <Card title={data.name} style={{ marginTop: 24 }}>
      <Row gutter={24}>
        <Col span={8}>
          <Statistic
            title="可用余额"
            value={data.usd_available}
            precision={2}
            prefix={<DollarOutlined />}
            suffix="USD"
            valueStyle={{ color: '#3f8600' }}
          />
        </Col>
        <Col span={8}>
          <Statistic
            title="已使用"
            value={data.usd_used}
            precision={2}
            prefix={<DollarOutlined />}
            suffix="USD"
            valueStyle={{ color: '#cf1322' }}
          />
        </Col>
        <Col span={8}>
          <Statistic
            title="总额度"
            value={data.usd_total}
            precision={2}
            prefix={<DollarOutlined />}
            suffix="USD"
          />
        </Col>
      </Row>

      <Progress
        percent={Number(usedPercent.toFixed(1))}
        status={usedPercent > 90 ? 'exception' : 'active'}
        style={{ marginTop: 24 }}
      />

      <Descriptions column={2} style={{ marginTop: 24 }} size="small">
        <Descriptions.Item label="原始可用额度">
          {data.total_available.toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="原始已用额度">
          {data.total_used.toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="原始总额度">
          {data.total_granted.toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="无限额度">
          {data.unlimited_quota ? <Tag color="green">是</Tag> : <Tag>否</Tag>}
        </Descriptions.Item>
        {data.expires_at > 0 && (
          <Descriptions.Item label="过期时间">
            {new Date(data.expires_at * 1000).toLocaleString()}
          </Descriptions.Item>
        )}
      </Descriptions>
    </Card>
  )
}
