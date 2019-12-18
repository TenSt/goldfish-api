FROM golang:1.13-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o server main.go

EXPOSE 3000

CMD ["./server"]