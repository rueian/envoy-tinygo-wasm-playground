package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"strings"
	"sync"
)

type rootContext struct {
	proxywasm.DefaultRootContext
}

type httpContext struct {
	proxywasm.DefaultHttpContext
}

type replacer struct {
	*strings.Replacer
	mu sync.RWMutex
}

var rp = replacer{}

func newRootContext(contextID uint32) proxywasm.RootContext {
	return &rootContext{}
}

func newHttpContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	return &httpContext{}
}

func (ctx *rootContext) OnVMStart(vmConfigurationSize int) bool {
	if err := proxywasm.SetSharedData("dict", []byte("headers|請求標頭|Cookie|美味餅乾"), 0); err != nil {
		proxywasm.LogWarnf("fail to set dict to share date: %v", err)
	}
	ctx.OnTick()
	proxywasm.SetTickPeriodMilliSeconds(1000)
	return true
}

func (ctx *rootContext) OnTick() {
	proxywasm.LogInfo("OnTick, updating dictionary")

	if data, _, err := proxywasm.GetSharedData("dict"); err == nil {
		dict := strings.Split(string(data), "|")
		rp.mu.Lock()
		rp.Replacer = strings.NewReplacer(dict...)
		rp.mu.Unlock()
	} else {
		proxywasm.LogWarnf("fail to set dict to share date: %v", err)
	}
}

func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if bodySize == 0 {
		return types.ActionContinue
	}
	resp, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("fail to get resp body: %s", err)
		return types.ActionContinue
	}
	cp := string(resp)
	rp.mu.RLock()
	ret := rp.Replace(cp)
	rp.mu.RUnlock()
	proxywasm.SetHttpResponseBody([]byte(ret))

	return types.ActionContinue
}

func main() {
	proxywasm.SetNewRootContext(newRootContext)
	proxywasm.SetNewHttpContext(newHttpContext)
}
