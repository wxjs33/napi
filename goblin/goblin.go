package goblin

import (
	//"fmt"
	//"time"
	//"io/ioutil"
	//"net/http"
	//"encoding/json"
	//"bytes"
	//"github.com/wxjs33/napi/variable"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/hserver"
	"github.com/wxjs33/napi/config"
)

type GoblinContext struct {
	api_loc   string
	loc       string

	log       *log.Log
	gh        *GoblinHandler
	gmc       *GoblinMysqlContext
	cf         *GoblinConfig
}

func InitGoblinContext(conf *config.Config, hs *hserver.HttpServer, log *log.Log) (*GoblinContext, error) {
	gc := &GoblinContext{}

	gc.cf = &GoblinConfig{}
	err := gc.cf.ParseConfig(conf)
	if err != nil {
		log.Error("goblin parse config failed")
		return nil, err
	}

	gc.gh = InitGoblinHandler(gc.cf.goblin_loc, log)
	gc.gmc = InitGoblinMysqlContext(gc.cf.maddr, gc.cf.dbname, gc.cf.dbuser, gc.cf.dbpwd, log)
	gc.log = log


	//TODO: parse conf
	gc.api_loc = gc.cf.api_loc
	hs.AddRouter(gc.api_loc + RULE_ADD_LOCATION, &AddHandler{gh: gc.gh, gmc: gc.gmc, log: log})
	hs.AddRouter(gc.api_loc + RULE_DELETE_LOCATION, &DeleteHandler{gh: gc.gh, gmc: gc.gmc, log: log})
	hs.AddRouter(gc.api_loc + RULE_READ_LOCATION, &ReadHandler{gmc: gc.gmc, log: log})
	hs.AddRouter(gc.api_loc + SERVER_ADD_LOCATION, &AddServerHandler{gmc: gc.gmc, log: log})
	hs.AddRouter(gc.api_loc + SERVER_DELETE_LOCATION, &DeleteServerHandler{gmc: gc.gmc, log: log})
	hs.AddRouter(gc.api_loc + SERVER_READ_LOCATION, &ReadServerHandler{gmc: gc.gmc, log: log})

	return gc, nil
}
