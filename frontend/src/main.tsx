import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { PersistGate } from 'redux-persist/integration/react'
import App from '@/App'
import 'antd/dist/reset.css'
import '@ant-design/v5-patch-for-react-19'
import store, { persistor } from '@redux/store'
import LoadingOverlay from '@components/LoadingOverlay/LoadingOverlay'

const rootElement = document.getElementById('root')
if (!rootElement) throw new Error('Failed to find the root element')

createRoot(rootElement).render(
  <Provider store={store}>
    <PersistGate loading={<LoadingOverlay loading={true} children={null} />} persistor={persistor}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </PersistGate>
  </Provider>,
)
