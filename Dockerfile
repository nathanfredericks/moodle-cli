FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION}" -o /out/moodle ./cmd/moodle

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/moodle /usr/local/bin/moodle

USER nonroot:nonroot
ENTRYPOINT ["moodle"]
CMD ["mcp", "serve"]
