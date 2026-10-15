# Informe de recuperacion ante desastre

## 1. Objetivos declarados

Para una futura prueba completa se proponen RTO de hasta 30 minutos y RPO
de hasta seis horas, acorde con la frecuencia de los respaldos. Estos objetivos
no quedaron registrados antes de los ensayos descritos abajo; no se presentan
como metas ya demostradas ni como tiempos de reconstruccion del cluster.

## 2. Escenario ejecutado

El 15 de octubre de 2026 se inserto un marcador nuevo en PostgreSQL, se hizo
un backup y se elimino solamente esa fila de prueba de la base activa. Luego
se restauro en `sa-p9-restored`, con un volumen diferente. El marcador no
estaba en el script de inicializacion. No se elimino el namespace activo.

Tambien se bloqueo temporalmente la programacion en un nodo y se desalojaron
solo sus pods de pedidos, mientras se consultaba el gateway. Esta prueba
comprueba reprogramacion selectiva, no la perdida de una maquina completa.

## 3. Tiempos medidos

Todos los horarios siguientes son UTC, del 15 de octubre de 2026.

| Evento | Hora o resultado |
|---|---|
| Backup completado | 23:05:31; 141/141 recursos |
| Eliminacion de la fila de prueba | 23:05:42 |
| Inicio del restore | 23:05:47 |
| Restore completado | 23:06:29; 3/3 recursos seleccionados |
| Dato recuperado verificado | 23:06:37 |
| Recuperacion aislada | 50 s desde el inicio del restore; 55 s desde la eliminacion |
| Segunda prueba de reprogramacion | 23:24:19 a 23:25:10; aproximadamente 52 s |

La primera prueba de reprogramacion tuvo 22 consultas correctas de 30.
Se agregaron 15 segundos de espera antes de terminar pedidos y un plazo
de terminacion de 45 segundos. Al repetirla, las 30 consultas fueron
correctas y las dos replicas quedaron listas en nodos distintos. En ambos
casos el nodo volvio a quedar habilitado.

## 4. Perdida medida

Se recupero `p9-verificacion|dato-recuperado-20261015T230448Z` desde el backup.
No hubo perdida del marcador respaldado. Esto verifica ese dato concreto,
no garantiza RPO cero para todas las operaciones posteriores al respaldo.

## 5. Puntos unicos de fallo detectados

PostgreSQL y RabbitMQ tienen una sola instancia. No se afirma alta
disponibilidad de esos componentes. El primer ensayo tambien mostro errores
durante el reemplazo de pedidos; el segundo no los reprodujo despues del
ajuste de terminacion. Las capturas historicas de observabilidad dependen
de su retencion y almacenamiento.

## 6. Brecha y plan

Falta medir una destruccion y reconstruccion completa y una perdida de nodo
completa. No se ejecutaron porque el cluster tambien sirve a P10. Deben
programarse en un entorno separado o en una ventana autorizada, declarando
antes los objetivos y verificando un backup utilizable. Replicar bases y
mensajeria exige configurar su replicacion; aumentar replicas sin ella no
resuelve el problema.

Registros reales: [restauracion](../evidencias/restauracion.txt),
[primer ensayo con errores](../evidencias/reprogramacion-inicial.txt) y
[segundo ensayo correcto](../evidencias/reprogramacion.txt).
