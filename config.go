//config

package main

import (
	"time"
)

var (
	PPROFBIND = []string{"localhost:2045"}
	WEBSOCKETBIND    = []string{"0.0.0.0:8099"}
)

const (
	Debug            = true;

	WEBSOCKETTLSOPEN = false
	
	HANDSHAKETIMEOUT = 5 * time.Second
	WRITETIMEOUT     = 5 * time.Second
	HEARTBEAT        = 60 * time.Second
	
	TIMERSIZE        = 1000
	
	CLIPROTO         = 5
	SVRPROTO         = 80
	
	// bucket
	BUCKETNUM        = 512	
	CHANNELSIZE      = 1024	
	ROOMSIZE         = 10
	
	//goroutine
	ROUTINEAMOUNT    = 128
	ROUTINESIZE      = 20
	
	AUTHKEY          = "e2a7edKLD31XkmgoMZBBS"

	REDISHOST        = "127.0.0.1:6379";
	REDISDB          = 0;
	REDISPWD         = "";

)