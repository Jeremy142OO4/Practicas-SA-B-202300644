import ast
from collections import Counter
from pathlib import Path


RAIZ_PRUEBAS = Path(__file__).resolve().parents[1] / "tests"
categorias = {"unit", "integration"}
encontradas = Counter()

for ruta in RAIZ_PRUEBAS.rglob("test_*.py"):
    arbol = ast.parse(ruta.read_text(encoding="utf-8"), filename=str(ruta))
    for nodo in ast.walk(arbol):
        if not isinstance(nodo, (ast.FunctionDef, ast.AsyncFunctionDef)) or not nodo.name.startswith("test_"):
            continue
        marcas = set()
        for decorador in nodo.decorator_list:
            if isinstance(decorador, ast.Attribute) and isinstance(decorador.value, ast.Attribute):
                if isinstance(decorador.value.value, ast.Name) and decorador.value.value.id == "pytest":
                    marcas.add(decorador.attr)
        if len(marcas & categorias) != 1:
            raise SystemExit(f"{ruta}:{nodo.lineno}: cada prueba debe tener exactamente una categoría")
        encontradas.update(marcas & categorias)

esperadas = {"unit": 7, "integration": 2}
if dict(encontradas) != esperadas:
    raise SystemExit(
        "Cantidad de pruebas inválida. "
        f"se esperaba {esperadas} y se encontró {dict(encontradas)}"
    )

print(
    f"Pruebas válidas: {encontradas['unit']} unitarias y "
    f"{encontradas['integration']} de integración."
)
