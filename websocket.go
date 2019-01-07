package main

import (
	"gochat/libs/define"
	"gochat/libs/proto"
	itime "gochat/libs/time"
	"net"
	"net/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWebsocket() (err error) {
	var (
		bind         string
		listener     *net.TCPListener
		addr         *net.TCPAddr
		httpServeMux = http.NewServeMux()
		server       *http.Server
	)
	httpServeMux.HandleFunc("/live", ServeWebSocket)

	for _, bind = range WEBSOCKETBIND {
		if addr, err = net.ResolveTCPAddr("tcp4", bind); err != nil {
			log.Printf("net.ResolveTCPAddr(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp4", addr); err != nil {
			log.Printf("net.ListenTCP(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		server = &http.Server{Handler: httpServeMux}
		if Debug {
			log.Printf("start websocket listen: \"%s\"", bind)
		}

		go func(host string) {
			if err = server.Serve(listener); err != nil {
				log.Printf("server.Serve(\"%s\") error(%v)", host, err)
				panic(err)
			}
		}(bind)
	}
	return;
}

func ServeWebSocket(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("Websocket Upgrade error(%v), userAgent(%s)", err, req.UserAgent())
		return
	}
	defer ws.Close()
	var (
		lAddr = ws.LocalAddr()
		rAddr = ws.RemoteAddr()
		tr    = DefaultServer.round.Timer(rand.Int())
	)
	if Debug {
		log.Printf("start websocket serve \"%s\" with \"%s\"", lAddr, rAddr)
	}
	DefaultServer.serveWebsocket(ws, tr);
}

func (server *Server) serveWebsocket(conn *websocket.Conn, tr *itime.Timer) {
	var (
		err error
		key string
		p   *proto.Proto
		b   *Bucket
		trd *itime.TimerData
		ch  = NewChannel(server.Options.CliProto, server.Options.SvrProto, define.NoRoom)
	)

	// handshake
	trd = tr.Add(server.Options.HandshakeTimeout, func() {
		conn.Close()
	});

	// must not setadv, only used in auth
	if p, err = ch.CliProto.Set(); err == nil {
		if key, ch.RoomId, err = server.authWebsocket(conn, p); err == nil {
			b = server.Bucket(key)
			err = b.Put(key, ch)
		}
	}
	if err != nil {
		conn.Close()
		tr.Del(trd)
		log.Printf("handshake failed error(%v)", err)
		return
	}
	if Debug {
		log.Printf("key %v", key);
	}

	trd.Key = key
	tr.Set(trd, HEARTBEAT)

	// hanshake ok start dispatch goroutine
	go server.dispatchWebsocket(key, conn, ch, tr, trd)

	//读取上报的信息，并广播或pong
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			log.Printf("proto err %v", err)
			break
		}
		if err = p.ReadWebsocket(conn); err != nil {
			log.Printf("read ws err %v", err)
			break
		}
		//p.Time = *globalNowTime
		if p.Operation == define.OP_HEARTBEAT {
			// heartbeat
			tr.Set(trd, HEARTBEAT)
			p.Body = nil
			p.Operation = define.OP_HEARTBEAT_REPLY
		}
		//信息 广播
		// if p.Operation == define.OP_SEND_SMS {
		// 	var bucket *Bucket
		// 	p.Operation = define.OP_SEND_SMS_REPLY;
		// 	for _, bucket = range DefaultServer.Buckets {
		// 		go bucket.Broadcast(p)
		// 	}
		// }
		ch.CliProto.SetAdv()
		ch.Signal()
	}

	log.Printf("key: %s server websocket failed error(%v)", key, err)
	conn.Close()
	ch.Close()
	tr.Del(trd)
	b.Del(key)

	if Debug {
		log.Printf("key: %s server websocket goroutine exit", key)
	}
	return
}

// dispatch accepts connections on the listener and serves requests
// for each incoming connection.  dispatch blocks; the caller typically
// invokes it in a go statement.
func (server *Server) dispatchWebsocket(key string, conn *websocket.Conn, ch *Channel, tr *itime.Timer, trd *itime.TimerData) {
	var (
		p   *proto.Proto
		err error
	)
	if Debug {
		log.Printf("key: %s start dispatch websocket goroutine", key)
	}
	for {
		p = ch.Ready()
		switch p {
		case proto.ProtoFinish:
			if Debug {
				log.Printf("key: %s wakeup exit dispatch goroutine", key)
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
				ch.CliProto.GetAdv();

				tr.Set(trd, HEARTBEAT);
			}
		default:
			// TODO room-push support
			// just forward the message
			if err = p.WriteWebsocket(conn); err != nil {
				goto failed
			}
			tr.Set(trd, HEARTBEAT);
		}
	}
failed:
	if err != nil {
		log.Printf("key: %s dispatch websocket error(%v)", key, err)
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
		log.Printf("key: %s dispatch goroutine exit", key)
	}
	return
}

func (server *Server) authWebsocket(conn *websocket.Conn, p *proto.Proto) (key string, rid int32, err error) {
	if err = p.ReadWebsocket(conn); err != nil {
		log.Printf("read ws err %v", err)
		return
	}
	if p.Operation != define.OP_AUTH {
		log.Printf("auth op err %v %v", p.Operation, define.OP_AUTH);
		err = ErrOperation
		return
	}

	auther := NewDefaultAuther();

	var uid string;
	uid, rid = auther.Auth(string(p.Body))
	key = encode(uid, rid);

	p.Body = emptyJSONBody
	p.Operation = define.OP_AUTH_REPLY
	err = p.WriteWebsocket(conn)
	return
}
