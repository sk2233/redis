/*
@author: sk
@date: 2024/5/13
*/
package main

import (
	"fmt"
	"net"
	"time"
)

type Service struct {
	Address string
	Kvs     map[string]*Data // 两个可能key重叠
	ZSets   map[string]*ZSet
	Expires *Heap
}

func (s *Service) Run() {
	go s.listen()
	go s.expire()
}

func (s *Service) expire() {
	for {
		if s.Expires.IsEmpty() {
			time.Sleep(time.Second)
			continue
		}
		item := s.Expires.Peek()
		wait := item.Expire - time.Now().Unix()
		if wait > 0 {
			time.Sleep(time.Second * time.Duration(wait))
		} else {
			s.Expires.Pop()
			delete(s.Kvs, item.Name)
			delete(s.ZSets, item.Name)
		}
	}
}

func (s *Service) listen() {
	listen, err := net.Listen("tcp", s.Address)
	HandleErr(err)
	for {
		accept, err := listen.Accept()
		HandleErr(err)
		s.HandleAccept(accept)
	}
}

func (s *Service) HandleAccept(accept net.Conn) {
	go s.handleAccept(accept)
}

func (s *Service) handleAccept(accept net.Conn) {
	defer accept.Close()
	cmd := &Cmd{}
	ReadAny(accept, cmd)
	data := s.HandleCmd(cmd)
	WriteAny(accept, data)
}

func (s *Service) HandleCmd(cmd *Cmd) *Data {
	switch cmd.Cmd {
	case CmdDel:
		return s.StrDel(cmd)
	case CmdSet:
		return s.StrSet(cmd)
	case CmdGet:
		return s.StrGet(cmd)
	case CmdKeys:
		return s.Keys(cmd)
	case CmdZAdd:
		return s.ZSetAdd(cmd)
	case CmdZRem:
		return s.ZSetRem(cmd)
	case CmdZScore:
		return s.ZSetScore(cmd)
	case CmdZQuery:
		return s.ZSetQuery(cmd)
	case CmdExpire:
		return s.Expire(cmd)
	case CmdTTL:
		return s.TTL(cmd)
	default:
		return NewErrData(fmt.Sprintf("unknown command: %v", cmd.Cmd))
	}
}

func (s *Service) StrDel(cmd *Cmd) *Data {
	delete(s.Kvs, cmd.Args[0].MustStr())
	return NewOkData()
}

func (s *Service) StrSet(cmd *Cmd) *Data {
	s.Kvs[cmd.Args[0].MustStr()] = cmd.Args[1]
	return NewOkData()
}

func (s *Service) StrGet(cmd *Cmd) *Data {
	if val, ok := s.Kvs[cmd.Args[0].MustStr()]; ok {
		return val
	} else {
		return NewErrData(fmt.Sprintf("key not found: %v", cmd.Args[0].MustStr()))
	}
}

func (s *Service) Keys(_ *Cmd) *Data {
	arr := make([]*Data, 0)
	for key := range s.Kvs {
		arr = append(arr, NewStrData(key))
	}
	return NewArrData(arr)
}

func (s *Service) ZSetAdd(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	score := cmd.Args[1].MustInt()
	name := cmd.Args[2].MustStr()
	if !Has(s.ZSets, key) {
		s.ZSets[key] = NewZSet()
	}
	s.ZSets[key].Add(score, name)
	return NewOkData()
}

func (s *Service) ZSetRem(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	name := cmd.Args[1].MustStr()
	if Has(s.ZSets, key) {
		s.ZSets[key].Rem(name)
	}
	return NewOkData()
}

func (s *Service) ZSetScore(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	name := cmd.Args[1].MustStr()
	if !Has(s.ZSets, key) {
		return NewErrData(fmt.Sprintf("no zset key = %v", key))
	}
	if res, ok := s.ZSets[key].Score(name); !ok {
		return NewErrData(fmt.Sprintf("zset = %v no name = %v", key, name))
	} else {
		return NewNumData(res)
	}
}

func (s *Service) ZSetQuery(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	score := cmd.Args[1].MustInt()
	name := cmd.Args[2].MustStr()
	offset := cmd.Args[3].MustInt()
	limit := cmd.Args[4].MustInt()
	if !Has(s.ZSets, key) {
		return NewErrData(fmt.Sprintf("no zset key = %v", key))
	}
	items := s.ZSets[key].Query(score, name, offset, limit)
	datas := make([]*Data, 0)
	for _, item := range items {
		datas = append(datas, NewStrData(item.Name))
		datas = append(datas, NewNumData(item.Score))
	}
	return NewArrData(datas)
}

func (s *Service) Expire(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	sec := cmd.Args[1].MustInt()
	s.Expires.Rem(key)
	if sec > 0 {
		s.Expires.Add(key, time.Now().Unix()+int64(sec))
	}
	return NewOkData()
}

func (s *Service) TTL(cmd *Cmd) *Data {
	key := cmd.Args[0].MustStr()
	data := s.Expires.Get(key)
	return NewNumData(int(data))
}

func NewService(address string) *Service {
	return &Service{Address: address, Kvs: make(map[string]*Data), ZSets: make(map[string]*ZSet), Expires: NewHeap()}
}
