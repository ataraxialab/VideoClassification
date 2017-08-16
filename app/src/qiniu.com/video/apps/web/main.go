package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/judwhite/go-svc/svc"
	"qiniu.com/video/web"
)

type program struct {
	httpServer *web.HTTPServer
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	httpServer, err := web.NewHTTPServer(context.Background(), 8000)
	if err != nil {
		return err
	}

	p.httpServer = httpServer
	httpServer.Serve()

	return nil
}

func (p *program) Stop() error {
	if p.httpServer != nil {
		p.httpServer.Close()
	}
	return nil
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}
