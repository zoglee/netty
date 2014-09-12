package main

import (
	"net"
	"unsafe"
)

const (
	uint32zero = uint32(0)
	uint16zero = uint16(0)

	Uint32Size = int(unsafe.Sizeof(uint32zero))
	Uint16Size = int(unsafe.Sizeof(uint16zero))
	PkgLenSize = Uint32Size + Uint16Size

	// 不包含自身的头部字节数
	MaxAllowedPkgHeadLen      = (1 << 10) // 1KB
	MaxAllowedPkgTotalHeadLen = int(Uint16Size) + int(MaxAllowedPkgHeadLen)
	MinAllowedPkgHeadLen      = 0 // 头部最少字节数
	MinAllowedPkgTotalHeadLen = int(Uint16Size) + int(MinAllowedPkgHeadLen)

	MaxAllowedPkgBodyLen = (8 << 20) // 8M (除开len,hlen,head) body的最大字节数
	MinAllowedPkgBodyLen = 0         // (除开len,hlen,head) body的最少字节数(0即允许空包)

	MaxAllowedPkgLen = PkgLenSize + MaxAllowedPkgHeadLen + MaxAllowedPkgBodyLen
	MinAllowedPkgLen = PkgLenSize + MinAllowedPkgHeadLen + MinAllowedPkgBodyLen //6|0|nil|nil
)

type Request struct {
	Head []byte
	Body []byte
}

type Response struct {
	Head []byte
	Body []byte
}

type Player struct {
}

type FdContext struct {
	fd      int
	tcpconn *net.TCPConn
	inreqc  <-chan *Request
	player  *Player
}

var fdmap = make(map[int]*FdContext, 1000)
