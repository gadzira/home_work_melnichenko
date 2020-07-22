package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

// List ...
type List interface {
	Len() int                      // длина списка
	Front() *Node                  // первый Item вернуть первый элемент списка
	Back() *Node                   // последний Item вернуть последний элемент списка
	PushFront(v interface{}) *Node // добавить значение в начало
	PushBack(v interface{}) *Node  // добавить значение в конец
	Remove(i *Node)                // удалить элемент
	MoveToFront(i *Node)           // переместить элемент в начало
}

// Node ...
type Node struct {
	Prev, Next *Node
	Value      interface{}
}

// Видимо тот самый двусвязанный лист
type list struct {
	Lock        *sync.RWMutex
	front, back *Node
	Lenght      int
}

// длина списка
func (l *list) Len() int {
	return l.Lenght
}

// первый Item
func (l *list) Front() *Node {
	return l.front
}

// последний Item
func (l *list) Back() *Node {
	return l.back
}

func (l *list) PushFront(v interface{}) *Node {
	// l.Lock.Lock()
	// defer l.Lock.Lock()

	n := &Node{}
	if l.Lenght == 0 {
		n.Prev = nil
		n.Next = nil
		n.Value = v

		l.front = n
		l.back = n
		l.Lenght++

		return n
	}

	currentHead := l.front
	n.Prev = nil
	n.Next = currentHead
	n.Value = v

	currentHead.Prev = n

	l.front = n
	if l.back == nil {
		l.back = currentHead
	}
	l.Lenght++
	return n
}

func (l *list) PushBack(v interface{}) *Node {
	// l.Lock.Lock()
	// defer l.Lock.Lock()

	n := &Node{}
	if l.Lenght == 0 {
		n.Prev = nil
		n.Next = nil
		n.Value = v

		l.front = n
		l.back = n
		l.Lenght++
		return n
	}
	currentBack := l.back
	n.Prev = currentBack
	n.Next = nil
	n.Value = v

	currentBack.Next = n

	l.back = n
	if l.front == nil {
		l.front = currentBack
	}
	l.Lenght++
	return n
}

func (l *list) Remove(n *Node) {

	if n.Prev == nil {
		n.Next.Prev = nil
		n.Value = nil
		l.front = n.Next
		l.Lenght--
		return
	}

	if n.Next == nil {
		n.Prev.Next = nil
		n.Value = nil
		l.back = n.Prev
		l.Lenght--
		return
	}

	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
	n.Value = nil
	l.Lenght--
}

// переместить элемент в начало
func (l *list) MoveToFront(n *Node) {

	transportNode := n

	if n.Prev == nil { // that first node
		return
	}

	if n.Next == nil { // that last node
		// transportNode := n
		transportNode.Prev.Next = nil

		transportNode.Prev = nil
		transportNode.Next = l.front

		l.front.Prev = transportNode
		l.front = transportNode
		return
	}

	// замыкаем prev и next друг на друга, те изымаем текущую ноду из списка
	transportNode.Prev.Next = transportNode.Next
	transportNode.Next.Prev = transportNode.Prev

	// "делаем" перемещаемую ноду первой: а) обновляем ее указатели
	transportNode.Prev = nil
	transportNode.Next = l.front

	// б) фиксуруем изменения в списке
	l.front.Prev = transportNode
	l.front = transportNode

}

// NewList ...
func NewList() List {
	return &list{}
}
