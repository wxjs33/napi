package goblin

import (
	"fmt"
	"bytes"
	//"time"
	"net/http"
	"github.com/wxjs33/napi/variable"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/errors"
)

type GoblinHandler struct {
	add_loc    string
	delete_loc string
	log        *log.Log
}

func InitGoblinHandler(loc string, log *log.Log) *GoblinHandler {
	gh := &GoblinHandler{}
	gh.log = log
	gh.add_loc = loc + GOBLIN_ADD_LOCATION
	gh.delete_loc = loc + GOBLIN_DELETE_LOCATION

	return gh
}

func (gh *GoblinHandler) RuleOperate(addr string, args *bytes.Buffer, op int) error {
	var err error
	var resp *http.Response

	switch op {
	case ADD_RULE:
		gh.log.Debug("add rule args is ", args)
		resp, err = http.Post("http://" + addr + gh.add_loc, variable.DEFAULT_CONTENT_HEADER, args)
		//resp, err = http.Get("http://" + addr + gh.add_loc + "?" + args)
		break
	case DELETE_RULE:
		gh.log.Debug("delete rule args is ", args)
		resp, err = http.Post("http://" + addr + gh.delete_loc, variable.DEFAULT_CONTENT_HEADER, args)
		//resp, err = http.Get("http://" + addr + gh.delete_loc + "?" + args)
		break
	default: /* Should not reach here */
		gh.log.Error("Unknown operate code: ", op)
		return errors.InternalServerError
	}

	if err != nil {
		gh.log.Error("Opereate service to nginx failed: ", err)
		return errors.BadGatewayError
	}

	defer resp.Body.Close()

	if resp.StatusCode != variable.HTTP_OK {
		gh.log.Error("Opereate http status error: %d", resp.StatusCode)
		return errors.BadGatewayError
	}

	return nil
}

func (gh *GoblinHandler) RuleAdd(addr, ip, uid, uuid string, expire int, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	//if uid != "" {
	//	args = fmt.Sprint(args, "&uid=", uid)
	//}
	if uuid != "" {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}
	args = fmt.Sprint(args, "&expire=", expire, "&punish=", action, "&punish_arg=0\r\n")
	gh.log.Error(args)

	/* post data */
	data := bytes.NewBufferString(args)

	return gh.RuleOperate(addr, data, ADD_RULE)
	//return gh.RuleOperate(addr, args, ADD_RULE)
}

func (gh *GoblinHandler) RuleDelete(addr, ip, uid, uuid, action string) error {
	if ip == "" {
		ip = EMPTY_IP
	}
	args := fmt.Sprint("startip=", ip, "&endip=", ip)
	//if uid != "" {
	//	args = fmt.Sprint(args, "&uid=", uid)
	//}
	if uuid != "" {
		args = fmt.Sprint(args, "&uuid=", uuid)
	}

	args = fmt.Sprint(args, "&punish=", action)

	gh.log.Error(args)
	data := bytes.NewBufferString(args)
	return gh.RuleOperate(addr, data, DELETE_RULE)
	//return gh.RuleOperate(addr, args, DELETE_RULE)
}
