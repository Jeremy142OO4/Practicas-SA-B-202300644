@echo off
mode con cols=132 lines=46
title P10 - Evidencias reales por SSH
if "%1"=="telemetria" goto telemetria
if "%1"=="pods" goto pods
if "%1"=="incidente" goto incidente
ssh -tt -o IdentitiesOnly=yes -i C:\Users\Jeremy\.ssh\id_ed25519_codex_fedora jeremy@192.168.1.23 "set -x; kubectl config current-context; kubectl get nodes; kubectl get applications -n argocd; kubectl get pods -n sa-p9; kubectl get rollout api-gateway -n sa-p9; set +x; exec bash"
goto :eof
:telemetria
ssh -tt -o IdentitiesOnly=yes -i C:\Users\Jeremy\.ssh\id_ed25519_codex_fedora jeremy@192.168.1.23 "python3 ~/Escritorio/SA/Practicas-SA-B-202300644/P10/scripts/validar_observabilidad.py; exec bash"
goto :eof
:pods
ssh -tt -o IdentitiesOnly=yes -i C:\Users\Jeremy\.ssh\id_ed25519_codex_fedora jeremy@192.168.1.23 "set -x; kubectl get pods -n observability; kubectl get svc -n observability; set +x; exec bash"
goto :eof
:incidente
ssh -tt -o IdentitiesOnly=yes -i C:\Users\Jeremy\.ssh\id_ed25519_codex_fedora jeremy@192.168.1.23 "set -x; kubectl get analysisruns -n sa-p9; kubectl get rollout api-gateway -n sa-p9; set +x; python3 ~/Escritorio/SA/Practicas-SA-B-202300644/P10/scripts/resumen_incidente.py; exec bash"
