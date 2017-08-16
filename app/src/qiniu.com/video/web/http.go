package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/server"
	"qiniu.com/video/target"

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

func newHTTPHandler(ctx context.Context, server server.Server) *httpServer {
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
	router.GET("/:target/:pattern/:from/:count",
		doDecorate(s.getBuildResult, s.jsonDecorator))
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
	target := target.GetTarget(pTarget)
	if !target.IsValid() {
		text := fmt.Sprintf("unknow target:%s", pTarget)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	params := switchParam{}
	if err := parseJSONParam(req.Body, &params); err != nil {
		s.logger.Errorf("parse request body error:%v", err)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: "request body is not json",
		}
	}

	pattern := pattern.GetPattern(params.Pattern)
	if !pattern.IsValid() {
		text := "unknown pattern:" + params.Pattern
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	start, stop := "start", "stop"
	if params.Op != start && params.Op != stop {
		text := "unknown op:" + params.Op
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	if params.Op == stop {
		if err := s.server.StopBuilding(target, pattern); err != nil {
			s.logger.Errorf("stop build error:%v", err)
			return nil, &httpError{
				Code: http.StatusForbidden,
				Text: err.Error(),
			}
		}

		return nil, nil
	}

	if params.Params.Count <= 0 {
		text := fmt.Sprintf("invalid count:%d", params.Params.Count)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	if params.Params.Offset <= 0 {
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: fmt.Sprintf("invalid offset:%f", params.Params.Offset),
		}
	}

	if err := s.server.StartBuilding(target, pattern, params.Params); err != nil {
		s.logger.Errorf("start build error:%v", err)
		return nil, &httpError{
			Code: http.StatusForbidden,
			Text: err.Error(),
		}
	}

	return nil, nil
}

func (s *httpServer) getBuildResult(w http.ResponseWriter,
	req *http.Request,
	ps httprouter.Params,
) (interface{}, *httpError) {
	pTarget := ps.ByName("target")
	target := target.GetTarget(pTarget)
	if !target.IsValid() {
		text := fmt.Sprintf("unknow target:%s", pTarget)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	p := ps.ByName("pattern")
	pattern := pattern.GetPattern(p)
	if !pattern.IsValid() {
		text := fmt.Sprintf("unknow pattern:%s", p)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	p = ps.ByName("from")
	from, err := strconv.Atoi(p)
	if err != nil || from < 0 {
		text := fmt.Sprintf("bad from:%s", p)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	p = ps.ByName("count")
	count, err := strconv.Atoi(p)
	if err != nil || count <= 0 {
		text := fmt.Sprintf("bad count:%s", p)
		s.logger.Errorf(text)
		return nil, &httpError{
			Code: http.StatusBadRequest,
			Text: text,
		}
	}

	result, err := s.server.GetResult(target, pattern, uint(from),
		uint(math.Min(1000, count)))
	if err != nil {
		s.logger.Errorf("get result error:%v", err)
		return nil, &httpError{
			Code: http.StatusForbidden,
			Text: err.Error(),
		}
	}

	return result, nil
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

// HTTPServer http server
type HTTPServer struct {
	httpListener net.Listener
	httpServer   *httpServer
	mq           mq.MQ
}

// NewHTTPServer create http server
func NewHTTPServer(ctx context.Context, port int) (*HTTPServer, error) {
	mq := &mq.EmbeddedMQ{}
	if err := mq.Open(); err != nil {
		return nil, err
	}

	server, err := server.CreateServer(builder.Cmd, mq)
	if err != nil {
		mq.Close()
		return nil, err
	}

	httpListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	srv := &HTTPServer{
		httpListener: httpListener,
		httpServer:   newHTTPHandler(ctx, server),
		mq:           mq,
	}

	return srv, nil
}

// Serve http service
func (s *HTTPServer) Serve() {
	go func() {
		logger := s.httpServer.logger
		logger.Infof("start http server on:%s", s.httpListener.Addr().String())

		server := http.Server{
			Handler: s.httpServer,
		}

		err := server.Serve(s.httpListener)
		if err != nil {
			logger.Errorf("server error:%v", err)
		}
	}()
}

// Close close http server
func (s *HTTPServer) Close() {
	logger := s.httpServer.logger

	err := s.httpListener.Close()
	if err != nil {
		logger.Errorf("close http server error:%v", err)
	}

	err = s.httpServer.server.Close()
	if err != nil {
		logger.Errorf("close server error:%v", err)
	}

	err = s.mq.Close()
	if err != nil {
		logger.Errorf("close mq error:%v", err)
	}
	logger.Infof("http server closed")
}
