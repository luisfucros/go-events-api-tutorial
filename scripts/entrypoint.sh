#!/bin/sh
set -e

# Wait for DB to be ready
: "Waiting for DB..."
DB_HOST=${DB_HOST:=db}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-events_password}
MYSQL_DATABASE=${MYSQL_DATABASE:-events}

until mysqladmin ping -h "$DB_HOST" -uroot -p"$MYSQL_ROOT_PASSWORD" --silent; do
  echo "Waiting for database at $DB_HOST..."
  sleep 1
done

# Run migrations
echo "Running migrations..."
/usr/local/bin/migrate -path=/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(${DB_HOST}:3306)/${MYSQL_DATABASE}" -verbose up

# Start the app
exec /usr/local/bin/app
