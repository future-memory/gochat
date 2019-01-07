package main

import (
	"gochat/libs/hash/cityhash"
	"time"
)

var (
	maxInt        = 1<<31 - 1
	emptyJSONBody = []byte("{}")
)

type ServerOptions struct {
	CliProto         int
	SvrProto         int
	HandshakeTimeout time.Duration
}

type Server struct {
	Buckets   []*Bucket // subkey bucket
	round     *Round // accept round store
	bucketIdx uint32
	Options   ServerOptions
}

// NewServer returns a new Server.
func NewServer(b []*Bucket, r *Round, options ServerOptions) *Server {
	s := new(Server)
	s.Buckets = b
	s.round = r
	s.bucketIdx = uint32(len(b))
	s.Options = options
	return s
}

func (server *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % server.bucketIdx;
	return server.Buckets[idx]
}
