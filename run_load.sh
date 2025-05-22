#!/usr/bin/env bash
set -e

RATE=440/1m
DUR=5m

for i in 1 2 3; do
  PORT=$((18080 + i))
  URL="http://localhost:${PORT}/?token=495386773&user=-2104222730&config=${i}"
  OUT=load_config${i}.bin

  echo "=== Конфигурация $i (порт $PORT) ==="
  echo "URL: $URL"
  echo "Запуск Vegeta attack..."
  echo "GET $URL" | vegeta attack -rate=$RATE -duration=$DUR > $OUT
  echo "Результат в $OUT"
  echo
done

