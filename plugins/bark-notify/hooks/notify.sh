#!/bin/sh
# Stop hook: push a Bark notification after Claude finishes a response.
# Reads the hook JSON payload from stdin and POSTs a summary to $BARK_URL.

set -u

INPUT=$(cat)

command -v jq   >/dev/null 2>&1 || exit 0
command -v curl >/dev/null 2>&1 || exit 0

STOP_HOOK_ACTIVE=$(printf '%s' "$INPUT" | jq -r '.stop_hook_active // false')
if [ "$STOP_HOOK_ACTIVE" = "true" ]; then
    exit 0
fi

if [ -z "${BARK_URL:-}" ]; then
    exit 0
fi

CWD=$(printf '%s' "$INPUT" | jq -r '.cwd // "."')
SESSION_ID=$(printf '%s' "$INPUT" | jq -r '.session_id // ""')

PROJECT=$(basename "$CWD")
SHORT_ID=$(printf '%s' "$SESSION_ID" | cut -c1-8)

BODY="$PROJECT"
if [ -n "$SHORT_ID" ]; then
    BODY="$BODY ($SHORT_ID)"
fi

PAYLOAD=$(jq -n --arg t "Claude 回复完成" --arg b "$BODY" '{title:$t, body:$b}')

curl -s -m 5 -X POST "${BARK_URL%/}" \
    -H "Content-Type: application/json; charset=utf-8" \
    -d "$PAYLOAD" > /dev/null 2>&1 || true

exit 0
