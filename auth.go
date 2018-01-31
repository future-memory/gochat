package main

import (
	"gochat/libs/define"
	"encoding/json"
    "crypto/md5" 
    "encoding/hex"	
	"time"
	log "github.com/thinkboy/log4go"
)

// developer could implement "Auth" interface for decide how get userId, or roomId
type Auther interface {
	Auth(token string) (userId string, roomId int32)
}

type AuthCode struct {
	userId string        `json:"uid"`
	code   string        `json:"code"`
	time   int64         `json:"time"`
	roomId int32         `json:"id"`
}

type DefaultAuther struct {
}

func NewDefaultAuther() *DefaultAuther {
	return &DefaultAuther{}
}

func (a *DefaultAuther) Auth(token string) (userId string, roomId int32) {
	var err error
	tmp := AuthCode{}
	log.Debug("token: %v", token)

	if err = json.Unmarshal([]byte(token), &tmp); err != nil {
		userId = "0"
		roomId = define.NoRoom
		log.Debug("decode err: %v", err)
		return
	}

	userId = tmp.userId;
	roomId = tmp.roomId;

	log.Debug("tmp: %v", tmp)

	//验证时间
	t  := time.Now();
	tt := t.Unix();
	if(tt-tmp.time>60){
		userId = "0"
		roomId = define.NoRoom
		log.Debug("time err: %v, %v", tt, tmp.time)
		return		
	}

	mc := Md5(AUTHKEY+userId+string(tmp.time));
	if(mc!=tmp.code){
		userId = "0"
		roomId = define.NoRoom
		log.Debug("md5 err: %v, %v", mc, tmp.code)		
	}
	return
}

func Md5(str string) string {  
    h := md5.New()  
    h.Write([]byte(str))  
    data := h.Sum([]byte(""))  
    return hex.EncodeToString(data)  
} 