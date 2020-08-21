# TGBOT

**GO+LUA+TDLIB 跨平台TG机器人框架**


GitHub上开源的TG一些库,bot-api用的不习惯 `pyrogram` `telethon`比如不支持新的代理加密方式MTProto2.0 ,一些 dd ee 开头的 secret 再此就显得很鸡肋了 不如官方的TDlib香 还可以跨平台 公益代理还极其不稳定, 半身不遂的TG机器人🤖️ 属实是。。。&&……%！#¥@……！

> 采用Glang开发  也是第一次尝试使用Go链接c++ 在配上Lua 写一些小插件就很完美了 多个平台运行不用cao太多❤️ 

> 项目自动维护了TG的公益代理 定时检测并对失效的代理自动剔除和切换

>Lua 只绑定了几个功能 就当是抛砖引玉吧

> 协程池处理机制

## 环境配置

* 将 TGCLI和tdlib 两个目录放到你的GoPath目录下 从CrossLib中选出你编译平台的库文件放到TGCLI中
* 编译启用 CGO_ENABLED=1
* 并不支持一个平台编译出 多个平台  工具链配置可能比较繁琐 建议在使用环境下编译或使用我编译好的自行开发插件即可
* Linux 下编译需要在TGCLI目录下创建一个软链 `ln -s libtdjson.so.1.6.0 libtdjson.so`
* Windows 需要 `MinGW-W64` 环境 将 TGCLI 目录下的td 目录复制到tdlib中 在将编译好的tdjson.dll放入 tdlib/td目录下进行链接编译即可 或自己配置 头文件引用目录也可以
## 开箱即用

* 只封装了发消息 发图片 撤回消息几个接口  对于TG爬🐛 消息监听 足够用   有需要可以自行二开

[🔗下载地址🔗](https://github.com/OPQBOT/TGBOT/releases)

```javascript
function test() {
	console.log("欢迎大佬 拍砖 吐槽");
}
```