package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/isaqueveras/power-sso/config"
	"github.com/isaqueveras/power-sso/internal/middleware"
	"github.com/isaqueveras/power-sso/pkg/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration file: ", err)
	}

	var cfg *config.Config
	if cfg, err = config.ParseConfig(cfgFile); err != nil {
		log.Fatal("Error parsing configuration file: ", err)
	}

	logg := logger.NewLogger(cfg)
	logg.InitLogger()

	router := gin.New()
	router.Use(
		middleware.VersionInfo(),
		middleware.RecoveryWithZap(logg.ZapLogger(), true),
		middleware.GinZap(logg.ZapLogger(), *cfg),
	)

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to PowerSSO", "date": time.Now()})
	})

	group := errgroup.Group{}
	group.Go(func() error {
		return endless.ListenAndServe("0.0.0.0"+cfg.Server.Port, router)
	})

	if err = group.Wait(); err != nil {
		logg.Fatal("Error while serving the application", err)
	}
}