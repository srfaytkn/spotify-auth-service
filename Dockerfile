FROM gcr.io/distroless/static-debian11:nonroot

EXPOSE 8080

ENV CLIENT_ID ""
ENV CLIENT_SECRET ""
ENV REDIRECT_URI ""

WORKDIR /app-container

COPY spotify-auth-service /app-container

ENTRYPOINT ["/app-container/spotify-auth-service"]
