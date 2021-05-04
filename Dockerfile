FROM debian:buster-slim

EXPOSE 8080

ENV CLIENT_ID ""
ENV CLIENT_SECRET ""
ENV REDIRECT_URI ""

RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates

RUN mkdir /app-container
COPY spotify-auth-service /app-container
WORKDIR /app-container

RUN chmod +x spotify-auth-service

ENTRYPOINT ["./spotify-auth-service"]
