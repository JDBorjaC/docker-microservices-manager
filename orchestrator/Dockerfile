FROM golang:1.25.6

WORKDIR /usr/src/app

RUN mkdir -p /usr/src/app/storage

#pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generate swagger docs inside the container so the build doesn't fail on missing /docs folder
RUN swag init -g cmd/api/main.go

RUN go build -v -o /usr/local/bin/app ./cmd/api/
CMD ["app"]
