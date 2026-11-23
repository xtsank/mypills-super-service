FROM node:22.23.4-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
ARG VITE_BASE_PATH=/legacy/
ARG VITE_API_BASE_URL
RUN npm run build

FROM nginx:1.30.6-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
ARG VITE_BASE_PATH=/legacy/
COPY --from=build /app/dist/ /usr/share/nginx/html${VITE_BASE_PATH}

