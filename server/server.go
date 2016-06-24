package server

import (
	"github.com/wxjs33/napi/hserver"
	"github.com/wxjs33/napi/config"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/modules"
)

type Server struct {
	addr    string

	hs      *hserver.HttpServer

	domain  string

	log     *log.Log
}

func InitServer(conf *config.Config, log *log.Log) (*Server, error) {
	s := &Server{}

	s.log = log
	s.addr = conf.Addr

	hs, err := hserver.InitHttpServer(conf.Addr, s.log)
	if err != nil {
		s.log.Error("Init http server failed")
		return nil, err
	}
	s.hs = hs

	s.log.Debug("Init http server done")

	modules.InitModules(conf, hs, log)

	return s, nil
}

func (s *Server) Run() error {
	err := s.hs.Run()
	if err != nil {
		s.log.Error("Server run failed: ", err)
		return err
	}

	return nil
}

