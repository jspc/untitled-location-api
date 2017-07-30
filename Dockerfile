FROM alpine
MAINTAINER jspc <james@zero-internet.org.uk>

EXPOSE 8000
ADD app /app

ENTRYPOINT ["/app"]
