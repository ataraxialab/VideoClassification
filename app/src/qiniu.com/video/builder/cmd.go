package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"qiniu.com/video/logger"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

const (
	// Cmd call external command to build data
	Cmd        Implement = "cmd"
	outputRoot string    = "output"
	program    string    = "./export_frames"
)

type cmdRandom struct {
	name   string
	logger *logger.Logger
}

func createOutputDir(prefix string) (string, error) {
	d := path.Join(outputRoot, prefix, time.Now().Format("20060102-150405"))
	return d, os.MkdirAll(d, os.ModePerm)
}

func selectVideo() string {
	// TODO
	return "TODO"
}

func (f *cmdRandom) Build(params interface{}) ([]interface{}, error) {
	p, ok := params.(Params)
	if !ok {
		f.logger.Errorf("unknown params:%v", p)
		return nil, fmt.Errorf("unknown params")
	}

	output, e := createOutputDir(f.name)
	if e != nil {
		f.logger.Errorf("create output dir error:%v", e)
		return nil, fmt.Errorf("create output dir failed:%s", e.Error())
	}

	arg := fmt.Sprintf("-i %s -o %s -postfix jpg -c %d -ss %f -s 256x256",
		selectVideo(), output, p.Count, p.Offset)

	f.logger.Debugf("run cmd:%s %s", program, arg)
	cmd := exec.Command(program, arg)
	e = cmd.Run()
	if e != nil {
		f.logger.Errorf("run cmd error:%v", e)
		return nil, e
	}

	return f.buildLabels(f.getFiles(output)), nil
}

func (f *cmdRandom) buildLabels(files []string) []interface{} {
	ret := make([]interface{}, 0, len(files))
	for _, file := range files {
		l, err := buildLabel(file)
		if err != nil {
			f.logger.Errorf("build label error:%v", err)
			continue
		}
		ret = append(ret, l)
	}
	return ret
}

func buildLabel(file string) (interface{}, error) {
	// FIXME
	return nil, nil
}

func (f *cmdRandom) getFiles(d string) []string {
	files := make([]string, 0, 10)
	filepath.Walk(d, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			f.logger.Errorf("walk %s error:%v", p, err)
		}
		if strings.Contains(info.Name(), f.name) {
			files = append(files, p)
		}

		return nil
	})
	return files
}

func (f *cmdRandom) Clean(result interface{}) error {
	return nil
}

type cmdFrameRandom struct {
	*cmdRandom
}

type cmdFlowRandom struct {
	*cmdRandom
}

func init() {
	l := logger.New(os.Stderr, "[cmd] ", logger.Ldefault)
	Register(Cmd, target.Frame, pattern.Random, &cmdRandom{logger: l, name: "frame"})
	Register(Cmd, target.Flow, pattern.Random, &cmdRandom{logger: l, name: "flow"})
}
