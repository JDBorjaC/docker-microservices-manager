import { BrowserRouter, Routes, Route} from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import { StrictMode } from 'react'

import Dash from './pages/Dash.tsx'
import ServiceEditor from './pages/ServiceEditor.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/admin" element={<Dash />} />
        <Route path="/edit" element={<ServiceEditor/>} />
      </Routes>
    </BrowserRouter>
  </StrictMode>
)
