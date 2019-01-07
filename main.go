package main

import (
	"gochat/libs/perf"
	"runtime"
)

var (
	DefaultServer *Server
	DefaultSeq    *Sequence
)

func main() {
	maxproc := runtime.NumCPU()
	runtime.GOMAXPROCS(maxproc)
	
	//pprof
	perf.Init(PPROFBIND);

	// init buckets
	buckets := make([]*Bucket, BUCKETNUM)
	for i := 0; i < BUCKETNUM; i++ {
		buckets[i] = NewBucket(BucketOptions{
			ChannelSize:   CHANNELSIZE,
			RoomSize:      ROOMSIZE,
			RoutineAmount: ROUTINEAMOUNT,
			RoutineSize:   ROUTINESIZE,
		})
	}

	timer := runtime.NumCPU()
	round := NewRound(RoundOptions{
		Timer:     timer,
		TimerSize: TIMERSIZE,
	})

	// new server
	DefaultServer = NewServer(buckets, round, ServerOptions{
		CliProto:         CLIPROTO,
		SvrProto:         SVRPROTO,
		HandshakeTimeout: HANDSHAKETIMEOUT,
	});
	//序列
	DefaultSeq = NewSequence();
	// websocket comet
	if err := InitWebsocket(); err != nil {
		panic(err)
	}

	//sub
	go InitSubscribe();

	// block until a signal is received.
	InitSignal()
}
