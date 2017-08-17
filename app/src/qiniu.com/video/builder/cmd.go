package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"qiniu.com/video/flow"
	"qiniu.com/video/frame"
	"qiniu.com/video/logger"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

const (
	// Cmd call external command to build data
	Cmd     Implement = "cmd"
	program string    = "./export_frames"
)

type cmdRandom struct {
	name       string
	logger     *logger.Logger
	buildLabel func(string) (interface{}, error)
}

func createOutputDir(prefix string) (string, error) {
	d := path.Join(outputRoot, prefix, time.Now().Format("20060102-150405."))
	return d, os.MkdirAll(d, os.ModePerm)
}

func selectVideo() string {
	// TODO
	return "test.mp4"
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

	arg := fmt.Sprintf("-interval 10 -i %s -o %s -postfix jpg -c %d -ss %f -s 256x256",
		selectVideo(), output, p.Count, p.Offset)

	f.logger.Errorf("run cmd:%s %s", program, arg)
	cmd := exec.Command(program, strings.Split(arg, " ")...)
	out, e := cmd.CombinedOutput()
	f.logger.Infof("output:%s, error:%v", out, e)

	if e != nil {
		f.logger.Errorf("run cmd error:%v", e)
		os.RemoveAll(output)
		return nil, e
	}

	return f.buildLabels(f.getFiles(output)), nil
}

func (f *cmdRandom) buildLabels(files []string) []interface{} {
	ret := make([]interface{}, 0, len(files))
	for _, file := range files {
		l, err := f.buildLabel(file)
		if err != nil {
			f.logger.Errorf("build label error:%v", err)
			continue
		}
		ret = append(ret, l)
	}
	return ret
}

func buildFrame(file string) (interface{}, error) {
	// FIXME
	return frame.Frame{
		Label:     float32(999),
		ImagePath: file,
	}, nil
}

func buildFlow(file string) (interface{}, error) {
	// FIXME
	return flow.Flow{
		Label:     float32(999),
		ImagePath: file,
	}, nil
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
	if r, ok := result.(frame.Frame); ok {
		return os.Remove(r.ImagePath)
	}

	if r, ok := result.(flow.Flow); ok {
		return os.Remove(r.ImagePath)
	}
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
	Register(Cmd, target.Frame, pattern.Random, &cmdRandom{logger: l,
		name:       "frame",
		buildLabel: buildFrame,
	})
	Register(Cmd, target.Flow, pattern.Random, &cmdRandom{logger: l,
		name:       "flow",
		buildLabel: buildFlow,
	})
}
