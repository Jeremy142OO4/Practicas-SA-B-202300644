import { FalloApi, type ErrorApi, type RespuestaApi } from '../types/api'

const URL_API = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1').replace(/\/$/, '')

type Opciones = RequestInit & { reintentar?: boolean }

export const EVENTO_SESION_RENOVADA = 'sesion-renovada'
export const EVENTO_SESION_TERMINADA = 'sesion-terminada'

async function leer<T>(respuesta: Response): Promise<T> {
  const cuerpo = (await respuesta.json().catch(() => null)) as RespuestaApi<T> | ErrorApi | null
  if (!respuesta.ok) {
    const error = cuerpo as ErrorApi | null
    throw new FalloApi(respuesta.status, error?.message || 'No se pudo completar la solicitud', error?.error?.code)
  }
  return (cuerpo as RespuestaApi<T>).data as T
}

async function enviar<T>(ruta: string, opciones: Opciones = {}): Promise<T> {
  const { reintentar = true, ...configuracion } = opciones
  const respuesta = await fetch(`${URL_API}${ruta}`, {
    ...configuracion,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...configuracion.headers },
  })

  if (respuesta.status === 401 && reintentar) {
    const error = (await respuesta.clone().json().catch(() => null)) as ErrorApi | null
    if (error?.error?.code === 'SESSION_EXPIRED') {
      const renovacion = await fetch(`${URL_API}/auth/refresh`, {
        method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' },
      })
      if (renovacion.ok) {
        const cuerpo = await leer<unknown>(renovacion)
        window.dispatchEvent(new CustomEvent(EVENTO_SESION_RENOVADA, { detail: cuerpo }))
        return enviar<T>(ruta, { ...configuracion, reintentar: false })
      }
      window.dispatchEvent(new Event(EVENTO_SESION_TERMINADA))
    }
  }
  try {
    return await leer<T>(respuesta)
  } catch (error) {
    if (error instanceof FalloApi && error.estado === 401 && error.codigo !== 'INVALID_CREDENTIALS') {
      window.dispatchEvent(new Event(EVENTO_SESION_TERMINADA))
    }
    throw error
  }
}

export const api = {
  obtener: <T>(ruta: string) => enviar<T>(ruta),
  publicar: <T>(ruta: string, datos?: unknown) => enviar<T>(ruta, {
    method: 'POST', body: datos === undefined ? undefined : JSON.stringify(datos),
  }),
  actualizar: <T>(ruta: string, datos: unknown) => enviar<T>(ruta, {
    method: 'PATCH', body: JSON.stringify(datos),
  }),
}
