//config

package main

import (
	"time"
)

const (	
	WEBSOCKETTLSOPEN = false
	
	HANDSHAKETIMEOUT = 5 * time.Second
	WRITETIMEOUT     = 5 * time.Second
	
	TIMERSIZE        = 1000
	
	CLIPROTO         = 5
	SVRPROTO         = 80
	
	// bucket
	BUCKETNUM        = 1024	
	BUCKETCHANNEL    = 1024	
	BUCKETROOM       = 10
	
	//goroutine
	ROUTINEAMOUNT    = 128
	ROUTINESIZE      = 20
	
	//tcp
	TCPSNDBUF        = 2048
	TCPRCVBUF        = 256
	TCPKEEPALIVE     = false
	TCPREADER        = 1024
	TCPREADBUF       = 1024
	TCPREADBUFSIZE   = 512
	TCPWRITER        = 1024
	TCPWRITEBUF      = 1024
	TCPWRITEBUFSIZE  = 4096

	//proto
	CliProto  = 5
	SvrProto  = 80

	SERVERID  = 1

	Ver = "0.1"

	AUTHKEY = "test"
)

var (
	TCPBIND       = []string{"0.0.0.0:8080"}
	WEBSOCKETBIND = []string{"0.0.0.0:8099"}
	PPROFBIND     = []string{"localhost:2044"}
	STATBIND      = []string{"localhost:2045"}

	Debug         = false;
)