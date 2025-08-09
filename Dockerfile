FROM golang:1.24

WORKDIR /usr/src/app

# Install pre-reqs for seeding
RUN apt-get update -y

# Install the ca certificate for dev
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN rm -f .env

RUN go build .

EXPOSE 5995

CMD ["./rcagent", "-D"]