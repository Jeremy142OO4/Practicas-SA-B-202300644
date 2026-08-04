import { useEffect, useState } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { usuariosApi } from '../api/usuarios'
import { TemporizadorSesion } from '../components/session/TemporizadorSesion'
import { useSesion } from '../hooks/useSesion'
import type { UsuarioAdministrado } from '../types/usuario'

export function AdministrarUsuarios() {
  const { sesion, cerrar } = useSesion()
  const [usuarios, setUsuarios] = useState<UsuarioAdministrado[]>([])
  const [cargando, setCargando] = useState(true)
  const [actualizando, setActualizando] = useState('')
  const [error, setError] = useState('')
  const navegar = useNavigate()

  useEffect(() => {
    usuariosApi.listar()
      .then(setUsuarios)
      .catch((fallo) => setError(fallo instanceof Error ? fallo.message : 'No se pudieron cargar los usuarios.'))
      .finally(() => setCargando(false))
  }, [])

  if (sesion?.role !== 'ADMIN') return <Navigate to="/rutas-protegidas" replace />

  async function cambiarEstado(usuario: UsuarioAdministrado) {
    const nuevoEstado = usuario.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE'
    setActualizando(usuario.id)
    setError('')
    try {
      const actualizado = await usuariosApi.cambiarEstado(usuario.id, nuevoEstado)
      setUsuarios((actuales) => actuales.map((item) => item.id === actualizado.id ? actualizado : item))
    } catch (fallo) {
      setError(fallo instanceof Error ? fallo.message : 'No se pudo actualizar el usuario.')
    } finally {
      setActualizando('')
    }
  }

  async function salir() {
    await cerrar()
    navegar('/iniciar-sesion', { replace: true })
  }

  return <main className="pagina-rutas">
    <header className="barra-superior">
      <div><span className="marca">Práctica 2</span><h1>Administración de usuarios</h1></div>
      <div className="acciones-sesion">
        <TemporizadorSesion />
        <span className="insignia-rol">Rol: ADMIN</span>
        <button className="cerrar-sesion" onClick={() => void salir()}>Cerrar sesión</button>
      </div>
    </header>

    <section className="contenido-rutas contenido-usuarios">
      <nav className="migas"><Link to="/rutas-protegidas">← Volver a rutas</Link></nav>
      <div className="introduccion">
        <p className="sobrelinea">Solo administradores</p>
        <h2>Usuarios registrados</h2>
        <p>Consulta las cuentas del sistema y controla si pueden iniciar sesión.</p>
      </div>

      {error && <div className="aviso error" role="alert">{error}</div>}
      {cargando ? <div className="cargando-tabla">Cargando usuarios…</div> :
        <div className="tabla-contenedor">
          <table className="tabla-usuarios">
            <thead><tr><th>Usuario</th><th>Rol</th><th>Estado</th><th>Registro</th><th><span className="solo-lectores">Acciones</span></th></tr></thead>
            <tbody>{usuarios.map((usuario) => {
              const esPropio = usuario.id === sesion.id
              return <tr key={usuario.id}>
                <td><strong>{usuario.name}</strong><small>{usuario.email}</small></td>
                <td><span className={`rol-tabla ${usuario.role.toLowerCase()}`}>{usuario.role}</span></td>
                <td><span className={`estado-usuario ${usuario.status === 'ACTIVE' ? 'activo' : 'inactivo'}`}><i />{usuario.status === 'ACTIVE' ? 'Activo' : 'Inactivo'}</span></td>
                <td>{new Intl.DateTimeFormat('es-GT', { dateStyle: 'medium' }).format(new Date(usuario.createdAt))}</td>
                <td><button className={`boton-estado ${usuario.status === 'ACTIVE' ? 'desactivar' : 'activar'}`} disabled={esPropio || Boolean(actualizando)} onClick={() => void cambiarEstado(usuario)} title={esPropio ? 'No puedes cambiar el estado de tu propia cuenta' : undefined}>{actualizando === usuario.id ? 'Guardando…' : esPropio ? 'Sesión actual' : usuario.status === 'ACTIVE' ? 'Desactivar' : 'Activar'}</button></td>
              </tr>
            })}</tbody>
          </table>
        </div>}
    </section>
  </main>
}
