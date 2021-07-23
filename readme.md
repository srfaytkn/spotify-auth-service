# Spotify Remote Player Auth Service

## Getting Started
- [Spotify Docs](https://developer.spotify.com/documentation/general/guides/app-settings/)

## Endpoints
| Description                     | Endpoint             | Method   | Payload               | On Success                                                            | On Error                                               |
|---------------------------------|----------------------|----------|-----------------------|-----------------------------------------------------------------------|--------------------------------------------------------|
| App health service              | /health              | GET      | {}                    | 200                                                                   | 400                                                    |
| Authorization code swap service | /callback            | POST     | { code }              | 200 -> { access_token, token_type, expires_in, refresh_token, scope } | (400 &#124; SPOTIFY_ERROR_STATUS) -> { message }       |
| Token refresh service           | /refresh             | POST     | { refresh_token }     | 200 -> { access_token, token_type, expires_in, refresh_token, scope } | (400 &#124; SPOTIFY_ERROR_STATUS) -> { message }       |

## Usage
```shell
# https://hub.docker.com/r/srfaytkn/spotify-auth-service

docker run -d --restart=always \
  --name spotify-auth-service \
  --hostname spotify-auth-service \
  -p 8080:8080 \
  -e SPOTIFY_CLIENT_ID="cannot_be_empty" \
  -e SPOTIFY_CLIENT_SECRET="cannot_be_empty" \
  -e SPOTIFY_REDIRECT_URI="cannot_be_empty" \
  srfaytkn/spotify-auth-service
```


