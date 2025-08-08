FROM golang:latest AS go_builder

# Set the working directory and build
WORKDIR /build
RUN mkdir /build/build
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go mod download
COPY api api
COPY database database
COPY lib lib
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" bash -c ". /root/.bashrc && go build -o ./build/"
COPY ./database/migrations /build/migrations

FROM oven/bun:latest AS bun_builder
WORKDIR /build
RUN mkdir /build/build
COPY ./client ./client
ENV NODE_ENV=production
RUN cd ./client && bun install && bun run build --emptyOutDir

FROM debian:stable-slim AS reesource_tracker
# Set the working directory and run
WORKDIR /app
RUN apt-get update && apt-get install ca-certificates -y
RUN update-ca-certificates
COPY --from=go_builder /build/build .
COPY --from=bun_builder /build/build .
RUN mkdir /app/database
COPY --from=go_builder /build/migrations /app/migrations
EXPOSE 80
ENTRYPOINT ["/app/reesource-tracker"]
