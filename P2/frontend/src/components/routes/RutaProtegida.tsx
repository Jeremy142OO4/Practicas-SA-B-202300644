import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useSesion } from '../../hooks/useSesion'

export function RutaProtegida() {
  const { sesion, cargando, sesionTerminada } = useSesion()
  const ubicacion = useLocation()
  if (cargando) return <main className="pantalla-carga"><span /><p>Comprobando sesión…</p></main>
  if (!sesion) return <Navigate to="/iniciar-sesion" replace state={{ desde: ubicacion.pathname, sesionTerminada }} />
  return <Outlet />
}
