FROM python:3.13-slim AS builder

WORKDIR /app

COPY requirements.txt .

RUN python -m venv /opt/venv \
    && /opt/venv/bin/pip install --no-cache-dir \
    -r requirements.txt \
    && rm -rf /opt/venv/lib/python3.13/site-packages/pip* \
        /opt/venv/lib/python3.13/site-packages/setuptools*


FROM python:3.13-slim

ENV PATH="/opt/venv/bin:$PATH" \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

RUN groupadd --gid 10001 appgroup \
    && useradd --uid 10001 --gid appgroup \
        --no-create-home --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=builder /opt/venv /opt/venv

COPY --chown=10001:10001 . .

USER 10001:10001

EXPOSE 4003

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "4003"]
