package main

import (
	"github.com/LF-Engineering/insights-datasource-shared/http"
	"github.com/asaberwd/event-messages/api"
	"github.com/asaberwd/event-messages/internal/auth"
	"github.com/asaberwd/event-messages/internal/event"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"time"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	httpProvider := http.NewClientProvider(time.Second * 60)
	authProvider := auth.NewProvider(httpProvider, os.Getenv("AUTH0AUDIENCE"))
	e.Use(authProvider.JWTAuth)

	eventMgr := event.NewManager()
	eventHandler := api.NewEventHandler(*eventMgr)
	api.Router(e, eventHandler)

	adapter := echoadapter.New(e)
	lambda.Start(adapter.Proxy)
}
