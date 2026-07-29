-- ============================================================
-- Base de datos: Academia ficticia
-- Tabla: solicitudes_operativas
-- Motor: PostgreSQL / Supabase
-- ============================================================

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.solicitudes_operativas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    titulo VARCHAR(150) NOT NULL,

    area_solicitante VARCHAR(100) NOT NULL,

    prioridad SMALLINT NOT NULL,

    costo_estimado NUMERIC(12, 2) NOT NULL DEFAULT 0.00,

    estado VARCHAR(20) NOT NULL DEFAULT 'registrada',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_solicitudes_titulo_no_vacio
        CHECK (char_length(btrim(titulo)) > 0),

    CONSTRAINT chk_solicitudes_area_no_vacia
        CHECK (char_length(btrim(area_solicitante)) > 0),

    CONSTRAINT chk_solicitudes_prioridad
        CHECK (prioridad BETWEEN 1 AND 5),

    CONSTRAINT chk_solicitudes_costo_no_negativo
        CHECK (costo_estimado >= 0),

    CONSTRAINT chk_solicitudes_estado
        CHECK (estado IN ('registrada', 'en_proceso', 'finalizada'))
);

CREATE INDEX IF NOT EXISTS idx_solicitudes_operativas_estado
    ON public.solicitudes_operativas (estado);

CREATE INDEX IF NOT EXISTS idx_solicitudes_operativas_prioridad
    ON public.solicitudes_operativas (prioridad);

CREATE INDEX IF NOT EXISTS idx_solicitudes_operativas_area
    ON public.solicitudes_operativas (area_solicitante);

CREATE INDEX IF NOT EXISTS idx_solicitudes_operativas_created_at
    ON public.solicitudes_operativas (created_at DESC);

CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_solicitudes_operativas_updated_at
    ON public.solicitudes_operativas;

CREATE TRIGGER trg_solicitudes_operativas_updated_at
BEFORE UPDATE ON public.solicitudes_operativas
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();

COMMENT ON TABLE public.solicitudes_operativas IS
    'Almacena las solicitudes operativas registradas por las áreas de la academia.';

COMMENT ON COLUMN public.solicitudes_operativas.id IS
    'Identificador único UUID generado automáticamente.';

COMMENT ON COLUMN public.solicitudes_operativas.titulo IS
    'Título breve de la solicitud operativa.';

COMMENT ON COLUMN public.solicitudes_operativas.area_solicitante IS
    'Área de la academia que realiza la solicitud.';

COMMENT ON COLUMN public.solicitudes_operativas.prioridad IS
    'Nivel de prioridad entre 1 y 5.';

COMMENT ON COLUMN public.solicitudes_operativas.costo_estimado IS
    'Costo estimado de la solicitud con dos decimales.';

COMMENT ON COLUMN public.solicitudes_operativas.estado IS
    'Estado actual: registrada, en_proceso o finalizada.';

COMMIT;