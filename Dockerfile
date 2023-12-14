FROM golang:1.21

WORKDIR /

COPY . .

RUN go mod download

# Install the package
RUN go build -o ./app ./cmd/.

ARG PORT

EXPOSE ${PORT}

# Run the executable
CMD ["./app"]