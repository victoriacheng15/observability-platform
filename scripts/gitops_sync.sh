#!/bin/bash
set -euo pipefail

REPO_NAME=${1:-""}
ALLOWED_REPOS=("observability-hub")
BASE_DIR="/home/server/software"

# Log helper to make Grafana/Loki filtering easy
log() {
    local level=$1
    local msg=$2
    echo "repo=${REPO_NAME:-unknown} level=${level} msg=\"${msg}\""
}

# 1. Validation Logic (Security Barrier)
if [[ -z "$REPO_NAME" ]]; then
    log "ERROR" "No repository name provided."
    exit 1
fi

IS_ALLOWED=false
for repo in "${ALLOWED_REPOS[@]}"; do
    if [[ "$repo" == "$REPO_NAME" ]]; then
        IS_ALLOWED=true
        break
    fi
done

if [[ "$IS_ALLOWED" == "false" ]]; then
    log "CRITICAL" "Repository not in allowlist. Access denied."
    exit 1
fi

REPO_PATH="${BASE_DIR}/${REPO_NAME}"

if [[ ! -d "$REPO_PATH/.git" ]]; then
    log "ERROR" "Path ${REPO_PATH} is not a valid git repository."
    exit 1
fi

# 2. Sync Logic
cd "$REPO_PATH"
if ! git fetch origin main --quiet; then
    log "ERROR" "Failed to fetch from origin. Check network/permissions."
    exit 1
fi

LOCAL_HASH=$(git rev-parse HEAD)
REMOTE_HASH=$(git rev-parse origin/main)

if [[ "$LOCAL_HASH" != "$REMOTE_HASH" ]]; then
    log "SUCCESS" "New changes detected. Synchronizing..."
    git pull origin main
else
    # INFO level can be filtered out in Grafana to reduce noise
    log "INFO" "Already in sync."
fi
