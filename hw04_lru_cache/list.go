package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *Node
	Back() *Node
	PushFront(v interface{}) *Node
	PushBack(v interface{}) *Node
	Remove(i *Node)
	MoveToFront(i *Node)
}

type Node struct {
	// Node - element of List
	Prev  *Node
	Next  *Node
	Value interface{}
}

type list struct {
	// Double Linked list
	front, back *Node
	Length      int
}

func (l *list) Len() int {
	// Length of list
	return l.Length
}

func (l *list) Front() *Node {
	// First Item
	return l.front
}

func (l *list) Back() *Node {
	// Last Item
	return l.back
}

func (l *list) PushFront(v interface{}) *Node {
	n := &Node{}
	if l.Length == 0 {
		n.Prev = nil
		n.Next = nil
		n.Value = v

		l.front = n
		l.back = n
		l.Length++

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
	l.Length++
	return n
}

func (l *list) PushBack(v interface{}) *Node {
	n := &Node{}
	if l.Length == 0 {
		n.Prev = nil
		n.Next = nil
		n.Value = v

		l.front = n
		l.back = n
		l.Length++
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
	l.Length++
	return n
}

func (l *list) Remove(n *Node) {
	if n.Prev == nil {
		n.Next.Prev = nil
		n.Value = nil
		l.front = n.Next
		l.Length--
		return
	}

	if n.Next == nil {
		n.Prev.Next = nil
		n.Value = nil
		l.back = n.Prev
		l.Length--
		return
	}

	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
	n.Value = nil
	l.Length--
}

func (l *list) MoveToFront(n *Node) {
	// Move a node to the front
	transportNode := n
	if n.Prev == nil { // that first node
		return
	}

	if n.Next == nil { // that last node
		head := l.front
		beforeBack := l.back.Prev

		transportNode.Prev.Next = nil
		transportNode.Prev = nil
		transportNode.Next = head
		head.Prev = transportNode
		l.back = beforeBack
		l.front = transportNode
		return
	}

	transportNode.Prev.Next = transportNode.Next
	transportNode.Next.Prev = transportNode.Prev
	transportNode.Prev = nil
	transportNode.Next = l.front
	l.front.Prev = transportNode
	l.front = transportNode
}

func NewList() List {
	// Return NewList
	return &list{}
}
