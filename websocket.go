package main

import (
	"gochat/libs/define"
	"gochat/libs/proto"
	itime "gochat/libs/time"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/thinkboy/log4go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebsocket(addrs []string) (err error) {
	var (
		bind         string
		listener     *net.TCPListener
		addr         *net.TCPAddr
		httpServeMux = http.NewServeMux()
		server       *http.Server
	)

	log.Debug("init ws")

	httpServeMux.HandleFunc("/live", ServeWebSocket)

	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp4", bind); err != nil {
			log.Error("net.ResolveTCPAddr(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp4", addr); err != nil {
			log.Error("net.ListenTCP(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		server = &http.Server{Handler: httpServeMux}
		if Debug {
			log.Debug("start websocket listen: \"%s\"", bind)
		}
		go func(host string) {
			if err = server.Serve(listener); err != nil {
				log.Error("server.Serve(\"%s\") error(%v)", host, err)
				panic(err)
			}
		}(bind)
	}
	return
}

func ServeWebSocket(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Websocket Upgrade error(%v), userAgent(%s)", err, req.UserAgent())
		return
	}
	defer ws.Close()
	var (
		lAddr = ws.LocalAddr()
		rAddr = ws.RemoteAddr()
		tr    = DefaultServer.round.Timer(rand.Int())
	)
	log.Debug("start websocket serve \"%s\" with \"%s\"", lAddr, rAddr)
	DefaultServer.serveWebsocket(ws, tr)
}

func (server *Server) serveWebsocket(conn *websocket.Conn, tr *itime.Timer) {
	var (
		err error
		key string
		hb  time.Duration // heartbeat
		p   *proto.Proto
		b   *Bucket
		trd *itime.TimerData
		ch  = NewChannel(server.Options.CliProto, server.Options.SvrProto, define.NoRoom)
	)
	// handshake
	log.Debug("kk %v", server.Options.HandshakeTimeout)
	trd = tr.Add(server.Options.HandshakeTimeout, func() {
		log.Debug("c  %v", server.Options.HandshakeTimeout)
		conn.Close()
	})
	// must not setadv, only used in auth
	if p, err = ch.CliProto.Set(); err == nil {
		log.Debug("ff")
		if key, ch.RoomId, hb, err = server.authWebsocket(conn, p); err == nil {
			log.Debug("ee %v", hb)
			b = server.Bucket(key)
			err = b.Put(key, ch)
		}
	}
	if err != nil {
		conn.Close()
		tr.Del(trd)
		log.Error("handshake failed error(%v)", err)
		return
	}
	log.Debug("k %v", key);
	trd.Key = key
	tr.Set(trd, hb)
	// hanshake ok start dispatch goroutine
	go server.dispatchWebsocket(key, conn, ch)

	//读取上报的信息，并广播或pong
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			log.Debug("a")
			break
		}
		if err = p.ReadWebsocket(conn); err != nil {
			log.Debug("b %v", err)
			break
		}
		//p.Time = *globalNowTime
		if p.Operation == define.OP_HEARTBEAT {
			// heartbeat
			tr.Set(trd, hb)
			p.Body = nil
			p.Operation = define.OP_HEARTBEAT_REPLY
		}
		//信息 广播
		if p.Operation == define.OP_SEND_SMS {
			var bucket *Bucket
			p.Operation = define.OP_SEND_SMS_REPLY;
			for _, bucket = range DefaultServer.Buckets {
				go bucket.Broadcast(p)
			}
		}

		ch.CliProto.SetAdv()
		ch.Signal()
	}
	log.Error("key: %s server websocket failed error(%v)", key, err)
	tr.Del(trd)
	conn.Close()
	ch.Close()
	b.Del(key)
	if Debug {
		log.Debug("key: %s server websocket goroutine exit", key)
	}
	return
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (server *Server) dispatchWebsocket(key string, conn *websocket.Conn, ch *Channel) {
	var (
		p   *proto.Proto
		err error
	)
	if Debug {
		log.Debug("key: %s start dispatch websocket goroutine", key)
	}
	for {
		p = ch.Ready()
		switch p {
		case proto.ProtoFinish:
			if Debug {
				log.Debug("key: %s wakeup exit dispatch goroutine", key)
			}
			goto failed
		case proto.ProtoReady:
			for {
				if p, err = ch.CliProto.Get(); err != nil {
					err = nil // must be empty error
					break
				}
				if err = p.WriteWebsocket(conn); err != nil {
					goto failed
				}
				p.Body = nil // avoid memory leak
				ch.CliProto.GetAdv()
			}
		default:
			// TODO room-push support
			// just forward the message
			if err = p.WriteWebsocket(conn); err != nil {
				goto failed
			}
		}
	}
failed:
	if err != nil {
		log.Error("key: %s dispatch websocket error(%v)", key, err)
	}
	conn.Close()
	// must ensure all channel message discard, for reader won't blocking Signal
	for {
		if p == proto.ProtoFinish {
			break
		}
		p = ch.Ready()
	}
	if Debug {
		log.Debug("key: %s dispatch goroutine exit", key)
	}
	return
}

func (server *Server) authWebsocket(conn *websocket.Conn, p *proto.Proto) (key string, rid int32, heartbeat time.Duration, err error) {
	if err = p.ReadWebsocket(conn); err != nil {
		log.Debug("a e %v", err)
		return
	}
	if p.Operation != define.OP_AUTH {
		log.Debug("a a %v %v", p.Operation, define.OP_AUTH);
		err = ErrOperation
		return
	}

	auther := NewDefaultAuther();

	log.Debug("pbody:%v", string(p.Body));
	var uid int64;
	uid, rid = auther.Auth(string(p.Body))
	var seq int32 = 1;
	key = encode(uid, seq)
	heartbeat = 5 * 60 * time.Second

	p.Body = emptyJSONBody
	p.Operation = define.OP_AUTH_REPLY
	err = p.WriteWebsocket(conn)
	return
}
