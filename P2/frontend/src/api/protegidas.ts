import { api } from './cliente'
import type { RespuestaProtegida } from '../types/rutaProtegida'

export const protegidasApi = {
  administrador: () => api.obtener<RespuestaProtegida>('/protected/admin'),
  compartida: () => api.obtener<RespuestaProtegida>('/protected/shared'),
}
