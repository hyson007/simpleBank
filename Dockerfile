# BUILD STAGE
FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
# first dot is relative path in local directory, 
# second dot is the relative path inside the WORK DIR
COPY . .
# in this case, the main is inside the /app directory
RUN go build -o main main.go

# if we use this example, then the binary is saved into root directory
# RUN go build -o /docker-gs-ping


# RUN STAGE
FROM alpine:3.16 
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

# EXPOSE doesn't do anything except for documentation
EXPOSE 3000
CMD ["/app/main"]