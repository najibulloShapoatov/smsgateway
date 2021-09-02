package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	routes "smsc/mini/router"
	"smsc/pkg/config"
	"smsc/pkg/db"
	"smsc/pkg/log"
	"strings"
	"syscall"
	"time"
)

//main - run func
func main() {
	confFile := getCurrentDir() + "default.conf"
	config.Init(confFile)
	log.InitW(getCurrentDir() +fmt.Sprintf("mini_%s.log", time.Now().Format("20060102150405")))
	db.SetConfigDB(config.LoadDBConfigs())
	db.Init()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	var router = routes.Init(getCurrentDir())

	log.Printf("Starting Server on address http://%s:%s", config.GetHTTPHost(), config.GetHTTPPort())

	srv := &http.Server{
		Addr: config.GetHTTPHost() + ":" + config.GetHTTPPort(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func getCurrentDir() string {

	dir, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	var ss []string
	if runtime.GOOS == "windows" {
		ss = strings.Split(dir, "\\")
		return strings.Join(ss[0:len(ss)-1], "\\") + "\\"
	} else {
		ss = strings.Split(dir, "/")
		return strings.Join(ss[0:len(ss)-1], "/") + "/"
	}
}
