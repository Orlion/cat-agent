package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/config"
	"github.com/Orlion/cat-agent/handler"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/server"
	"github.com/Orlion/cat-agent/status"
)

var confFilename string

func init() {
	flag.StringVar(&confFilename, "conf", "", "please enter a configuration file name")
}

func main() {
	flag.Parse()

	conf, err := config.ParseConfig(confFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration file parse error: "+err.Error())
		os.Exit(1)
	}

	log.Init(conf.Log)

	err = cat.Init(conf.Cat)
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration file parse error: "+err.Error())
		os.Exit(1)
	}

	status.Init()

	srv := createServer(conf.Server)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, "server listen and serve error: "+err.Error())
			os.Exit(1)
		}
	}()

	waitGracefulStop(srv)
}

func createServer(config *server.Config) *server.Server {
	srv := server.NewServer(config)
	srv.Handle(server.CmdCreateMessageId, handler.CreateMessageId)
	srv.Handle(server.CmdSendMessage, handler.SendMessage)
	return srv
}

func waitGracefulStop(srv *server.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Infof("received signal: %s will stop...", s.String())
			ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
			srv.Shutdown(ctx)
			cat.Shutdown()
			log.Shutdown()
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
