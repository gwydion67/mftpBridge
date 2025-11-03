
FROM golang:1.25.3

RUN mkdir /app 

WORKDIR /app 

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 80090

CMD [ "/app/main" ]


