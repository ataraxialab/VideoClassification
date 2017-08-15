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
	return []int{1, 2, 3}, nil
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

func post(url string, body string) (*http.Response, error) {
	return http.Post(url, "applicaton/json", bytes.NewReader([]byte(body)))
}

func get(url string) (*http.Response, error) {
	return http.Get(url)
}

func TestSwitch(t *testing.T) {
	srv := mockServer(0)
	s := newHTTPServer(context.TODO(), &srv)

	ts := httptest.NewServer(s)
	defer ts.Close()

	resp, err := post(ts.URL+"/abc", "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "target")

	resp, err = post(ts.URL+"/frame", "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "json")

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"abc",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "pattern")

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"++start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "op")

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"start",
		"params":{
			"interval":0,
			"duration":0
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "duration")

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, int(srv))

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"random",
		"op":"stop"
		}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, int(srv))

	resp, err = post(ts.URL+"/frame", `{
		"pattern":"sample",
		"op":"start",
		"params":{
			"interval":0,
			"duration":1
		}}`)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "interval")

	resp, err = post(ts.URL+"/frame", `{
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

func TestGetResult(t *testing.T) {
	srv := mockServer(0)
	s := newHTTPServer(context.TODO(), &srv)

	ts := httptest.NewServer(s)
	defer ts.Close()

	resp, err := get(ts.URL + "/t/p/0/1")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "target")

	resp, err = get(ts.URL + "/frame/p/-1/0")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "pattern")

	resp, err = get(ts.URL + "/frame/random/-1/0")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "from")

	resp, err = get(ts.URL + "/frame/random/0/a")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "count")

	resp, err = get(ts.URL + "/frame/random/0/-1")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readAll(resp), "count")

	resp, err = get(ts.URL + "/frame/random/0/1")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "[1,2,3]", readAll(resp))
}
