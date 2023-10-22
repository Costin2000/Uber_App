#base go image
FROM alpine:latest

RUN mkdir /app

COPY carApp /app

CMD [ "/app/carApp" ]