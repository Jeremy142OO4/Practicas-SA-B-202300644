import { useContext } from 'react'
import { ContextoSesion } from '../contexts/sesion'

export function useSesion() {
  const contexto = useContext(ContextoSesion)
  if (!contexto) throw new Error('useSesion debe utilizarse dentro de ProveedorSesion')
  return contexto
}
