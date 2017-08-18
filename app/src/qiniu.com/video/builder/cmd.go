package builder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
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
	Cmd Implement = "cmd"
)

var (
	program = "./export_frames"
)

type cmdRandom struct {
	name            string
	logger          *logger.Logger
	buildLabel      func(string, int) (interface{}, error)
	videos          []string
	videoName2Label map[string]int
}

func createOutputDir(prefix string) (string, error) {
	d := path.Join(outputRoot, prefix, time.Now().Format("20060102-150405.000"))
	return d, os.MkdirAll(d, os.ModePerm)
}

func (f *cmdRandom) selectVideo() string {
	return f.videos[rand.Intn(len(f.videos))]
}

// interpret label name from video path
func labelName(video string) string {
	basename := path.Base(video)
	dotIndex := strings.IndexByte(basename, '.')
	if dotIndex < 0 {
		return basename
	}

	return basename[:dotIndex]
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

	video := f.selectVideo()
	labelName := labelName(video)
	label, exist := f.videoName2Label[labelName]
	if !exist {
		return nil, fmt.Errorf("not exist of label:%s, video:%s", labelName, video)
	}

	arg := fmt.Sprintf("-i %s -o %s -postfix jpg -c %d -ss %f -s 256x256",
		video, output, p.Count, p.Offset)

	f.logger.Errorf("run cmd:%s %s", program, arg)
	cmd := exec.Command(program, strings.Split(arg, " ")...)
	out, e := cmd.CombinedOutput()
	f.logger.Infof("output:%s, error:%v", out, e)

	if e != nil {
		f.logger.Errorf("run cmd error:%v", e)
		os.RemoveAll(output)
		return nil, e
	}

	return f.buildLabels(f.getFiles(output), label), nil
}

func (f *cmdRandom) buildLabels(files []string, label int) []interface{} {
	ret := make([]interface{}, 0, len(files))
	for _, file := range files {
		l, err := f.buildLabel(file, label)
		if err != nil {
			f.logger.Errorf("build label error:%v", err)
			continue
		}
		ret = append(ret, l)
	}
	return ret
}

func buildFrame(file string, label int) (interface{}, error) {
	return frame.Frame{
		Label:     float32(label),
		ImagePath: file,
	}, nil
}

func buildFlow(file string, label int) (interface{}, error) {
	return flow.Flow{
		Label:     float32(label),
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

func (f *cmdRandom) Init() error {
	if err := Valid(); err != nil {
		return err
	}

	if f.videos == nil {
		videos := make([]string, 0, 1024)
		filepath.Walk(videoRoot,
			func(p string, info os.FileInfo, err error) error {
				if err != nil {
					f.logger.Errorf("error :%v", err)
					return nil // IGNORE
				}
				if videoRoot == p {
					return nil
				}
				videos = append(videos, p)
				return nil
			})
		f.videos = videos
	}

	videoCount := len(f.videos)
	if len(f.videos) == 0 {
		return fmt.Errorf("no videos under:%s", videoRoot)
	}

	f.logger.Infof("video count:%d", videoCount)

	m, err := loadLabels(trainLabelFile, valLabelFile)
	if err != nil {
		f.logger.Errorf("load labels error:%v", err)
		return fmt.Errorf("load labels error:%s", err.Error())
	}

	f.logger.Infof("label count:%d", len(m))
	f.videoName2Label = m
	return nil
}

func loadLabel(labelFile string) (map[string]int, error) {
	content, err := ioutil.ReadFile(labelFile)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(content, []byte("\n"))
	m := make(map[string]int, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		elems := bytes.Split(line, []byte(","))
		if len(elems) < 2 {
			return nil, fmt.Errorf("bad label data:%s", line)
		}
		m[string(elems[0])], err = strconv.Atoi(string(elems[1]))
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func loadLabels(labelfiles ...string) (labels map[string]int, err error) {
	for _, f := range labelfiles {
		m, err := loadLabel(f)
		if err != nil {
			return nil, err
		}

		if labels == nil {
			labels = m
			continue
		}

		for k, v := range m {
			labels[k] = v
		}
	}

	return labels, nil
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
