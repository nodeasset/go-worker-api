# build stage
FROM golang:alpine AS build-env

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
ADD . /src

RUN go get github.com/streadway/amqp
RUN go get github.com/tkanos/gonfig
RUN go get github.com/gorilla/mux

RUN cd /src && go build -o go-worker-api

# final stage
FROM alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=build-env /src/go-worker-api /app/
COPY --from=build-env /src/config.json /app/
COPY --from=build-env /src/wait-for-it.sh /app/
#ENTRYPOINT ./go-worker
#CMD ["./wait-for-it.sh", "queue:5672", "-t", "45", "--", "./go-worker-api"]
CMD ["/app/go-worker-api"]