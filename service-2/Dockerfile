# BUN DOCKER IMAGE https://hub.docker.com/r/oven/bun/tags

FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json .
RUN npm ci
COPY . .
RUN npm run build
RUN npm prune --production

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .
EXPOSE 8081
ENV NODE_ENV=production
CMD ["node", "build"]