package bootstrap

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/kooroshh/fiber-boostrap/app/ws"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"github.com/kooroshh/fiber-boostrap/pkg/env"
	"github.com/kooroshh/fiber-boostrap/pkg/router"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogFile()
	database.SetupDatabase()
	database.SetupMongoDB()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New()) // menghandle panic agar tidak error
	app.Use(logger.New())  // untuk logger
	app.Get("/dashboard", monitor.New())

	go ws.ServeWSMessaging(app) // Initialize WebSocket messaging

	router.InstallRouter(app)

	return app
}

func SetupLogFile() {
	logFile, err := os.OpenFile("./logs/simple_messaging.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Set log format
}
