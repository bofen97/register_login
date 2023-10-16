# syntax=docker/dockerfile:1

FROM golang:1.21 AS BUILD

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY README.md ./

#sqlurl & listener addr 
RUN CGO_ENABLED=0 GOOS=linux go build -o /register_login

FROM scratch
COPY --from=BUILD /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=BUILD /register_login /register_login
COPY *.jpg /
EXPOSE 8080
CMD [ "/register_login" ]
