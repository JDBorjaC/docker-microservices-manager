import { BrowserRouter, Routes, Route} from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import { StrictMode } from 'react'

import {Dash} from './pages/Dash.tsx'
import ServiceCreator from './pages/ServiceCreator.tsx'
import { ServiceDetail } from './pages/ServiceDetail.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/admin" element={<Dash />} />
        <Route path="/edit" element={<ServiceCreator/>} />
        <Route path="/deets/:serviceId" element={<ServiceDetail/>}/>
      </Routes>
    </BrowserRouter>
  </StrictMode>
)
