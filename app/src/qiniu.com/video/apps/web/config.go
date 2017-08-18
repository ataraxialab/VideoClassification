package main

import "qiniu.com/video/config"

type conf struct {
	config.Config
	Port int `json:"port"`
}
