local log = require("log")
local json = require("json")
local http = require("http")
local mysql = require("mysql")

function ReceiveTGMsg(data)
    log.info("%s", "\nReceiveTGMsg")

    str =
        string.format(
        "ChatID %d\nMessageID %d\nSenderUserID %d\nMsgType %s\nContent %s",
        data.ChatID,
        data.MessageID,
        data.SenderUserID,
        data.MsgType,
        data.Content
        -- PkgCodec.EncodeBase64(data.ImgBuf)
    )
    log.notice("From log.lua Log\n%s", str)
    return 1
end

function ReceiveEvents(data)
    return 1
end
