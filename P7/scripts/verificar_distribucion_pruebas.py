import ast
from collections import Counter
from pathlib import Path


RAIZ_PRUEBAS = Path(__file__).resolve().parents[1] / "tests"
categorias = {"unit", "integration", "e2e"}
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

total = sum(encontradas.values())
esperadas = {"unit": total * 70 // 100, "integration": total * 20 // 100, "e2e": total * 10 // 100}
if total == 0 or total % 10 or dict(encontradas) != esperadas:
    raise SystemExit(
        "Distribución inválida: el total debe ser múltiplo de 10 y conservar 70/20/10; "
        f"se esperaba {esperadas} y se encontró {dict(encontradas)}"
    )

print(
    f"Distribución válida: {encontradas['unit']} unitarias (70%), "
    f"{encontradas['integration']} integración (20%), {encontradas['e2e']} E2E (10%)."
)
