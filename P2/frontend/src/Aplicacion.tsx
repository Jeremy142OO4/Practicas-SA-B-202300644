import { BrowserRouter } from 'react-router-dom'
import { ProveedorSesion } from './contexts/ContextoSesion'
import { Rutas } from './router/Rutas'
import './aplicacion.css'

function Aplicacion() {
  return <BrowserRouter><ProveedorSesion><Rutas /></ProveedorSesion></BrowserRouter>
}

export default Aplicacion
