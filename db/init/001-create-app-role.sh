#!/bin/sh
set -eu

APP_PASSWORD="$(cat /run/secrets/postgres_app_password)"

psql \
  --set=ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  --set=app_password="$APP_PASSWORD" <<'EOSQL'
CREATE ROLE community_app
  WITH LOGIN
  NOSUPERUSER
  NOCREATEDB
  NOCREATEROLE
  NOINHERIT
  NOREPLICATION
  PASSWORD :'app_password';

GRANT CONNECT ON DATABASE community_onion TO community_app;
GRANT USAGE, CREATE ON SCHEMA public TO community_app;
EOSQL
