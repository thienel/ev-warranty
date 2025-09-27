import { createRoot } from 'react-dom/client'
import App from '@/App.jsx'
import { ConfigProvider } from 'antd'
import antdTheme from '@styles/antdTheme'
import '@styles/reset.css'

createRoot(document.getElementById('root')).render(
  <ConfigProvider theme={antdTheme}>
    <App />
  </ConfigProvider>
)
