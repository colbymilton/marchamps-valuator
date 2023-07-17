FROM node:latest

WORKDIR /app
COPY ./web .
RUN npm ci
ENTRYPOINT CHOKIDAR_USEPOLLING=true npm run dev -- --host
