package goblin

import (
	//"io"
	//"strings"
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/wxjs33/napi/log"
	"github.com/wxjs33/napi/utils"
	"github.com/wxjs33/napi/hserver"
)

type AddHandler struct {
	gh  *GoblinHandler
	gmc *GoblinMysqlContext
	log *log.Log
}
type DeleteHandler struct {
	gh  *GoblinHandler
	gmc *GoblinMysqlContext
	log *log.Log
}
type ReadHandler struct {
	gmc *GoblinMysqlContext
	log *log.Log
}
type ReadServerHandler struct {
	gmc *GoblinMysqlContext
	log *log.Log
}
type AddServerHandler struct {
	gmc *GoblinMysqlContext
	log *log.Log
}
type DeleteServerHandler struct {
	gmc *GoblinMysqlContext
	log *log.Log
}

func (h *AddHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed")
		http.Error(w, "Read from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed")
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Create record request: ", data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" || data.Action == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}
	if data.Expire <= 0 {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Name invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}

	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	servers, err := h.gmc.QueryGetServer(db, data.Product)
	if err != nil { //TODO: error should be 400 product not found
		hserver.ReturnError(w, err, h.log)
		return
	}

	data.Ruleid, err = utils.NewUUID()
	if err != nil { //TODO: error should be 500 generate ruleid faild
		hserver.ReturnError(w, err, h.log)
		return
	}

	data.Expire = int(time.Now().Unix() + int64(data.Expire))
	err = h.gmc.QueryInsert(db, data.Ip, data.Uid, data.Uuid, data.Product, data.Expire, data.Action, data.Ruleid)
	if err != nil {
		hserver.ReturnError(w, err, h.log)
		return
	}

	errFlag := 0
	for _, server := range *servers {
		err = h.gh.RuleAdd(server, data.Ip, data.Uid, data.Uuid, data.Expire, data.Action)
		if err != nil {
			h.log.Error("Rule add to %s failed, data is %s", server, data)
			errFlag = 1
		}
	}

	if errFlag == 1 {
		hserver.ReturnError(w, err, h.log)
		return
	}

	err = h.gmc.QueryUpdateResult(db, data.Ruleid)
	if err != nil {
		hserver.ReturnError(w, err, h.log)
		return
	}

	hserver.ReturnResponse(w, nil, h.log)
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Delete record request: ", data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Ip, Uid, Uuid invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Product invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}

	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed")
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	servers, err := h.gmc.QueryGetServer(db, data.Product)
	if err != nil { //TODO: error should be 400 product not found
		h.log.Error("Get server failed: ", err)
		hserver.ReturnError(w, err, h.log)
		return
	}
	h.log.Debug("Get server result is: ", servers)

	resp, err := h.gmc.QueryRead(db, data.Ip, data.Uid, data.Uuid, data.Product)
	if err != nil {
		h.log.Error("Read rule failed: ", err)
		hserver.ReturnError(w, err, h.log)
		return
	}
	ruleid := resp.Ruleid
	action := resp.Action

	errFlag := 0
	for _, server := range *servers {
		err = h.gh.RuleDelete(server, data.Ip, data.Uid, data.Uuid, action)
		if err != nil {
			h.log.Error("Rule delete to %s failed, data is %s", server, data)
			errFlag = 1
		}
	}

	if errFlag == 1 {
		hserver.ReturnError(w, err, h.log)
		return
	}

	err = h.gmc.QueryUpdateDeleted(db, ruleid)
	if err != nil {
		h.log.Error("Update deleted failed: ", err)
		hserver.ReturnError(w, err, h.log)
		return
	}

	hserver.ReturnResponse(w, nil, h.log)
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &GoblinRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s ", req.RemoteAddr, data)

	/* Check input */
	if data.Ip == "" && data.Uid == "" && data.Uuid == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Ip, Uid, Uuid invalid", http.StatusBadRequest)
		return
	}
	if data.Product == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Product invalid", http.StatusBadRequest)
		return
	}

	if data.Ip == "" {
		data.Ip = "0.0.0.0"
	}
	if data.Uid == "" {
	   data.Uid = "0"
	}
	if data.Uuid == "" {
		data.Uuid = "0"
	}
	//TODO: Deal wildcard */
	//name := data.Name
	//if strings.Contains(data.Name, "*") {
	//	name = strings.Replace(data.Name, "*", "%", -1)
	//	like = 1
	//}
	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	cresps, err := h.gmc.QueryRead(db, data.Ip, data.Uid, data.Uuid, data.Product)
	if err != nil {
		hserver.ReturnError(w, err, h.log)
		return
	}

	//resps := &GoblinReadResponse{}
	//for _, cresp := range *cresps {
	//	resp := GoblinRequest{}
	//	resp.Ip      = cresp.Ip
	//	resp.Uid     = cresp.Uid
	//	resp.Uuid    = cresp.Uuid
	//	resp.Product = cresp.Product
	//	resp.Expire  = cresp.Expire
	//	resp.Action  = resp.Action
	//	*resps = append(*resps, resp)
	//}

	hserver.ReturnResponse(w, cresps, h.log)
}

func (h *ReadServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s, data is %s", req.RemoteAddr, data)
	if data.Addr == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Addr invalid", http.StatusBadRequest)
		return
	}

	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	resps, err := h.gmc.QueryReadServer(db, data.Addr)
	if err != nil {
		h.log.Error("Query read server failed: %s", err)
		hserver.ReturnError(w, err, h.log)
		return
	}

	hserver.ReturnResponse(w, resps, h.log)
}

func (h *AddServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed:%s ", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s ", req.RemoteAddr, data)
	if data.Addr == "" || data.Product == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Addr or product invalid", http.StatusBadRequest)
		return
	}

	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	err = h.gmc.QueryAddServer(db, data.Addr, data.Product)
	if err != nil {
		h.log.Error("Query add server failed: %s", err)
		hserver.ReturnError(w, err, h.log)
		return
	}

	hserver.ReturnResponse(w, nil, h.log)
}

func (h *DeleteServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		h.log.Error("Method invalid: %s", req.Method)
		http.Error(w, "Method invalid", http.StatusBadRequest)
		return
	}

	result, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error("Read from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	req.Body.Close()

	data := &ServerRequest{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		h.log.Error("Parse from request body failed: %s", err)
		http.Error(w, "Parse from body failed", http.StatusBadRequest)
		return
	}
	h.log.Info("Read record request from %s ", req.RemoteAddr, data)
	if data.Addr == "" || data.Product == "" {
		h.log.Error("Post arguments invalid")
		http.Error(w, "Addr or product invalid", http.StatusBadRequest)
		return
	}

	db, err := h.gmc.Open()
	if err != nil {
		h.log.Error("Mysql open failed: %s", err)
		http.Error(w, "Mysql open failed", http.StatusBadGateway)
		return
	}
	defer h.gmc.Close(db)

	err = h.gmc.QueryDeleteServer(db, data.Addr, data.Product)
	if err != nil {
		h.log.Error("Query delete server failed: %s", err)
		hserver.ReturnError(w, err, h.log)
		return
	}

	hserver.ReturnResponse(w, nil, h.log)
}
