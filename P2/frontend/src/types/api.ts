export type ErrorApi = {
  success: false
  message: string
  error: { code: string }
}

export type RespuestaApi<T> = {
  success: true
  message: string
  data?: T
}

export class FalloApi extends Error {
  estado: number
  codigo: string

  constructor(estado: number, mensaje: string, codigo = 'UNKNOWN_ERROR') {
    super(mensaje)
    this.name = 'FalloApi'
    this.estado = estado
    this.codigo = codigo
  }
}
