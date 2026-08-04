import type { ReactNode } from 'react'

export function DisenoAutenticacion({ titulo, descripcion, children }: { titulo: string; descripcion: string; children: ReactNode }) {
  return (
    <main className="pagina-autenticacion">
      <section className="presentacion">
        <span className="marca">Práctica 2 · Software Avanzado</span>
        <h1>Acceso seguro, sin complicaciones.</h1>
        <p>Autenticación con JWT, cookies HTTP-only y autorización basada en roles.</p>
        <div className="tecnologias"><span>React</span><span>Go</span><span>PostgreSQL</span></div>
      </section>
      <section className="panel-formulario">
        <div className="encabezado-formulario"><h2>{titulo}</h2><p>{descripcion}</p></div>
        {children}
      </section>
    </main>
  )
}
