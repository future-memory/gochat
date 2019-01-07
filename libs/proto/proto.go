package proto

import (
	"encoding/json"
	"errors"
	"fmt"
	"gochat/libs/bufio"
	"gochat/libs/bytes"
	"gochat/libs/define"
	"gochat/libs/encoding/binary"

	"github.com/gorilla/websocket"
)

const (
	// size
	PackSize      = 4
	HeaderSize    = 2
	OperationSize = 4
	RawHeaderSize = PackSize + HeaderSize + OperationSize;
	// offset
	PackOffset      = 0
	HeaderOffset    = PackOffset + PackSize
	OperationOffset = HeaderOffset + HeaderSize
)

var (
	emptyProto    = Proto{}
	emptyJSONBody = []byte("{}")

	ErrProtoPackLen   = errors.New("default server codec pack length error")
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	ProtoReady  = &Proto{Operation: define.OP_PROTO_READY}
	ProtoFinish = &Proto{Operation: define.OP_PROTO_FINISH}
)

// Proto is a request&response written before every connect.  It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
// tcp:
// binary codec
// websocket & http:
// raw codec, with http header stored ver, operation, seqid
type Proto struct {
	Operation int32           `json:"op"`   // operation for request
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)
}

func (p *Proto) Reset() {
	*p = emptyProto
}

func (p *Proto) String() string {
	return fmt.Sprintf("\n-------- proto --------\nop: %d\nbody: %v\n-----------------------", p.Operation, p.Body)
}

func (p *Proto) WriteTo(b *bytes.Writer) {
	var (
		packLen = RawHeaderSize + int32(len(p.Body))
		buf     = b.Peek(RawHeaderSize)
	)
	binary.BigEndian.PutInt32(buf[PackOffset:], packLen)   //+4
	binary.BigEndian.PutInt16(buf[HeaderOffset:], int16(RawHeaderSize))  //+2
	binary.BigEndian.PutInt32(buf[OperationOffset:], p.Operation)

	if p.Body != nil {
		b.Write(p.Body)
	}
}

func (p *Proto) ReadWebsocket(wr *websocket.Conn) (err error) {
	err = wr.ReadJSON(p)
	return
}

func (p *Proto) WriteBodyTo(b *bytes.Writer) (err error) {
	var (
		ph  Proto
		js  []json.RawMessage
		j   json.RawMessage
		jb  []byte
		bts []byte
	)
	offset := int32(PackOffset)
	buf := p.Body[:]
	for {
		if (len(buf[offset:])) < RawHeaderSize {
			// should not be here
			break
		}
		packLen := binary.BigEndian.Int32(buf[offset : offset+HeaderOffset])
		
		//log.Printf("packLen:%d, %d", packLen, offset+HeaderOffset);

		packBuf := buf[offset : offset+packLen]
		// packet
		ph.Operation = binary.BigEndian.Int32(packBuf[OperationOffset:RawHeaderSize])
		ph.Body = packBuf[RawHeaderSize:]
		if jb, err = json.Marshal(&ph); err != nil {
			return
		}
		j = json.RawMessage(jb)
		js = append(js, j)
		offset += packLen
	}
	if bts, err = json.Marshal(js); err != nil {
		return
	}
	b.Write(bts)
	return
}

func (p *Proto) WriteWebsocket(wr *websocket.Conn) (err error) {
	if p.Body == nil {
		p.Body = emptyJSONBody
	}
	// [{"ver":1,"op":8,"seq":1,"body":{}}, {"ver":1,"op":3,"seq":2,"body":{}}]
	if p.Operation == define.OP_RAW {
		// batch mod
		var b = bytes.NewWriterSize(len(p.Body) + 40*RawHeaderSize)
		if err = p.WriteBodyTo(b); err != nil {
			return
		}
		err = wr.WriteMessage(websocket.TextMessage, b.Buffer())
		return
	}
	err = wr.WriteJSON([]*Proto{p})
	return
}
