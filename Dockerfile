FROM golang:1.22-alpine
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
RUN go build -o /watch-now
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=0 /watch-now /watch-now
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD [ "/watch-now" ]
