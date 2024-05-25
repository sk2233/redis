/*
@author: sk
@date: 2024/5/13
*/
package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type ByteData struct {
	Size uint64
	Data []byte
}

func NewByteData(data []byte) *ByteData {
	return &ByteData{Data: data, Size: uint64(len(data))}
}

type DataType int

const (
	TypeStr DataType = iota + 1
	TypeErr
	TypeNum
	TypeArr
)

type Data struct {
	Type DataType
	Str  string
	Err  string
	Num  int
	Arr  []*Data
}

func NewArrData(arr []*Data) *Data {
	return &Data{Arr: arr, Type: TypeArr}
}

func NewNumData(num int) *Data {
	return &Data{Num: num, Type: TypeNum}
}

func NewStrData(str string) *Data {
	return &Data{Str: str, Type: TypeStr}
}

func NewErrData(err string) *Data {
	return &Data{Err: err, Type: TypeErr}
}

func (d *Data) MustStr() string {
	if d.Type != TypeStr {
		panic(fmt.Sprintf("data type %d not str", d.Type))
	}
	return d.Str
}

func (d *Data) MustErr() string {
	if d.Type != TypeErr {
		panic(fmt.Sprintf("data type %d not err", d.Type))
	}
	return d.Err
}

func (d *Data) MustInt() int {
	if d.Type != TypeNum {
		panic(fmt.Sprintf("data type %d not num", d.Type))
	}
	return d.Num
}

func (d *Data) MustArr() []*Data {
	if d.Type != TypeArr {
		panic(fmt.Sprintf("data type %d not arr", d.Type))
	}
	return d.Arr
}

func NewOkData() *Data {
	return &Data{Type: TypeStr, Str: "ok"}
}

func (d *Data) String() string {
	switch d.Type {
	case TypeStr:
		return d.Str
	case TypeErr:
		return fmt.Sprintf("(%s)", d.Err)
	case TypeNum:
		return fmt.Sprintf("<%d>", d.Num)
	case TypeArr:
		buff := strings.Builder{}
		buff.WriteString("[")
		for i, item := range d.Arr {
			if i > 0 {
				buff.WriteString(",")
			}
			buff.WriteString(item.String())
		}
		buff.WriteString("]")
		return buff.String()
	default:
		panic(fmt.Sprintf("Unknown data type %d", d.Type))
	}
}

type Cmd struct {
	Cmd  string
	Args []*Data
}

type SkipItem struct {
	Name  string
	Score int
	Pre   [SkipDep]*SkipItem
	Next  [SkipDep]*SkipItem
}

func (s *SkipItem) String() string {
	return fmt.Sprintf("<%s:%d>", s.Name, s.Score)
}

func NewSkipItem(name string, score int) *SkipItem {
	return &SkipItem{Name: name, Score: score}
}

type SkipList struct { // 单独使用没什么意义，要放在 zset中使用才行
	Head *SkipItem // 按 score从小到大排序
}

func (l *SkipList) Add(score int, name string) {
	if l.Head == nil {
		l.Head = NewSkipItem(name, score)
	} else if score < l.Head.Score {
		old := l.Head
		l.Head = NewSkipItem(name, score)
		l.buildRef(nil, l.Head, old)
	} else {
		l.add(score, name, l.Head, SkipDep-1)
	}
}

func (l *SkipList) add(score int, name string, item *SkipItem, dep int) {
	lastItem := item
	for item != nil && item.Score < score {
		lastItem = item
		item = item.Next[dep]
	}
	if item == nil { // 到末尾了
		if dep > 0 { // 还有机会下一层
			l.add(score, name, lastItem, dep-1)
		} else { // 没有机会了就是末尾
			l.buildRef(lastItem, NewSkipItem(name, score), nil)
		}
		return
	}
	if item.Score == score { // 相同直接插入不用下一层了
		l.buildRef(item, NewSkipItem(name, score), item.Next[0]) // 这里必须使用都有的那一层
		return
	}
	if dep > 0 { // 当前节点较大，找到上一个节点进入下一层
		l.add(score, name, lastItem, dep-1)
	} else { // 确定在当前节点与其前一个节点之间
		l.buildRef(item.Pre[0], NewSkipItem(name, score), item)
	}
}

func (l *SkipList) Get(score int, name string) *SkipItem {
	if l.Head == nil {
		return nil
	}
	return l.get(score, name, l.Head, SkipDep-1)
}

func (l *SkipList) get(score int, name string, item *SkipItem, dep int) *SkipItem {
	lastItem := item
	for item != nil && item.Score < score {
		lastItem = item
		item = item.Next[dep]
	}
	if item == nil {
		if dep > 0 {
			return l.get(score, name, lastItem, dep-1)
		} else {
			return nil
		}
	}
	if item.Score == score { // 相同进行遍历第一层
		for item.Pre[0] != nil && item.Pre[0].Score == score {
			item = item.Pre[0]
		}
		for item != nil {
			if item.Score != score {
				break
			}
			if item.Name == name {
				return item
			}
			item = item.Next[0]
		}
		return nil
	}
	if dep > 0 {
		return l.get(score, name, lastItem, dep-1)
	} else {
		return nil
	}
}

func (l *SkipList) Del(score int, name string) {
	item := l.Get(score, name)
	if item == nil {
		return
	}
	for i := 0; i < SkipDep; i++ {
		if item.Next[i] != nil {
			item.Next[i].Pre[i] = item.Pre[i]
		}
		if item.Pre[i] != nil {
			item.Pre[i].Next[i] = item.Next[i]
		}
	}
}

// curr不能为 nil 其他的可以为 nil
func (l *SkipList) buildRef(pre, curr, next *SkipItem) {
	sum := 1
	for i := 0; i < SkipDep; i++ {
		if rand.Intn(sum) == 0 {
			curr.Next[i] = next
			if next != nil {
				next.Pre[i] = curr
			}
			curr.Pre[i] = pre
			if pre != nil {
				pre.Next[i] = curr
			}
		}
		sum *= 2
	}
}

func NewSkipList() *SkipList {
	return &SkipList{}
}

type ZSet struct {
	Scores map[string]int
	Order  *SkipList
}

func (s *ZSet) Add(score int, name string) {
	if Has(s.Scores, name) {
		s.Order.Del(s.Scores[name], name) // 若已经存在需要先移除
	}
	s.Scores[name] = score
	s.Order.Add(score, name)
}

func (s *ZSet) Rem(name string) {
	if score, ok := s.Scores[name]; ok {
		delete(s.Scores, name)
		s.Order.Del(score, name)
	}
}

func (s *ZSet) Score(name string) (int, bool) {
	res, ok := s.Scores[name]
	return res, ok
}

func (s *ZSet) Query(score int, name string, offset int, limit int) []*SkipItem {
	item := s.Order.Get(score, name)
	res := make([]*SkipItem, 0)
	for item != nil && offset > 0 {
		item = item.Next[0]
		offset--
	}
	for item != nil && limit > 0 {
		res = append(res, item)
		item = item.Next[0]
		limit--
	}
	return res
}

func NewZSet() *ZSet {
	return &ZSet{Scores: make(map[string]int), Order: NewSkipList()}
}

type HeapItem struct {
	Name   string // 同时作用于所有key
	Expire int64  // 过期时间
}

func NewHeapItem(name string, expire int64) *HeapItem {
	return &HeapItem{Name: name, Expire: expire}
}

type Heap struct {
	Items []*HeapItem
}

func (h *Heap) Add(name string, expire int64) {
	h.Items = append(h.Items, NewHeapItem(name, expire))
	h.up(len(h.Items) - 1)
}

func (h *Heap) up(index int) {
	if index == 0 {
		return
	}
	parent := (index - 1) / 2
	if h.Items[index].Expire < h.Items[parent].Expire {
		h.Items[index], h.Items[parent] = h.Items[parent], h.Items[index]
		h.up(parent)
	}
}

func (h *Heap) Rem(name string) {
	for i := 0; i < len(h.Items); i++ {
		if h.Items[i].Name == name {
			h.rem(i)
			break
		}
	}
}

func (h *Heap) rem(index int) {
	last := len(h.Items) - 1
	h.Items[index], h.Items[last] = h.Items[last], h.Items[index]
	h.Items = h.Items[:last]
	h.down(index)
}

func (h *Heap) down(index int) {
	left := index*2 + 1
	right := index*2 + 2
	target := index
	if left < len(h.Items) && h.Items[left].Expire < h.Items[target].Expire {
		target = left
	}
	if right < len(h.Items) && h.Items[right].Expire < h.Items[target].Expire {
		target = right
	}
	if target == index {
		return
	}
	h.Items[index], h.Items[target] = h.Items[target], h.Items[index]
	h.down(target)
}

func (h *Heap) Pop() *HeapItem {
	res := h.Items[0]
	h.rem(0)
	return res
}

func (h *Heap) Peek() *HeapItem {
	return h.Items[0]
}

func (h *Heap) IsEmpty() bool {
	return len(h.Items) == 0
}

func (h *Heap) Get(name string) int64 {
	for _, item := range h.Items {
		if item.Name == name {
			return item.Expire
		}
	}
	return -1
}

func NewHeap() *Heap {
	return &Heap{Items: make([]*HeapItem, 0)}
}
