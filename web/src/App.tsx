import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import Home from './pages/Home'
import Admin from './pages/Admin'

const { Header, Content } = Layout

function NavMenu() {
  const location = useLocation()
  return (
    <Menu
      theme="dark"
      mode="horizontal"
      selectedKeys={[location.pathname]}
      items={[
        { key: '/', icon: <SearchOutlined />, label: <Link to="/">额度查询</Link> },
      ]}
    />
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <Layout style={{ minHeight: '100vh' }}>
        <Header style={{ display: 'flex', alignItems: 'center' }}>
          <div style={{ color: '#fff', fontSize: 18, fontWeight: 'bold', marginRight: 24 }}>
            OneProxy
          </div>
          <NavMenu />
        </Header>
        <Content style={{ padding: '24px 48px' }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/admin" element={<Admin />} />
          </Routes>
        </Content>
      </Layout>
    </BrowserRouter>
  )
}
