package main

import (
	"net/url"
	"strconv"

	tcp "NetworkFramework"

	"tdlib"

	lua "github.com/yuin/gopher-lua"
	//luar "layeh.com/gopher-luar"
)

type LuaCoreApiModule struct {
}

func NewCoreApiModule() *LuaCoreApiModule {
	return &LuaCoreApiModule{}
}
func (l *LuaCoreApiModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"SendText":       l.Api_SendText,
		"SendPhoto":      l.Api_SendPhoto,
		"AddProxy":       l.Api_AddProxy,
		"EnableProxy":    l.Api_EnableProxy,
		"RemoveProxy":    l.Api_RemoveProxy,
		"DeleteMessages": l.Api_DeleteMessages,
	})
	L.Push(mod)
	return 1
}
func (l *LuaCoreApiModule) Api_SendText(L *lua.LState) int {
	ChatID := L.CheckInt64(1)
	text := L.CheckString(2)
	if text != "" {
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, false)
		client.SendMessage(ChatID, 0, false, true, nil, inputMsgTxt)
	}
	return 0
}
func (l *LuaCoreApiModule) Api_SendPhoto(L *lua.LState) int {
	ChatID := L.CheckInt64(1)
	path := L.CheckString(2)
	text := L.CheckString(3)
	if path != "" {
		inputMsg := tdlib.NewInputMessagePhoto(tdlib.NewInputFileLocal(path), nil, nil, 400, 400,
			tdlib.NewFormattedText(text, nil), 0)
		client.SendMessage(ChatID, 0, false, false, nil, inputMsg)

	}
	return 0
}
func (l *LuaCoreApiModule) Api_AddProxy(L *lua.LState) int {
	//ChatID := L.CheckInt64(1)
	link := L.CheckString(1)
	isUser := L.CheckBool(2)
	//https://t.me/proxyserver=13.76.43.181&port=443&secret=ee0000e78a4fa0072b8b30d21cd02de9ad617a7572652e6d6963726f736f66742e636f6d
	if link != "" {
		val, err := url.ParseQuery(link)

		if err == nil {
			port, _ := strconv.Atoi(val.Get("port"))
			_, err := client.AddProxy(val.Get("server"), int32(port), isUser, tdlib.NewProxyTypeMtproto(val.Get("secret")))
			tcp.Logger.Trace(" addproxy %v", err)
		}
	}
	return 0
}
func (l *LuaCoreApiModule) Api_EnableProxy(L *lua.LState) int {
	ID := L.CheckInt64(1)
	client.EnableProxy(int32(ID))
	return 0
}
func (l *LuaCoreApiModule) Api_RemoveProxy(L *lua.LState) int {
	ID := L.CheckInt64(1)

	client.RemoveProxy(int32(ID))
	return 0
}
func (l *LuaCoreApiModule) Api_DeleteMessages(L *lua.LState) int {
	ChatID := L.CheckInt64(1)
	MessageID := L.CheckInt64(2)
	revoke := L.CheckBool(3)
	client.DeleteMessages(ChatID, []int64{MessageID}, revoke)
	return 0
}
