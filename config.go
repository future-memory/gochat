//config

package main

import (
	"runtime"
)

const (
	WEBSOCKETBIND    = "0.0.0.0:8099"
	PPROFBIND        = "localhost:2044"
	STATBIND         = "localhost:2045"
	
	WEBSOCKETTLSOPEN = false
	
	HANDSHAKETIMEOUT = 5 * time.Second
	WRITETIMEOUT     = 5 * time.Second
	
	TIMER            = runtime.NumCPU()
	MaxProc          = runtime.NumCPU()
	TIMERSIZE        = 1000
	
	CLIPROTO         = 5
	SVRPROTO         = 80
	
	// bucket
	BUCKETNUM        = 1024	
	BUCKETCHANNEL    = 1024	
	BUCKETROOM       = 1024
	
	//goroutine
	ROUTINEAMOUNT    = 128
	ROUTINESIZE      = 20
	
	//tcp
	TCPBIND          = "0.0.0.0:8080"
	TCPSNDBUF        = 2048
	TCPRCVBUF        = 256
	TCPKEEPALIVE     = 0
	TCPREADER        = 1024
	TCPREADBUF       = 1024
	TCPREADBUFSIZE   = 512
	TCPWRITER        = 1024
	TCPWRITEBUF      = 1024
	TCPWRITEBUFSIZE  = 4096

	//proto
	CliProto  = 5
	SvrProto  = 80

)
