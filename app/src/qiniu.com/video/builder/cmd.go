package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"qiniu.com/video/logger"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

const (
	// Cmd call external command to build data
	Cmd        Implement = "cmd"
	outputRoot string    = "output"
)

type cmdFrameRandom struct {
	logger *logger.Logger
}

func createOutputDir() (string, error) {
	d := path.Join(outputRoot, time.Now().Format("20060102-150405"))
	return d, os.MkdirAll(d, os.ModePerm)
}

func (f cmdFrameRandom) Build(params interface{}) ([]interface{}, error) {
	p, ok := params.(Params)
	if !ok {
		f.logger.Errorf("unknown params:%v", p)
		return nil, fmt.Errorf("unknown params")
	}

	output, e := createOutputDir()
	if e != nil {
		f.logger.Errorf("create output dir error:%v", e)
		return nil, fmt.Errorf("create output dir failed:%s", e.Error())
	}

	cmd := fmt.Sprintf(
		"./export_frames -i %s -o %s -postfix jpg -c %d -ss %f -s 256x256",
		"TODO", output, p.Count, p.Offset)

	f.logger.Debugf("run cmd:%s", cmd)
	c := exec.Command(cmd)
	e = c.Run()
	if e != nil {
		f.logger.Errorf("run cmd error:%v", e)
		return nil, e
	}

	return nil, nil
}

func (f cmdFrameRandom) Clean(result interface{}) error {
	return nil
}

type cmdFlowRandom struct {
	cmdFrameRandom
}

func init() {
	l := logger.New(os.Stderr, "[cmd] ", logger.Ldefault)
	Register(Cmd, target.Frame, pattern.Random, cmdFrameRandom{logger: l})
	Register(Cmd, target.Flow, pattern.Random, cmdFlowRandom{
		cmdFrameRandom: cmdFrameRandom{logger: l},
	})
}
