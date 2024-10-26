FROM alpine:latest

RUN mkdir /app

COPY adminApp /app

CMD ["/app/adminApp"]
