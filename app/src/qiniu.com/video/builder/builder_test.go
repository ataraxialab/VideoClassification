package builder

import (
	"os"
	"testing"

	"qiniu.com/video/flow"
	"qiniu.com/video/frame"
	"qiniu.com/video/logger"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"

	"github.com/stretchr/testify/assert"
)

func newFile(t *testing.T, filename, data string) string {
	f, err := os.Create(filename)
	f.WriteString(data)
	f.Close()
	assert.Nil(t, err)
	return filename
}

func TestValid(t *testing.T) {
	assert.NotNil(t, Valid())
	defer func() {
		os.Remove(trainLabelFile)
		os.Remove(valLabelFile)
	}()
	videoRoot = "."
	assert.NotNil(t, Valid())

	trainLabelFile = newFile(t, "trainLabelFile", "abc,123")
	assert.NotNil(t, Valid())

	valLabelFile = newFile(t, "valLabelFile", "ab1c,123")
	assert.Nil(t, Valid())

	cmd := &cmdRandom{
		logger: logger.Std,
	}
	assert.Nil(t, cmd.Init())
	assert.Equal(t, 2, len(cmd.videoName2Label))
}

func TestLoadLabel(t *testing.T) {
	f1 := newFile(t, "invalid_label", "abc,123\n111")
	_, err := loadLabel(f1)
	assert.NotNil(t, err)
	os.Remove(f1)
	f1 = newFile(t, "invalid_label", "abc,123\n111, abc")
	_, err = loadLabel(f1)
	assert.NotNil(t, err)
	_, err = loadLabel(f1)
	assert.NotNil(t, err)
	os.Remove(f1)
}

func TestBuilder(t *testing.T) {
	videoRoot = "."
	trainLabelFile = newFile(t, "trainLabelFile", "program,2\ntrainLabelFile,1\ncmd,123\n\nfactory,111")
	valLabelFile = newFile(t, "valLabelFile", "valLabelFile,2\nbuilder,123\nbulder_test,1")
	program = newFile(t, "./program", "#!/bin/bash\nmkdir -p $4 && touch $4/frame && touch $4/flow")
	os.Chmod(program, os.ModePerm)
	defer func() {
		os.Remove(trainLabelFile)
		os.Remove(valLabelFile)
		os.Remove(program)
		os.RemoveAll(outputRoot)
	}()

	cmd := GetBuilder(Cmd, target.Frame, pattern.Random)
	assert.Nil(t, cmd.Init())
	ret, err := cmd.Build(Params{
		Count:  10,
		Offset: 0.2,
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	_, ok := ret[0].(frame.Frame)
	assert.True(t, ok)

	os.RemoveAll(outputRoot)
	cmd = GetBuilder(Cmd, target.Flow, pattern.Random)
	assert.Nil(t, cmd.Init())
	assert.Nil(t, err)
	ret, err = cmd.Build(Params{
		Count:  100,
		Offset: 0.12,
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	_, ok = ret[0].(flow.Flow)
	assert.True(t, ok)
}

func TestDeleteEmptyDir(t *testing.T) {
	os.Mkdir("abc", os.ModePerm)

	_, e := os.Open("abc")
	assert.Nil(t, e)

	cmd := &cmdRandom{
		logger: logger.Std,
	}
	cmd.deleteEmptyDir("abc")
	_, e = os.Open("abc")
	assert.True(t, os.IsNotExist(e))
}
