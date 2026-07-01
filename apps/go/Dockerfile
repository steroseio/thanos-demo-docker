# Standard-library only: no dependencies to download, so the build cannot fail
# on a module proxy or checksum database. Produces a static binary that runs on
# scratch (no final image to pull).
FROM golang:1.22 AS build
WORKDIR /src
COPY main.go .
RUN go mod init demo && CGO_ENABLED=0 go build -o /app .

FROM scratch
COPY --from=build /app /app
EXPOSE 2112
ENTRYPOINT ["/app"]
