package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//取文件所在路径
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func GetLuaList(dirPth string, suffix string) (luascripts map[string]string, err error) {
	scripts := make(map[string]string)
	//	log.Printf("path %s", GetCurrentDirectory()+dirPth)

	dir, err := ioutil.ReadDir(GetAppPath() + dirPth)
	if err != nil {

		return scripts, err
	}
	//	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			userBuff := InFile(dirPth + fi.Name())
			//log.Println(dirPth + fi.Name())
			//log.Println(fi.Name())
			scripts[fi.Name()] = string(userBuff)

		}
	}
	return scripts, nil
}
func InFile(name string) []byte {

	if contents, err := ioutil.ReadFile(GetAppPath() + name); err == nil {

		return contents
	}

	return nil
}

func HttpRequest(method string, urls string, data []byte, jar http.CookieJar, head map[string]string) ([]byte, error, *http.Response) {

	c := &http.Client{

		Jar: jar,
		//Timeout: time.Second * 5,
	}

	// create a socks5 dialer
	// dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1086", nil, proxy.Direct)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	// 	//os.Exit(1)
	// }
	// // setup a http client
	// httpTransport := &http.Transport{}
	// httpTransport.Dial = dialer.Dial
	// c.Transport = httpTransport

	reqest, err := http.NewRequest(method, urls, bytes.NewBuffer(data))

	for k, v := range head { //遍历key和value。
		//fmt.Println("key=", k, "value=", v)

		//增加header选项

		reqest.Header.Add(k, v)

	}

	if err != nil {
		return nil, err, nil
	}
	//处理返回结果
	resp, err := c.Do(reqest)

	if err != nil {

		return nil, err, nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err, resp
	}
	resp.Body.Close()
	return body, nil, resp
}
