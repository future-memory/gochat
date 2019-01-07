package main

import (
	"github.com/tidwall/gjson"
    "crypto/md5" 
    "encoding/hex"	
    "log"
	"time"
	"strconv"
)

// developer could implement "Auth" interface for decide how get userId, or roomId
type Auther interface {
	Auth(token string) (userId string, roomId int32)
}

type DefaultAuther struct {
}

func NewDefaultAuther() *DefaultAuther {
	return &DefaultAuther{}
}

func (a *DefaultAuther) Auth(token string) (userId string, roomId int32) {
	userId = gjson.Get(token, "uid").String();
	code  := gjson.Get(token, "code").String();
	time1 := gjson.Get(token, "time").String();
	tmpid := gjson.Get(token, "roomid").Int();

	time2, _ := strconv.ParseInt(time1, 10, 64);
	roomId = int32(tmpid);

	//验证时间
	t  := time.Now();
	tt := t.Unix();

	if Debug {
		log.Printf("auth info,uid:%s,code:%s,time1:%s,time2:%d,tt:%d,token:%v,tmpid:%v", userId, code, time1, time2, tt, token, tmpid);
	}

	if(tt-time2>60){
		userId = guest();
		return		
	}

	mc := Md5(AUTHKEY+userId+time1);
	if Debug {
		log.Printf("auth info,mc:%s", mc);
	}

	if(mc!=code){
		userId = guest();
	}
	return
}

func guest() (userId string){
	guestid := DefaultSeq.nextSeq();
	userId = "g_"+strconv.Itoa(int(guestid));
	return
}

func Md5(str string) string {
    h := md5.New()
    h.Write([]byte(str))
    data := h.Sum([]byte(""))
    return hex.EncodeToString(data)
}