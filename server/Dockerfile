FROM golang:1.22

ENV GEMINI_API_KEY="AIzaSyBScRrn0BvsSRZsdOGFNNft1YfI_28fCL8"
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o server

EXPOSE 8080

CMD ["./server"]