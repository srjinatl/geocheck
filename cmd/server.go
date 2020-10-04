package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/srjinatl/geocheck/log"
	"github.com/srjinatl/geocheck/service"
	"go.uber.org/zap"
)

const (
	applicationName = "GeoCheck"
	defaultDbDir = "data/"
	defaultDbFileName = "GeoLite2-Country.mmdb"
	defaultPort = "8080"
	envHumanReadable = "HUMANREADABLE"
	envDbDir = "DB_DIR"
	envDbName = "DB_NAME"
	envPort = "PORT"
)


func envSettingExists(key string) (exists bool){
	if os.Getenv(key) != "" {
		exists = true
	}
	return
}

func getEnvSetting(key, def string) (setting string) {
	setting = def
	val := os.Getenv(key)
	if val != "" {
		setting = val
	}
	return
}

func main() {

	// load configuration settings in from environment variables if provided
	humanReadable := envSettingExists(envHumanReadable)
	dbDir := getEnvSetting(envDbDir, defaultDbDir)
	dbFileName := getEnvSetting(envDbName, defaultDbFileName)


	// create logger for the application
	log := log.NewLogger(applicationName, humanReadable)
	log.Zap.Info("Application starting...")
	defer log.Zap.Info("Application stopping...")

	// create router for application
	router := gin.Default()
	router.RedirectTrailingSlash = false

	// create config for the service
	cfg := service.NewGeoCheckConfig().WithDbDir(dbDir).WithDbFileName(dbFileName).WithLogger(log).WithRouter(router)

	// create our service
	svc := service.NewGeoCheckService(cfg)
	err := svc.Init()
	if err != nil {
		msg := fmt.Sprintf("Unable to start %s", applicationName)
		log.Zap.Fatal(msg, zap.Error(err))
		return
	}
	defer func() {
		// shut down the service
		svc.Shutdown()
	}()

	// create the server
	srv := &http.Server{
		Addr:    port(),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Zap.Fatal("Server abnormal close error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Zap.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Zap.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	log.Zap.Info("Server exiting...")

}

func port() string {
	port := os.Getenv(envPort)
	if len(port) == 0 {
		port = defaultPort
	}
	return ":" + port
}

