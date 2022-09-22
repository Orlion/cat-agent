package main

import (
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

	cat.Init(conf.Cat)

	srv := createServer(conf.Server)

	go waitGracefulStop(srv)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "server listen and serve error: "+err.Error())
		os.Exit(1)
	}
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
			srv.Shutdown()
			time.Sleep(3 * time.Second)
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
