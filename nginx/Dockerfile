FROM node:22.23.4-alpine AS documentation
WORKDIR /build
COPY nginx/documentation/package.json nginx/documentation/package-lock.json ./
RUN npm ci
COPY nginx/documentation/render.mjs ./
COPY docs/web/README.md ./
RUN node render.mjs README.md > /documentation.html

FROM nginx:1.30.6-alpine
RUN apk add --no-cache openssl
COPY --from=documentation /documentation.html /srv/web-documentation/index.html
COPY --chmod=755 nginx/load-admin-password.sh /docker-entrypoint.d/40-load-admin-password.sh
