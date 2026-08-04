import { Navigate, Route, Routes } from 'react-router-dom'
import { RutaProtegida } from '../components/routes/RutaProtegida'
import { Confirmacion } from '../pages/Confirmacion'
import { IniciarSesion } from '../pages/IniciarSesion'
import { NoEncontrada } from '../pages/NoEncontrada'
import { Registro } from '../pages/Registro'
import { RutasProtegidas } from '../pages/RutasProtegidas'
import { AdministrarUsuarios } from '../pages/AdministrarUsuarios'

export function Rutas() {
  return <Routes>
    <Route path="/" element={<Navigate to="/iniciar-sesion" replace />} />
    <Route path="/iniciar-sesion" element={<IniciarSesion />} />
    <Route path="/registro" element={<Registro />} />
    <Route element={<RutaProtegida />}>
      <Route path="/confirmacion" element={<Confirmacion />} />
      <Route path="/rutas-protegidas" element={<RutasProtegidas />} />
      <Route path="/usuarios" element={<AdministrarUsuarios />} />
    </Route>
    <Route path="*" element={<NoEncontrada />} />
  </Routes>
}
