package builder

import (
	"fmt"
	"os"
)

// Builder build the data from vido
type Builder interface {
	Build(params interface{}) ([]interface{}, error)
	Clean(interface{}) error
	Init() error
}

// Implement the build implementation
type Implement string

// Params building parameters, it is BAD design
type Params struct {
	Count  int     `json:"count"`
	Offset float32 `json:"offset"`
}

var outputRoot = "build-output"
var videoRoot = ""
var trainLabelFile = ""
var valLabelFile = ""

// SetOutputRoot set build output directory
func SetOutputRoot(d string) {
	outputRoot = d
}

// SetVideoRoot set the video directory
func SetVideoRoot(d string) {
	videoRoot = d
}

// SetTrainLabelFile set the file containing the labels for trainning
func SetTrainLabelFile(f string) {
	trainLabelFile = f
}

// SetValLabelFile set the file containing the labels for validing
func SetValLabelFile(f string) {
	valLabelFile = f
}

func exists(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}

// Valid the params
func Valid() error {
	if !exists(videoRoot) {
		return fmt.Errorf("not exists of video root:%s", videoRoot)
	}

	if !exists(trainLabelFile) {
		return fmt.Errorf("not exists of trainning label file:%s",
			trainLabelFile)
	}

	if !exists(valLabelFile) {
		return fmt.Errorf("not exists of validing label file:%s", valLabelFile)
	}

	return nil
}
