package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/judwhite/go-svc/svc"
	"qiniu.com/video/config"
	"qiniu.com/video/web"
)

type program struct {
	httpServer *web.HTTPServer
	conf       struct {
		config.Config
		Port int `json:"port"`
	}
}

func (p *program) Init(env svc.Environment) error {
	var configFile string
	flag.StringVar(&configFile, "conf", "", "configure")
	flag.Parse()
	if configFile == "" {
		return errors.New("no configure file")
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &p.conf)
	if err != nil {
		return err
	}

	fmt.Printf("configure:%s", data)

	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	httpServer, err := web.NewHTTPServer(context.Background(),
		p.conf.Port,
		p.conf.Config)
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
