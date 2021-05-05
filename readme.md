### Spotify Remote Player Auth Service

```shell
docker run -d --restart=always \
  --name spotify-auth-service \
  --hostname spotify-auth-service \
  -p 8080:8080 \
  -e SPOTIFY_CLIENT_ID="cannot_be_empty" \
  -e SPOTIFY_CLIENT_SECRET="cannot_be_empty" \
  -e SPOTIFY_REDIRECT_URI="cannot_be_empty" \
  srfaytkn/spotify-auth-service
```


