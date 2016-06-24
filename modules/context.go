package modules
import (
	"github.com/wxjs33/napi/config"
	"github.com/wxjs33/napi/hserver"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/modules/goblin"
	"github.com/wxjs33/napi/modules/sample"
)

func InitModules(conf *config.Config, hs *hserver.HttpServer, log *log.Log) {
	if err := goblin.InitContext(conf, hs, log); err != nil {
		log.Error("goblin module will not start")
	}

	if err := sample.InitContext(conf, hs, log); err != nil {
		log.Error("sample module will not start")
	}
}
