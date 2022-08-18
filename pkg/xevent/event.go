package xevent

import (
	"fmt"
	"log"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// 事件ID
const (
	LicenseXXXEvent = iota + 1 // 事件ID
)

var eventMap = make(map[int]*link)
var eventMutex = sync.RWMutex{}

type Handle func()

// Register 事件注册
func Register(eventID int, handle Handle) (uuid string) {
	eventMutex.Lock()
	defer eventMutex.Unlock()
	if v, ok := eventMap[eventID]; ok && v != nil {
		uuid = v.add(handle)
	} else {
		eventMap[eventID] = newLink(handle)
		uuid = eventMap[eventID].UUID
	}
	return
}

// Remove 删除注册事件
func Remove(eventID int, uuid string) {
	eventMutex.Lock()
	defer eventMutex.Unlock()
	if v, ok := eventMap[eventID]; ok && v != nil {
		if v.UUID == uuid {
			v = v.Next
		} else {
			v.remove(uuid)
		}
	}
}

func Debug() {
	for _, v := range eventMap {
		v.string()
		fmt.Println(v.len())
	}
}

// Event 事件触发
func Event(eventID int) {
	eventMutex.RLock()
	defer eventMutex.RUnlock()
	if v, ok := eventMap[eventID]; ok && v != nil {
		v.do()
	}
}

type link struct {
	UUID string
	Func Handle
	Next *link
}

func newLink(h Handle) *link {
	return &link{
		UUID: uuid.NewV4().String(),
		Func: h,
	}
}

func (l *link) add(h Handle) (uuid string) {
	if l.Next == nil {
		l.Next = newLink(h)
		uuid = l.Next.UUID
	} else {
		l.Next.add(h)
	}
	return
}

func (l *link) len() int {
	tmp := l
	counter := 1
	for {
		if tmp.Next != nil {
			tmp = tmp.Next
			counter++
		} else {
			break
		}
	}
	return counter
}

func (l *link) string() {
	tmp := l
	for {
		if tmp != nil {
			fmt.Println(tmp.UUID)
			tmp = tmp.Next
		} else {
			break
		}
	}
}

func (l *link) remove(uuid string) {
	tmp := l.Next
	prev := l
	for {
		if tmp == nil {
			break
		}
		if tmp.UUID == uuid {
			prev.Next = tmp.Next
			break
		} else {
			tmp = tmp.Next
			prev = prev.Next
		}
	}
}

func (l *link) do() {
	if l.Func != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Event %s occur panic:%v\n", l.UUID, r)
				}
			}()
			l.Func()
		}()
	}
	if l.Next != nil {
		l.Next.do()
	}
}
