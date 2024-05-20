#!/bin/bash

TMP_DIR="${1:-"./tmp"}"

LINE_FILE="line"
LOG_FILE="${2:-"tailwind.log"}"

LAST_LINE=0
CURRENT_LINES=0

if [ ! -d "$TMP_DIR" ]; then
    mkdir "$TMP_DIR"
fi

if [ -f "$TMP_DIR/$LINE_FILE" ]; then
    LAST_LINE="$(cat "$TMP_DIR/$LINE_FILE")"
fi

if [ -f "$TMP_DIR/$LOG_FILE" ]; then
    CURRENT_LINES=$(wc -l < "$TMP_DIR/$LOG_FILE")
fi

if [ "$LAST_LINE" -ge "$CURRENT_LINES" ]; then
    echo "Waiting"
    while read f; do 
        if [ "$f" = "$LOG_FILE" ] && [ $(($(wc -l < "$TMP_DIR/$LOG_FILE") % 4)) -eq 0 ]; then 
            pkill -P $$ inotifywait
            break
        fi
    done < <(inotifywait -m -q -e create,modify --format %f "$TMP_DIR")
fi

echo "$(wc -l < "$TMP_DIR/$LOG_FILE")" > "$TMP_DIR/$LINE_FILE"
