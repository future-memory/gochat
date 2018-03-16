package main

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"gochat/libs/proto"
	"gochat/libs/define"
)

type Proto struct {
	Operation int32  // operation for request
	Body      []byte // body
}

//从redis订阅
func InitSubscribe(){
	c, err := redis.Dial("tcp", "172.16.180.232:6379")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe("live")

	for {
		switch v := psc.Receive().(type) {
			case redis.Message:
				message := []byte(v.Data)
				var bucket *Bucket;
				var p *proto.Proto;
				p.Operation = define.OP_SEND_SMS_REPLY;
				p.Body = message;

				for _, bucket = range DefaultServer.Buckets {
					go bucket.Broadcast(p)
				}
			case redis.Subscription:
			case error:
				log.Printf("redis error: %v", v)
			return
		}
	}
 }