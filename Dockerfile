# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY README.md ./

#sqlurl & listener addr 
RUN CGO_ENABLED=0 GOOS=linux go build -o /register_login

EXPOSE 8080
CMD [ "/register_login" ]
