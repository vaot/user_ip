package user_ip

import (
  "testing"
  "net/http"
  "github.com/stretchr/testify/assert"
  "github.com/kataras/iris"
  "github.com/valyala/fasthttp"
)

func TestGetNoIp(t *testing.T) {
  req := &http.Request{Method: "GET"}
  ip := Get(req)

  assert.Equal(t, ip, "", "Should not find an ip address.")
}

func TestGetFromContextNoIp(t *testing.T) {
  var ctx iris.Context
  ctx.RequestCtx = &fasthttp.RequestCtx{}
  ip := GetFromContext(&ctx)

  assert.Equal(t, ip, "", "Should not find an ip address.")
}

func TestGetFromContextWithPrivateForwadedIp(t *testing.T) {
  var ctx iris.Context
  ctx.RequestCtx = &fasthttp.RequestCtx{}
  header := fasthttp.RequestHeader{}
  header.Add("X-Forwarded-For", "100.64.0.0, 192.168.0.0")
  ctx.RequestCtx.Request.Header = header
  ip := GetFromContext(&ctx)

  assert.Equal(t, ip, "", "Should not find an ip address.")
}

func TestGetWithPrivateForwadedIp(t *testing.T) {
  req := &http.Request{Method: "GET", Header: make(http.Header)}
  ip := Get(req)

  req.Header.Set("X-Forwarded-For", "100.64.0.0, 192.168.0.0")

  assert.Equal(t, ip, "", "Should not find an ip address.")
}

func TestGetFromContextWithGlobalAndPrivateForwardedIps(t *testing.T) {
  var ctx iris.Context
  ctx.RequestCtx = &fasthttp.RequestCtx{}
  header := fasthttp.RequestHeader{}
  header.Add("X-Forwarded-For", "65.55.37.104, 100.64.0.0, 192.168.0.0, 192.0.0.0")
  ctx.RequestCtx.Request.Header = header
  ip := GetFromContext(&ctx)

  assert.Equal(t, ip, "65.55.37.104", "Should be the first ip, which is the global ip address.")
}

func TestGetWithGlobalAndPrivateForwardedIps(t *testing.T) {
  req := &http.Request{Method: "GET", Header: make(http.Header)}
  req.Header.Set("X-Forwarded-For", "65.55.37.104, 100.64.0.0, 192.168.0.0, 192.0.0.0")

  ip := Get(req)
  assert.Equal(t, ip, "65.55.37.104", "Should be the first ip, which is the global ip address.")
}
