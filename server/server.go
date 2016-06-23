package server

import (
	//"time"
	//"strings"
	//"github.com/wxjs33/napi/variable"
	"github.com/wxjs33/napi/hserver"
	"github.com/wxjs33/napi/config"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/goblin"

	//"github/wxjs33/napi/"
)

type Server struct {
	addr    string

	hs      *hserver.HttpServer

	//gc      *GoblinContext
	//mc      *MysqlContext
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
	//hs.s = s

	s.log.Debug("Init http server done")
	//s.log.Debug(conf.c)

	//_, err = goblin.InitGoblinMysqlContext(conf, hs, log)
	//if err != nil {
	//	s.log.Error("Init mysql client faild")
	//	return nil, err
	//}
	//hs.AddRouter(mc.url, mc)

	_, err = goblin.InitGoblinContext(conf, hs, log)
	if err != nil {
		s.log.Error("Init goblin context failed")
		return nil, err
	}
	//hs.AddRouter(gc.url, gc)

	return s, nil
}

//func (s *Server) HttpServer() (*HttpServer) {
//	return s.hs
//}
//
//func (s *Server) MysqlContext() (*MysqlContext) {
//	return s.mc
//}

func (s *Server) Run() error {
	err := s.hs.Run()
	if err != nil {
		s.log.Error("Server run failed: ", err)
		return err
	}

	return nil
}

