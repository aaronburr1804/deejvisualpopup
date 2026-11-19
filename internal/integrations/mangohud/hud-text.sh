#!/bin/bash
FILE="$HOME/.local/state/deej/deej-hud.txt"

[ -f "$FILE" ] || exit 0

HUD_TIMEOUT=5
CURRENT_TIME=$(date +%s)
FILE_MOD_TIME=$(date -r "$FILE" +%s)
DELTA_TIME=$((CURRENT_TIME - FILE_MOD_TIME))

if [ "$DELTA_TIME" -lt "$HUD_TIMEOUT" ]; then
    cat "$FILE"
fi


