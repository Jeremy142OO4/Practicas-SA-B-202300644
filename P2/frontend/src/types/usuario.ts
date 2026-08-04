export type UsuarioSesion = {
  id: string
  name: string
  email: string
  role: 'ADMIN' | 'CLIENTE'
  status: 'ACTIVE' | 'INACTIVE' | 'BLOCKED'
}

export type SesionActual = {
  id: string
  role: UsuarioSesion['role']
  expiresAt: number
  renewableUntil: number
}

export type UsuarioAdministrado = UsuarioSesion & {
  createdAt: string
  updatedAt: string
}

export type Credenciales = {
  email: string
  password: string
}

export type DatosRegistro = Credenciales & {
  name: string
}
