FROM node:latest

WORKDIR /app

COPY . .

RUN npm i -g pnpm

RUN pnpm i

RUN pnpm run build