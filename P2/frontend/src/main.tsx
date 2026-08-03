import { StrictMode as ModoEstricto } from "react";
import { createRoot as crearRaiz } from "react-dom/client";
import Aplicacion from "./Aplicacion.tsx";
import "./estilos.css";

crearRaiz(document.getElementById("raiz")!).render(
  <ModoEstricto>
    <Aplicacion />
  </ModoEstricto>,
);
