import { useEffect, useState } from 'react'
import { useSesion } from '../../hooks/useSesion'

function formato(segundos: number) {
  return `00:${Math.max(0, segundos).toString().padStart(2, '0')}`
}

export function TemporizadorSesion() {
  const { sesion } = useSesion()
  const [ahora, setAhora] = useState(() => Date.now())

  useEffect(() => {
    const intervalo = window.setInterval(() => setAhora(Date.now()), 250)
    return () => window.clearInterval(intervalo)
  }, [])

  if (!sesion) return null
  const vencimiento = sesion.expiresAt * 1000
  const enGracia = ahora >= vencimiento
  const limite = (enGracia ? sesion.renewableUntil : sesion.expiresAt) * 1000
  const segundos = Math.max(0, Math.ceil((limite - ahora) / 1000))

  return <div className={`temporizador-sesion ${enGracia ? 'gracia' : ''}`} title="La sesión se renueva automáticamente al realizar una solicitud durante el periodo de gracia.">
    <span className="pulso-sesion" />
    <div>
      <small>{enGracia ? 'Renovación disponible' : 'Token activo'}</small>
      <strong>{formato(segundos)}</strong>
    </div>
  </div>
}
