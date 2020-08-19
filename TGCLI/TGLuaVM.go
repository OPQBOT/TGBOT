package main

import (
	tcp "NetworkFramework"

	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"time"

	"github.com/cjoudrey/gluahttp"
	mysql "github.com/tengattack/gluasql/mysql"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func TGLuaVMRun(parms []interface{}) {
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	//lua.Options{SkipOpenLibs: true}
	lua.OpenPackage(L) // Must be first
	lua.OpenBase(L)
	lua.OpenTable(L)
	lua.OpenIo(L)
	lua.OpenOs(L)
	lua.OpenString(L)
	lua.OpenMath(L)
	lua.OpenDebug(L)
	lua.OpenChannel(L)
	lua.OpenCoroutine(L)
	defer func() {
		L.Close()
		runtime.GC()
		if err := recover(); err != nil {
			stacks := string(tcp.PanicTrace(4))

			tcp.Logger.Error("TGLuaVMRun will exit panics: %v call:%v", err, stacks)

		}

	}()
	L.PreloadModule("mysql", mysql.Loader)
	L.PreloadModule("log", NewLogModule().Loader)
	L.PreloadModule("coreApi", NewCoreApiModule().Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("json", NewJsonModule().Loader)
	//L.PreloadModule("PkgCodec", NewPkgCodecModule().Loader)
	//	NewPkgCodecModule().Loader(L)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	L.SetContext(ctx)
	defer cancel()
	scripts, _ := GetLuaList("/Plugins/", ".lua")

	for fName, scLua := range scripts {
		var err error
		if err = L.DoString(scLua); err != nil {
			tcp.Logger.Error(fName)
			tcp.Logger.Error(scLua)
			panic(err)
			return
		}

		var Ret int

		var r lua.LValue

		if len(parms) == 1 {
			switch parms[0].(type) {
			case map[string]interface{}:
				fName = fmt.Sprintf("File %s when  call ReceiveFriendMsg params [0] %v", fName, parms[0])
				r, err = CallGlobalValue(L, "ReceiveTGMsg", parms[0])
				break

			}
		}

		if err != nil {
			tcp.Logger.Error("TGLuaVMRun CallGlobal err %v detail %v", err, fName)
			return

		}

		if n, ok := r.(lua.LNumber); ok {

			Ret = int(n)

		}

		if Ret == 1 || err != nil {
			continue
		}
		if Ret == 2 {
			return
		}

	}

}

func ToGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key))
				ret[keystr] = ToGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	case *lua.LUserData:
		return v.Value
	default:

		return v
	}
}
func ToGoBytes(in []interface{}) []byte {
	b := bytes.NewBuffer([]byte{})

	for _, v := range in {

		b.WriteByte(byte(v.(float64)))
	}
	return b.Bytes()

}

func CallGlobalValue(L *lua.LState, fnName string, args ...interface{}) (r lua.LValue, err error) {

	fn := L.GetGlobal(fnName)
	if fn.Type() != lua.LTFunction {
		err = errors.New(fmt.Sprintf("Unknow Lua Function:%v", fnName))
		return
	}

	// 组合参数列表
	lpValues := []lua.LValue{}
	argsArr := []interface{}(args)
	for _, v := range argsArr {
		lpValues = append(lpValues, luar.New(L, v))
	}

	err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, lpValues...)

	r = L.Get(-1)

	// if n, ok := r.(lua.LNumber); ok {
	// 	ret = int(n)
	// 	return
	// }

	L.Pop(1)
	return
}
