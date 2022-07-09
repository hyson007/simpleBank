# BUILD STAGE
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
# first dot is relative path in local directory, 
# second dot is the relative path inside the WORK DIR
COPY . .
# in this case, the main is inside the /app directory
RUN go build -o main main.go
# install curl
RUN apk add curl
# install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin/


# if we use this example, then the binary is saved into root directory
# RUN go build -o /docker-gs-ping


# RUN STAGE
FROM alpine:3.16 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration


# EXPOSE doesn't do anything except for documentation
EXPOSE 3000

# when CMD is used together with ENTRYPOINT, CMD will acts as additonal paramters

CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]