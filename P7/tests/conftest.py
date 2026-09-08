from importlib.util import module_from_spec, spec_from_file_location
from pathlib import Path


RAIZ = Path(__file__).resolve().parents[2]


def cargar_modulo(nombre: str, ruta_relativa: str):
    especificacion = spec_from_file_location(nombre, RAIZ / ruta_relativa)
    modulo = module_from_spec(especificacion)
    assert especificacion and especificacion.loader
    especificacion.loader.exec_module(modulo)
    return modulo
