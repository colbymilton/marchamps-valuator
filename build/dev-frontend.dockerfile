FROM node:latest

WORKDIR /app
COPY . .
ENTRYPOINT CHOKIDAR_USEPOLLING=true npm run dev -- --host
