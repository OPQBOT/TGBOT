local log = require("log")
local json = require("json")
local http = require("http")
local mysql = require("mysql")
local Api = require("coreApi")
function ReceiveTGMsg(data)
	--发文字
    if (string.find(data.Content, "复读机") == 1) then
        keyWord = data.Content:gsub("复读机", "")
        Api.SendText(data.ChatID, keyWord)
        return 1
    end
    --发图
    if (string.find(data.Content, "ph") == 1) then
        Api.SendPhoto(data.ChatID, "./1.png", "hello")
        return 1
    end
    --添加代理
    if string.find(data.Content, "secret") then
        link = data.Content:gsub("?", "&")
        log.info("%s", link)
        Api.AddProxy(link, true)
        return 1
    end
    --撤回消息
    if (string.find(data.Content, "de") == 1) then
        Api.DeleteMessages(data.ChatID, data.MessageID, true)
        return 1
    end

    return 1
end

function ReceiveEvents(data)
    return 1
end
