FROM oven/bun AS build
WORKDIR /app
COPY bun.lockb .
COPY package.json .
RUN bun install --frozen-lockfile
COPY . .
RUN bun build index.ts --compile --outfile api

FROM ubuntu:22.04
WORKDIR /app
COPY --from=build /app/api /app/api
EXPOSE 8081/tcp
CMD ["/app/api"]
