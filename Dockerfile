FROM golang:1.21.6-bullseye

RUN mkdir /code

COPY go.mod go.sum /code/

WORKDIR /code/

RUN go mod download && go mod verify

COPY . /code/

RUN go build /code/cmd/web

