import { Link, useNavigate } from 'react-router-dom'
import { useSesion } from '../hooks/useSesion'

export function Confirmacion() {
  const { sesion, cerrar } = useSesion()
  const navegar = useNavigate()

  async function salir() {
    await cerrar()
    navegar('/iniciar-sesion', { replace: true })
  }

  return <main className="pagina-confirmacion">
    <section className="confirmacion">
      <div className="icono-correcto">✓</div>
      <span className="marca">Sesión verificada</span>
      <h1>Inicio de sesión exitoso</h1>
      <p>Has ingresado correctamente. Tu sesión está activa y ya puedes comprobar las rutas disponibles para tu rol.</p>
      <div className="rol-confirmacion">
        <span>Ingresaste como</span>
        <strong>{sesion?.role}</strong>
      </div>
      <div className="acciones-confirmacion">
        <Link className="boton-principal" to="/rutas-protegidas">Ver rutas disponibles</Link>
        <button className="boton-secundario" onClick={() => void salir()}>Cerrar sesión</button>
      </div>
    </section>
  </main>
}
