package main

import (
	tcp "NetworkFramework"
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"tdlib"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	client  *tdlib.Client
	poolOne WorkPool
)

func main() {

	logs.SetLogger(logs.AdapterConsole, `{"level":8}`)

	if runtime.GOOS != "windows" {

		logs.SetLogger(logs.AdapterFile, `{"filename":"`+GetAppPath()+`/Logs/tg.log","level":3}`)

	} else {

		logs.SetLogger(logs.AdapterFile, `{"filename":"./Logs/tg.log","level":3}`)
	}
	//初始化go协程池
	poolOne.InitPool(60)
	//GO层实现日志接口 如果包不存在 替换所有日志输出的地方即可编译
	tcp.Logger = &TGLog{}

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./Logs/errors.txt")

	// Create new instance of client
	client = tdlib.NewClient(tdlib.Config{
		APIID:               "793416",
		APIHash:             "021de84fe4f1ac0361c333b0ba6198b6",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  false,
		UseFileDatabase:     false,
		UseChatInfoDatabase: false,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})
	//client.AddProxy("127.0.0.1", 443, true, tdlib.NewProxyTypeMtproto("ee4012999756fd6eb3fafd63fd17cb3c70617a7572652e6d6963726f736f66742e636f6d"))

	// Handle Ctrl+C
	// ch := make(chan os.Signal, 2)
	// signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-ch
	// 	client.DestroyInstance()
	// 	os.Exit(1)
	// }()

	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			go GetProxy()
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v\n", err)
			}

			// p, err := client.CheckAuthenticationBotToken(":AAFomEPDiMQ6hE4dpmDFkKpHrmawsvwA")
			// fmt.Println(p, err)
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v\n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v\n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization Ready! Let's rock")
			break
		}
	}
	go GetMSG()
	go CheckProxy()
	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {

		// Show all updates
		tcp.Logger.Info("raw %s", update.Raw)
		//fmt.Print("\n\n")
	}
}
func GetProxy() {
	// 会一直维护的
	body, _, _ := HttpRequest("GET", "http://129.204.103.68:58897/v1/TGProxys", nil, nil, nil)

	if body != nil {

		var _jArray []interface{}
		json.Unmarshal(body, &_jArray)

		tcp.Logger.Critical("%v", _jArray)

		for _, v := range _jArray {
			link := v.(string)
			link = strings.ReplaceAll(link, "?", "&")
			val, err := url.ParseQuery(link)
			if err == nil {
				port, _ := strconv.Atoi(val.Get("port"))
				_, err := client.AddProxy(val.Get("server"), int32(port), true, tdlib.NewProxyTypeMtproto(val.Get("secret")))
				tcp.Logger.Trace(" AddProxy %v", err)
			}

		}
	}

}
func CheckProxy() {

	heartbeat := time.NewTicker(30 * time.Second)
	pullproxy := time.NewTicker(1 * time.Hour) // 1个小时拉取一次代理列表
	proxyFlag := false

	pmap := make(map[int32]int)

	for {
		select {
		case <-heartbeat.C:
			tcp.Logger.Alert("正在检测 活跃代理")
			proxys, err := client.GetProxies()
			if err == nil {
				for _, v := range proxys.Proxies {
					t, err := client.PingProxy(v.ID)
					coust := 0.0
					ChackFlag := false
					if err != nil {
						ChackFlag = true
					} else {
						if t.Seconds == 0 {
							ChackFlag = true
						}
					}
					if ChackFlag {
						count := 0
						if _, ok := pmap[v.ID]; ok {
							pmap[v.ID] += 1
							count = pmap[v.ID]
						} else {
							pmap[v.ID] = 1
						}
						tcp.Logger.Error("Proxy %d Err %v Try %d", v.ID, err, count)
						if count == 15 {
							if v.IsEnabled {
								proxyFlag = true
							}
							delete(pmap, v.ID)
							client.RemoveProxy(v.ID)

						}
						continue
					} else {
						if _, ok := pmap[v.ID]; ok {
							pmap[v.ID] = 0
						}
						coust = t.Seconds
					}
					l, err := client.GetProxyLink(v.ID)
					if err == nil {
						if proxyFlag {
							client.EnableProxy(v.ID)
							proxyFlag = false
						}
						tcp.Logger.Emergency("Check Ok Proxy %d Ping %fs Link %s IsEnabled %v", v.ID, coust, l.Text, v.IsEnabled)
					}

				}

			}
			break
		case <-pullproxy.C:
			GetProxy()
			break

		}
	}

}
func GetMSG() {
	// Create an filter function which will be used to filter out unwanted tdlib messages
	eventFilter := func(msg *tdlib.TdMessage) bool {
		// updateMsg := (*msg).(*tdlib.UpdateNewMessage)

		// // For example, we want incomming messages from user with below id:
		// if updateMsg.Message.SenderUserID == 41507975 {
		// 	return true
		// }
		return true
	}

	// Here we can add a receiver to retreive any message type we want
	// We like to get UpdateNewMessage events and with a specific FilterFunc
	receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		//fmt.Println(newMsg)
		updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
		// We assume the message content is simple text: (should be more sophisticated for general use)

		if msgText, ok := updateMsg.Message.Content.(*tdlib.MessageText); ok {
			m := make(map[string]interface{})
			m["ChatID"] = updateMsg.Message.ChatID
			m["SenderUserID"] = updateMsg.Message.SenderUserID
			m["MessageID"] = updateMsg.Message.ID
			m["MsgType"] = "MessageText"
			m["Content"] = msgText.Text.Text
			//协程池执行lua插件
			poolOne.Run(TGLuaVMRun, m)

		}

	}

}
