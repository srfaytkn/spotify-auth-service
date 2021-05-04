package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	TOKEN_URL   = "https://accounts.spotify.com/api/token"
	SERVER_HOST = ":8080"
)

var RestyClient = resty.New()

type TokenSwapRequest struct {
	Code string `json:"code"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SpotifyAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type SpotifyErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func getAuthorizationHeader() string {
	rawStr := fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	base64Str := base64.StdEncoding.EncodeToString([]byte(rawStr))

	return "Basic " + base64Str
}

func healthHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func tokenSwapHandler(c *fiber.Ctx) error {
	tokenSwapRequest := &TokenSwapRequest{}

	if err := c.BodyParser(tokenSwapRequest); err != nil {
		zap.S().Errorw("TokenSwapRequest parse error", "error", err)
		return c.
			Status(400).
			JSON(&ErrorResponse{Message: "invalid request"})
	}

	spotifyAuthResponse := &SpotifyAuthResponse{}
	spotifyErrorResponse := &SpotifyErrorResponse{}

	res, err := RestyClient.R().
		SetHeader("Authorization", getAuthorizationHeader()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetResult(spotifyAuthResponse).
		SetError(spotifyErrorResponse).
		SetQueryParams(map[string]string{
			"grant_type":   "authorization_code",
			"code":         tokenSwapRequest.Code,
			"redirect_uri": os.Getenv("SPOTIFY_REDIRECT_URI"),
		}).
		Post(TOKEN_URL)

	if err != nil {
		zap.S().Errorw("tokenSwap spotify request error", "error", err)
		return c.
			Status(400).
			JSON(&ErrorResponse{Message: "spotify request error"})
	}

	if res.StatusCode() != 200 {
		zap.S().Errorw("tokenSwap spotify response error", "error", spotifyErrorResponse)
		return c.
			Status(res.StatusCode()).
			JSON(&ErrorResponse{
				Message: fmt.Sprintf("%s && %s", spotifyErrorResponse.Error, spotifyErrorResponse.ErrorDescription),
			})
	}

	zap.S().Info("token swap success")
	return c.Status(200).JSON(spotifyAuthResponse)
}

func tokenRefreshHandler(c *fiber.Ctx) error {
	tokenRefreshRequest := &TokenRefreshRequest{}

	if err := c.BodyParser(tokenRefreshRequest); err != nil {
		zap.S().Errorw("TokenRefreshRequest parse error", "error", err)
		return c.
			Status(400).
			JSON(&ErrorResponse{Message: "invalid request"})
	}

	spotifyAuthResponse := &SpotifyAuthResponse{}
	spotifyErrorResponse := &SpotifyErrorResponse{}

	res, err := RestyClient.R().
		SetHeader("Authorization", getAuthorizationHeader()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetResult(spotifyAuthResponse).
		SetError(spotifyErrorResponse).
		SetQueryParams(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": tokenRefreshRequest.RefreshToken,
		}).
		Post(TOKEN_URL)

	if err != nil {
		zap.S().Errorw("tokenRefresh spotify request error", "error", err)
		return c.
			Status(400).
			JSON(&ErrorResponse{Message: "spotify request error"})
	}

	if res.StatusCode() != 200 {
		zap.S().Errorw("tokenRefresh spotify response error", "error", spotifyErrorResponse)
		return c.
			Status(res.StatusCode()).
			JSON(&ErrorResponse{
				Message: fmt.Sprintf("%s && %s", spotifyErrorResponse.Error, spotifyErrorResponse.ErrorDescription),
			})
	}

	//spotify returns empty refresh_token
	spotifyAuthResponse.RefreshToken = tokenRefreshRequest.RefreshToken

	zap.S().Info("token refresh success")
	return c.Status(200).JSON(spotifyAuthResponse)
}

func initLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()

	app := fiber.New()

	app.Get("/health", healthHandler)
	app.Post("/callback", tokenSwapHandler)
	app.Post("/refresh", tokenRefreshHandler)

	go func() {
		log.Fatal(app.Listen(SERVER_HOST))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	zap.S().Info("Gracefully shutting down...")

	_ = app.Shutdown()

	zap.S().Info("Successful shutdown")
}
