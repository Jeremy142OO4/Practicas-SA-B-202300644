import type { InputHTMLAttributes } from 'react'

type Propiedades = InputHTMLAttributes<HTMLInputElement> & {
  etiqueta: string
  error?: string
}

export function Campo({ etiqueta, error, id, ...propiedades }: Propiedades) {
  return (
    <label className="campo" htmlFor={id}>
      <span>{etiqueta}</span>
      <input id={id} aria-invalid={Boolean(error)} aria-describedby={error ? `${id}-error` : undefined} {...propiedades} />
      {error && <small id={`${id}-error`}>{error}</small>}
    </label>
  )
}
