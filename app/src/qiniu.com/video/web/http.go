package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/server"

	"github.com/julienschmidt/httprouter"
)

type httpError struct {
	Code int    `json:"-"`
	Text string `json:"message"`
}

type decorator func(apiHandler) apiHandler
type apiHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, *httpError)

func doDecorate(f apiHandler, ds ...decorator) httprouter.Handle {
	decorated := f
	for _, d := range ds {
		decorated = d(decorated)
	}

	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		decorated(w, req, ps)
	}
}

type httpServer struct {
	ctx    context.Context
	router http.Handler
	logger *logger.Logger
	server server.Server
}

func newHTTPServer(ctx context.Context, server server.Server) *httpServer {
	router := httprouter.New()

	router.PanicHandler = panicHandler
	s := &httpServer{
		ctx:    ctx,
		router: router,
		logger: logger.Std,
		server: server,
	}

	s.logger.Level = logger.Ldebug
	s.logger.SetPrefix("[http] ")

	router.POST("/:target", doDecorate(s.switchOp, s.jsonDecorator))
	return s
}

func panicHandler(resp http.ResponseWriter, req *http.Request, p interface{}) {
	logger.Errorf("panic HTTP handler :%v", p)
	resp.WriteHeader(http.StatusInternalServerError)
}

type switchParam struct {
	Pattern string         `json:"pattern"`
	Op      string         `json:"op"`
	Params  builder.Params `json:"params"`
}

func (s *httpServer) switchOp(w http.ResponseWriter,
	req *http.Request,
	ps httprouter.Params,
) (interface{}, *httpError) {
	pTarget := ps.ByName("target")
	target := builder.GetTarget(pTarget)
	if !target.IsValid() {
		s.logger.Errorf("unknow target:%s", pTarget)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: fmt.Sprintf("unknow target:%s", pTarget),
		}
	}

	params := switchParam{}
	err := parseJSONParam(req.Body, &params)
	if err != nil {
		s.logger.Errorf("parse request body error:%v", err)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: "request body is not json",
		}
	}

	pattern := builder.GetPattern(params.Pattern)
	if !pattern.IsValid() {
		s.logger.Errorf("unknow pattern:%s", params.Pattern)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: "unknown pattern:" + params.Pattern,
		}
	}

	start, stop := "start", "stop"
	if params.Op != start && params.Op != stop {
		s.logger.Errorf("unknow op:%s", params.Op)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: "unknown op:" + params.Op,
		}
	}

	if params.Op == stop {
		if err = s.server.StopBuilding(target, pattern); err != nil {
			s.logger.Errorf("stop build error:%v", err)
			return nil, &httpError{
				Code: http.StatusForbidden,
				Text: err.Error(),
			}
		}

		return nil, nil
	}

	if params.Params.Duration <= 0 {
		text := fmt.Sprintf("invalid duration:%d", params.Params.Duration)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	if pattern.NeedInterval() && params.Params.Interval <= 0 {
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: fmt.Sprintf("invalid interval:%d", params.Params.Interval),
		}
	}

	if err = s.server.StartBuilding(target, pattern, params.Params); err != nil {
		s.logger.Errorf("start build error:%v", err)
		return nil, &httpError{
			Code: http.StatusForbidden,
			Text: err.Error(),
		}
	}

	return nil, nil
}

func (s *httpServer) jsonDecorator(handler apiHandler) apiHandler {
	return func(w http.ResponseWriter,
		req *http.Request,
		ps httprouter.Params,
	) (interface{}, *httpError) {
		ret, httpErr := handler(w, req, ps)
		if httpErr != nil {
			w.WriteHeader(httpErr.Code)
			bytes, _ := json.Marshal(httpErr)
			w.Write(bytes)
			return nil, nil
		}

		if ret == nil {
			return nil, nil
		}

		bytes, err := json.Marshal(ret)
		if err != nil {
			s.logger.Errorf("marhal data error:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			bytes, _ = json.Marshal(&httpError{
				Text: "INTERNAL_ERROR",
			})
		}
		w.Write(bytes)
		return nil, nil
	}
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func parseJSONParam(body io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
