import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
//import './styles/monitor.css'
import Dash from './Dash.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Dash/>
  </StrictMode>,
)
