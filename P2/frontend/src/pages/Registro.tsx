import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { autenticacionApi } from '../api/autenticacion'
import { Campo } from '../components/forms/Campo'
import { DisenoAutenticacion } from '../components/layout/DisenoAutenticacion'

export function Registro() {
  const [nombre, setNombre] = useState(''); const [correo, setCorreo] = useState(''); const [clave, setClave] = useState('')
  const [error, setError] = useState(''); const [enviando, setEnviando] = useState(false)
  const navegar = useNavigate()

  async function enviar(evento: FormEvent) {
    evento.preventDefault(); setError('')
    if (clave.length < 8) { setError('La contraseña debe tener al menos 8 caracteres'); return }
    setEnviando(true)
    try {
      await autenticacionApi.registrar({ name: nombre.trim(), email: correo.trim(), password: clave })
      navegar('/iniciar-sesion', { replace: true, state: { mensaje: 'Cuenta creada. Ya puedes iniciar sesión.' } })
    } catch (fallo) { setError(fallo instanceof Error ? fallo.message : 'No se pudo crear la cuenta') }
    finally { setEnviando(false) }
  }

  return <DisenoAutenticacion titulo="Crea tu cuenta" descripcion="Completa los datos para registrarte como cliente.">
    {error && <div className="aviso error" role="alert">{error}</div>}
    <form className="formulario" onSubmit={enviar}>
      <Campo id="nombre" etiqueta="Nombre completo" autoComplete="name" minLength={2} value={nombre} onChange={(e) => setNombre(e.target.value)} required />
      <Campo id="correo" etiqueta="Correo electrónico" type="email" autoComplete="email" value={correo} onChange={(e) => setCorreo(e.target.value)} required />
      <Campo id="clave" etiqueta="Contraseña" type="password" autoComplete="new-password" minLength={8} value={clave} onChange={(e) => setClave(e.target.value)} required />
      <button className="boton-principal" disabled={enviando}>{enviando ? 'Creando cuenta…' : 'Crear cuenta'}</button>
    </form>
    <p className="alternativa">¿Ya tienes cuenta? <Link to="/iniciar-sesion">Iniciar sesión</Link></p>
  </DisenoAutenticacion>
}
