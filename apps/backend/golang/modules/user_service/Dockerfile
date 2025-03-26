FROM golang:1.20
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]