import { api } from './cliente'
import type { Credenciales, DatosRegistro, SesionActual, UsuarioSesion } from '../types/usuario'

export const autenticacionApi = {
  registrar: (datos: DatosRegistro) => api.publicar<UsuarioSesion>('/auth/register', datos),
  iniciar: (datos: Credenciales) => api.publicar<UsuarioSesion>('/auth/login', datos),
  sesion: () => api.obtener<SesionActual>('/auth/session'),
  cerrar: () => api.publicar<void>('/auth/logout'),
}
