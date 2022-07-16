package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/johnHPX/sistemaDeNoticias/internal/routers"
	"github.com/johnHPX/sistemaDeNoticias/internal/util"
)

func main() {
	log.Println("initializing webapi")
	log.Println("reading the settings")
	c := util.NewConfigs()
	projectConfigs, err := c.ProjectConfigs()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("settings started...")
	log.Println(projectConfigs.Name)
	log.Println("initialized routes")

	if projectConfigs.Port == "" {
		projectConfigs.Port = "4083"
	}

	//init web service
	ctx := context.Background()
	wsvc := routers.NewWebService(ctx)
	wsvc.Init()
	loggedRouter := handlers.LoggingHandler(os.Stdout, wsvc.GetRouters())
	//server setup
	srv := &http.Server{
		Handler:        loggedRouter,
		Addr:           fmt.Sprintf("0.0.0.0:%s", projectConfigs.Port),
		WriteTimeout:   800 * time.Second,
		ReadTimeout:    800 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listening on port %s", projectConfigs.Port)
	log.Fatal(srv.ListenAndServe())
}
