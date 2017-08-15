package web

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"qiniu.com/video/builder"
)

type mockServer int

func (s *mockServer) StartBuilding(target builder.Target, pattern builder.Pattern, params interface{}) error {
	*s++
	return nil
}
func (s *mockServer) StopBuilding(target builder.Target, pattern builder.Pattern) error {
	*s--
	return nil
}
func (s *mockServer) GetResult(target builder.Target, pattern builder.Pattern, from, to uint) (interface{}, error) {
	return nil, nil
}

func (s *mockServer) Close() error {
	return nil
}

func readAll(resp *http.Response) string {
	body := resp.Body
	data, _ := ioutil.ReadAll(body)
	body.Close()
	return string(data)
}

func TestSwitch(t *testing.T) {
	srv := mockServer(0)
	s := newHTTPServer(context.TODO(), &srv)

	ts := httptest.NewServer(s)
	defer ts.Close()

	call := func(url string, body string) (*http.Response, error) {
		return http.Post(url, "applicaton/json", bytes.NewReader([]byte(body)))
	}

	resp, err := call(ts.URL+"/abc", "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "target")

	resp, err = call(ts.URL+"/frame", "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "json")

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"abc",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "pattern")

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"++start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "op")

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"start",
		"params":{
			"interval":0,
			"duration":0
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "duration")

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, int(srv))

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"stop"
		}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, int(srv))

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"sample",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "interval")

	resp, err = call(ts.URL+"/frame", `{
		"pattern":"sample",
		"op":"start",
		"params":{
			"interval":1,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, int(srv))
}
