#!/bin/sh

set -e
# to let script exit immediately if command return non-zero

echo "run db migration"
ls -al /app/app.env
cat /app/app.env
source /app/app.env
env
/app/migrate -path /app/migration -database "$DBSOURCE" -verbose up

echo "start the app"

# means take all para and run it
exec "$@"
