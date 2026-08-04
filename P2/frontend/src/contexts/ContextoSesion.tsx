import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react'
import { autenticacionApi } from '../api/autenticacion'
import { EVENTO_SESION_RENOVADA, EVENTO_SESION_TERMINADA } from '../api/cliente'
import type { Credenciales, SesionActual } from '../types/usuario'
import { ContextoSesion } from './sesion'

export function ProveedorSesion({ children }: { children: ReactNode }) {
  const [sesion, setSesion] = useState<SesionActual | null>(null)
  const [cargando, setCargando] = useState(true)
  const [sesionTerminada, setSesionTerminada] = useState(false)

  const comprobar = useCallback(async () => {
    try {
      setSesion(await autenticacionApi.sesion())
    } catch {
      setSesion(null)
    } finally {
      setCargando(false)
    }
  }, [])

  useEffect(() => { void comprobar() }, [comprobar])

  useEffect(() => {
    function renovar(evento: Event) {
      setSesion((evento as CustomEvent<SesionActual>).detail)
    }
    function terminar() {
      setSesion(null)
      setSesionTerminada(true)
    }
    window.addEventListener(EVENTO_SESION_RENOVADA, renovar)
    window.addEventListener(EVENTO_SESION_TERMINADA, terminar)
    return () => {
      window.removeEventListener(EVENTO_SESION_RENOVADA, renovar)
      window.removeEventListener(EVENTO_SESION_TERMINADA, terminar)
    }
  }, [])

  useEffect(() => {
    if (!sesion) return
    const restante = sesion.renewableUntil * 1000 - Date.now()
    if (restante <= 0) {
      setSesion(null)
      setSesionTerminada(true)
      return
    }
    const temporizador = window.setTimeout(() => {
      setSesion(null)
      setSesionTerminada(true)
    }, restante)
    return () => window.clearTimeout(temporizador)
  }, [sesion])

  const iniciar = useCallback(async (datos: Credenciales) => {
    await autenticacionApi.iniciar(datos)
    setSesion(await autenticacionApi.sesion())
    setSesionTerminada(false)
  }, [])

  const cerrar = useCallback(async () => {
    try { await autenticacionApi.cerrar() } finally {
      setSesion(null)
      setSesionTerminada(false)
    }
  }, [])

  const limpiarAviso = useCallback(() => setSesionTerminada(false), [])
  const valor = useMemo(() => ({
    sesion, cargando, sesionTerminada, iniciar, cerrar, comprobar, limpiarAviso,
  }), [sesion, cargando, sesionTerminada, iniciar, cerrar, comprobar, limpiarAviso])

  return <ContextoSesion.Provider value={valor}>{children}</ContextoSesion.Provider>
}
