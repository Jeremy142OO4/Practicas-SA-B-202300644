import { useState, type FormEvent } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { Campo } from '../components/forms/Campo'
import { DisenoAutenticacion } from '../components/layout/DisenoAutenticacion'
import { useSesion } from '../hooks/useSesion'

export function IniciarSesion() {
  const [correo, setCorreo] = useState('')
  const [clave, setClave] = useState('')
  const [error, setError] = useState('')
  const [enviando, setEnviando] = useState(false)
  const { iniciar, limpiarAviso } = useSesion()
  const navegar = useNavigate()
  const ubicacion = useLocation()
  const estado = ubicacion.state as { mensaje?: string; sesionTerminada?: boolean } | null
  const mensaje = estado?.mensaje

  async function enviar(evento: FormEvent) {
    evento.preventDefault(); setError(''); setEnviando(true)
    try {
      await iniciar({ email: correo.trim(), password: clave })
      limpiarAviso()
      navegar('/confirmacion', { replace: true })
    } catch (fallo) {
      setError(fallo instanceof Error ? fallo.message : 'No se pudo iniciar sesión')
    } finally { setEnviando(false) }
  }

  return <DisenoAutenticacion titulo="Bienvenido de nuevo" descripcion="Ingresa tus credenciales para continuar.">
    {mensaje && <div className="aviso exito" role="status">{mensaje}</div>}
    {estado?.sesionTerminada && <div className="aviso sesion-terminada" role="alert"><strong>Tu sesión terminó</strong><span>El tiempo de renovación finalizó. Ingresa nuevamente para continuar.</span></div>}
    {error && <div className="aviso error" role="alert">{error}</div>}
    <form className="formulario" onSubmit={enviar}>
      <Campo id="correo" etiqueta="Correo electrónico" type="email" autoComplete="email" value={correo} onChange={(e) => setCorreo(e.target.value)} required />
      <Campo id="clave" etiqueta="Contraseña" type="password" autoComplete="current-password" value={clave} onChange={(e) => setClave(e.target.value)} required />
      <button className="boton-principal" disabled={enviando}>{enviando ? 'Ingresando…' : 'Iniciar sesión'}</button>
    </form>
    <p className="alternativa">¿Aún no tienes cuenta? <Link to="/registro">Crear cuenta</Link></p>
  </DisenoAutenticacion>
}
