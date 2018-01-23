package main

import (
	"gochat/libs/perf"
	"runtime"
)

var (
	DefaultServer    *Server
)

func main() {
	runtime.GOMAXPROCS(MAXPROC)
	perf.Init(PPROFBIND)
	// new server
	buckets := make([]*Bucket, BUCKETNUM)
	for i := 0; i < BUCKETNUM; i++ {
		buckets[i] = NewBucket(BucketOptions{
			ChannelSize:   BUCKETCHANNEL,
			RoomSize:      BUCKETROOM,
			RoutineAmount: ROUTINEAMOUNT,
			RoutineSize:   ROUTINESIZE,
		})
	}

	round := NewRound(RoundOptions{
		Reader:       TCPREADER,
		ReadBuf:      TCPREADBUF,
		ReadBufSize:  TCPREADBUFSIZE,
		Writer:       TCPWRITER,
		WriteBuf:     TCPWRITEBUF,
		WriteBufSize: TCPWRITEBUFSIZE,
		Timer:        TIMER,
		TimerSize:    TIMERSIZE,
	})

	operator := new(DefaultOperator)
	DefaultServer = NewServer(buckets, round, operator, ServerOptions{
		CliProto:         CLIPROTO,
		SvrProto:         SVRPROTO,
		HandshakeTimeout: HANDSHAKETIMEOUT,
		TCPKeepalive:     TCPKEEPALIVE,
		TCPRcvbuf:        TCPRCVBUF,
		TCPSndbuf:        TCPSNDBUF,
	})

	// websocket comet
	if err := InitWebsocket(WEBSOCKETBIND); err != nil {
		panic(err)
	}

	// block until a signal is received.
	InitSignal()
}
