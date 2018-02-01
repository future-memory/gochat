package main

import (
	"gochat/libs/perf"
	"runtime"
)

var (
	DefaultServer    *Server
	DefaultSeq       *Sequence
)

func main() {
	maxproc := runtime.NumCPU()
	runtime.GOMAXPROCS(maxproc)
	perf.Init([]string(PPROFBIND))
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

	timer := runtime.NumCPU()
	round := NewRound(RoundOptions{
		Reader:       TCPREADER,
		ReadBuf:      TCPREADBUF,
		ReadBufSize:  TCPREADBUFSIZE,
		Writer:       TCPWRITER,
		WriteBuf:     TCPWRITEBUF,
		WriteBufSize: TCPWRITEBUFSIZE,
		Timer:        timer,
		TimerSize:    TIMERSIZE,
	})

	DefaultServer = NewServer(buckets, round, ServerOptions{
		CliProto:         CLIPROTO,
		SvrProto:         SVRPROTO,
		HandshakeTimeout: HANDSHAKETIMEOUT,
		TCPKeepalive:     TCPKEEPALIVE,
		TCPRcvbuf:        TCPRCVBUF,
		TCPSndbuf:        TCPSNDBUF,
	})
	//序列
	DefaultSeq = NewSequence();
	// websocket comet
	if err := InitWebsocket([]string(WEBSOCKETBIND)); err != nil {
		panic(err)
	}

	// block until a signal is received.
	InitSignal()
}
