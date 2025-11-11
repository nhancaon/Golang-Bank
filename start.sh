#!/bin/sh
set -e

echo "Running migrations..."
# source /app/app.env dùng cái này để load cấu hình trong env cx dc
# hoặc docker run -p 8080:8080 -e DB_SOURCE=... để truyền biến môi trường 
# sau khi pull image về

echo "Starting app..."
exec "$@"
