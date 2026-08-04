import { createContext } from 'react'
import type { Credenciales, SesionActual } from '../types/usuario'

export type ValorSesion = {
  sesion: SesionActual | null
  cargando: boolean
  sesionTerminada: boolean
  iniciar: (datos: Credenciales) => Promise<void>
  cerrar: () => Promise<void>
  comprobar: () => Promise<void>
  limpiarAviso: () => void
}

export const ContextoSesion = createContext<ValorSesion | null>(null)
