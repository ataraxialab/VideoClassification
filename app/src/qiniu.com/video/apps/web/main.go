package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/judwhite/go-svc/svc"
	"qiniu.com/video/builder"
	"qiniu.com/video/mq"
	"qiniu.com/video/server"
	"qiniu.com/video/web"
)

type program struct {
	httpServer *web.HTTPServer
	mq         mq.MQ
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Stop() error {
	if p.mq != nil {
		p.mq.Close()
	}
	if p.httpServer != nil {
		p.httpServer.Close()
	}
	return nil
}

func (p *program) Start() error {
	mq := &mq.EmbeddedMQ{}
	if err := mq.Open(); err != nil {
		return err
	}

	p.mq = mq

	server, err := server.CreateServer(builder.Cmd, mq)
	if err != nil {
		mq.Close()
		return err
	}

	httpServer, err := web.NewHTTPServer(context.Background(), 8000, server)
	if err != nil {
		mq.Close()
		return err
	}

	p.httpServer = httpServer
	if p.httpServer != nil {
		p.httpServer.Serve()
	}

	return nil
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}
