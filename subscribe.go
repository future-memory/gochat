package main

import (
	"log"
	"gochat/libs/proto"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

type redisSubData struct{
	Id int32
	Op int32
	Uid int32
	Body json.RawMessage 
}

//从redis订阅
func InitSubscribe(){
	c, err := redis.Dial("tcp", REDISHOST, redis.DialPassword(REDISPWD))
	if err != nil {
		log.Println(err)
		panic("redis connect fail");
		return
	}
	c.Do("SELECT", REDISDB)  

	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe("live")

	for {
		switch v := psc.Receive().(type) {
			case redis.Message:
				message := []byte(v.Data);
				if(Debug){
					//log.Printf("redis msg: %v %s", message, message);
				}

				var data redisSubData;
				json.Unmarshal(message, &data);
				
				if(Debug){
					//log.Printf("redis msg decode: %v, id: %d", data, data.Id);
				}

				var bucket *Bucket;

				//推送指定用户
				if data.Uid>0 {
					var ch *Channel;
					var p *proto.Proto;

					p = ch.CliProto.LoopSet();
					p.Operation = data.Op;
					p.Body = data.Body;
					ch.CliProto.SetAdv();

					key := encode(string(data.Uid), data.Id);
					bucket = DefaultServer.Bucket(key)
					if ch = bucket.Channel(key); ch != nil {
						if err = ch.Push(p); err != nil {
							log.Printf("redis error: %v", err)
						}
					}

				}else{
					
					for _, bucket = range DefaultServer.Buckets {
						if data.Id==1 { 
							//推送所有
							go bucket.BroadcastRedisMsg(data.Op, data.Body);
						}else{
							//推送房间
							go bucket.BroadcastRoomRedisMsg(data.Id, data.Op, data.Body);
						}
					}
				}

			case redis.Subscription:
				log.Printf("redis sub")
			case error:
				log.Printf("redis error: %v", v)
			return
		}
	}
 }