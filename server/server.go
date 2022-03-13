package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/srrmendez/private-api-offers/conf"
	"github.com/srrmendez/private-api-offers/docs"
	pkgHttp "github.com/srrmendez/services-interface-tools/pkg/http"
	log "github.com/srrmendez/services-interface-tools/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Env struct {
}

var env Env

func Init() {
	// Swagger doc
	docs.SwaggerInfo.Title = "Private offers API"
	docs.SwaggerInfo.Description = "Private offers api for ETECSA."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = conf.GetProps().App.Path
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	//Open log file
	f, err := os.OpenFile(conf.GetProps().App.LogAddress, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Creating  logger
	lg := log.NewLogger(log.Config{
		Level:     log.Info,
		Formatter: &logrus.JSONFormatter{},
		Output:    f,
	})

	ctx := context.TODO()

	mongoAddr := fmt.Sprintf("mongodb://%s:%d", conf.GetProps().Database.Host, conf.GetProps().Database.Port)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoAddr))
	if err != nil {
		panic(err)
	}

	defer mongoClient.Disconnect(ctx)

	env = Env{}

	// Creating http logger
	l := log.NewLogger(log.Config{
		Level:     log.Info,
		Formatter: &logrus.TextFormatter{},
		Output:    os.Stdout,
	})

	router := pkgHttp.MapRoutes(Routes, conf.GetProps().App.Path, l)

	router.PathPrefix(fmt.Sprintf("%s/api-docs/", conf.GetProps().App.Path)).Handler(httpSwagger.WrapHandler)

	port := fmt.Sprintf(":%d", conf.GetProps().App.Port)

	server := http.Server{
		Addr:         port,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
