FROM node:22.23.4-alpine AS documentation
WORKDIR /build
COPY nginx/documentation/package.json nginx/documentation/package-lock.json ./
RUN npm ci
COPY nginx/documentation/render.mjs ./
COPY docs/web/README.md ./
RUN node render.mjs README.md > /documentation.html

FROM nginx:1.30.6-alpine AS lab4
RUN apk add --no-cache openssl
COPY --from=documentation /documentation.html /srv/web-documentation/index.html
COPY --chmod=755 nginx/load-admin-password.sh /docker-entrypoint.d/40-load-admin-password.sh

FROM openresty/openresty:1.29.2.4-1-alpine AS lab5
RUN apk add --no-cache openssl gettext-envsubst \
    && { id nginx >/dev/null 2>&1 || { addgroup -S nginx && adduser -S -D -H -G nginx nginx; }; } \
    && ln -sf /usr/local/openresty/nginx/sbin/nginx /usr/local/bin/nginx \
    && mkdir -p /etc/nginx/conf.d /var/log/nginx \
    && ln -sf /dev/stdout /var/log/nginx/access.log
COPY --from=documentation /documentation.html /srv/web-documentation/index.html
COPY --chmod=755 nginx/load-admin-password.sh /usr/local/bin/load-admin-password.sh
