import { api } from './cliente'
import type { UsuarioAdministrado } from '../types/usuario'

export const usuariosApi = {
  listar: () => api.obtener<UsuarioAdministrado[]>('/users/'),
  cambiarEstado: (id: string, status: 'ACTIVE' | 'INACTIVE') =>
    api.actualizar<UsuarioAdministrado>(`/users/${id}/status`, { status }),
}
