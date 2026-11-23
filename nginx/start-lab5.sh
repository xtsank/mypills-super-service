#!/bin/sh
set -eu

for template in /etc/nginx/templates/*.template; do
    name=${template##*/}
    envsubst '${APP_PORT} ${HEALTH_CACHE_CONTROL}' < "$template" > "/etc/nginx/conf.d/${name%.template}"
done

load-admin-password.sh
exec nginx -g 'daemon off;'
