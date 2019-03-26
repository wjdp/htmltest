ARG GO_VERSION=1.11
ARG TARGET=alpine:3.9

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY ./ ./
RUN CGO_ENABLED=0 go build -installsuffix 'static' -ldflags "-X main.date=`date -u +%Y-%m-%dT%H:%M:%SZ` -X main.version=`git describe --tags`" -o /app .

FROM ${TARGET} AS final

WORKDIR /bin
COPY --from=builder /app ./htmltest
WORKDIR /test
CMD [ "htmltest", "./"]