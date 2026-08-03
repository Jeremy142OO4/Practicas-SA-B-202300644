import './aplicacion.css'

function Aplicacion() {
  return (
    <main className="inicio">
      <section className="tarjeta">
        <span className="etiqueta">Software Avanzado</span>
        <h1>Autenticacion y autorizacion</h1>
        <p>
          Proyecto preparado con React, Go y PostgreSQL. Los formularios de acceso se
          agregaran en la siguiente fase.
        </p>
        <div className="estado">
          <span aria-hidden="true" />
          Estructura inicial lista
        </div>
      </section>
    </main>
  )
}

export default Aplicacion

