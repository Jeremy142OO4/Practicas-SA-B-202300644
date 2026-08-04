import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { protegidasApi } from '../api/protegidas'
import { FalloApi } from '../types/api'
import { useSesion } from '../hooks/useSesion'
import { TemporizadorSesion } from '../components/session/TemporizadorSesion'
import { Link } from 'react-router-dom'

type TipoRuta = 'admin' | 'shared'
type Resultado = { tipo: 'permitido' | 'denegado' | 'error'; texto: string }

const rutas = [
  {
    tipo: 'shared' as const,
    icono: '↔',
    clase: 'compartido',
    permiso: 'ADMIN Y CLIENTE',
    titulo: 'Área compartida',
    descripcion: 'Disponible para todos los usuarios que hayan iniciado sesión.',
    metodo: 'GET',
    direccion: '/api/v1/protected/shared',
  },
  {
    tipo: 'admin' as const,
    icono: '◆',
    clase: 'admin',
    permiso: 'SOLO ADMIN',
    titulo: 'Área administrativa',
    descripcion: 'Reservada para usuarios con el rol de administrador.',
    metodo: 'GET',
    direccion: '/api/v1/protected/admin',
  },
]

export function RutasProtegidas() {
  const { sesion, cerrar } = useSesion()
  const [resultados, setResultados] = useState<Partial<Record<TipoRuta, Resultado>>>({})
  const [cargando, setCargando] = useState<TipoRuta | null>(null)
  const navegar = useNavigate()

  async function probar(tipo: TipoRuta) {
    setCargando(tipo)
    setResultados((actuales) => ({ ...actuales, [tipo]: undefined }))

    try {
      const datos = tipo === 'admin' ? await protegidasApi.administrador() : await protegidasApi.compartida()
      setResultados((actuales) => ({
        ...actuales,
        [tipo]: { tipo: 'permitido', texto: `El servidor autorizó el acceso con el rol ${datos.role}.` },
      }))
    } catch (error) {
      const denegado = error instanceof FalloApi && error.estado === 403
      const texto = denegado
        ? 'Tu rol actual no tiene permiso para entrar a esta ruta.'
        : error instanceof Error ? error.message : 'No se pudo comprobar la ruta.'
      setResultados((actuales) => ({
        ...actuales,
        [tipo]: { tipo: denegado ? 'denegado' : 'error', texto },
      }))
    } finally {
      setCargando(null)
    }
  }

  async function salir() {
    await cerrar()
    navegar('/iniciar-sesion', { replace: true })
  }

  return <main className="pagina-rutas">
    <header className="barra-superior">
      <div><span className="marca">Práctica 2</span><h1>Control de acceso</h1></div>
      <div className="acciones-sesion">
        <TemporizadorSesion />
        {sesion?.role === 'ADMIN' && <Link className="enlace-usuarios" to="/usuarios">Usuarios</Link>}
        <span className="insignia-rol">Rol: {sesion?.role}</span>
        <button className="cerrar-sesion" onClick={() => void salir()}>Cerrar sesión</button>
      </div>
    </header>

    <section className="contenido-rutas">
      <div className="introduccion">
        <p className="sobrelinea">Autorización por roles</p>
        <h2>¿A qué rutas tienes acceso?</h2>
        <p>Prueba cada opción. El servidor comprobará tu sesión mediante una cookie HTTP-only; el JWT permanece protegido y nunca se muestra en pantalla.</p>
      </div>

      <div className="resumen-seguridad">
        <span className="punto-seguro" aria-hidden="true" />
        <div><strong>Sesión protegida</strong><small>Autenticación activa como {sesion?.role}</small></div>
      </div>

      <div className="rejilla-rutas">
        {rutas.map((ruta) => {
          const resultado = resultados[ruta.tipo]
          const comprobando = cargando === ruta.tipo
          return <article className={`tarjeta-ruta ${resultado ? `resultado-${resultado.tipo}` : ''}`} key={ruta.tipo}>
            <div className="cabecera-ruta">
              <div className={`icono-ruta ${ruta.clase}`}>{ruta.icono}</div>
              <span className={`estado-ruta ${ruta.tipo === 'admin' ? 'exclusivo' : ''}`}>{ruta.permiso}</span>
            </div>
            <h3>{ruta.titulo}</h3>
            <p>{ruta.descripcion}</p>
            <div className="direccion-ruta"><span>{ruta.metodo}</span><code>{ruta.direccion}</code></div>

            <div className={`resultado-ruta ${resultado?.tipo ?? 'pendiente'}`} role="status">
              <span className="icono-resultado">{resultado?.tipo === 'permitido' ? '✓' : resultado?.tipo === 'denegado' ? '×' : resultado?.tipo === 'error' ? '!' : '?'}</span>
              <div>
                <strong>{resultado?.tipo === 'permitido' ? 'Acceso permitido' : resultado?.tipo === 'denegado' ? 'Acceso denegado' : resultado?.tipo === 'error' ? 'Error de conexión' : 'Sin comprobar'}</strong>
                <small>{resultado?.texto ?? 'Presiona el botón para consultar al servidor.'}</small>
              </div>
            </div>

            <button onClick={() => void probar(ruta.tipo)} disabled={Boolean(cargando)}>
              {comprobando ? 'Comprobando…' : resultado ? 'Comprobar nuevamente' : 'Comprobar acceso'}
            </button>
          </article>
        })}
      </div>
    </section>
  </main>
}
