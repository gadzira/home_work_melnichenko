package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()
		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push front", func(t *testing.T) {
		l := NewList()
		l.PushFront(7)
		l.PushFront(14)
		l.PushFront(21)
		require.Equal(t, l.Len(), 3)
		require.Equal(t, 21, l.Front().Value)
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()
		l.PushBack(3)
		l.PushBack(13)
		require.Equal(t, l.Len(), 2)
		require.Equal(t, 13, l.Back().Value)
	})

	t.Run("remove some node", func(t *testing.T) {
		l := NewList()
		l.PushFront(7)                  // 7
		l.PushBack(3)                   // [7 3]
		l.PushBack(13)                  // [7 3 13]
		nodeForDelete := l.Front().Next // 3
		l.Remove(nodeForDelete)         // [7 13]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, 7, l.Front().Value)
		require.Equal(t, 13, l.Back().Value)
	})

	t.Run("remove front node", func(t *testing.T) {
		l := NewList()
		l.PushFront(37)        // [37]
		l.PushFront(36)        // [36 37]
		l.PushBack(1300)       // [36 37 1300]
		frontNode := l.Front() //  36
		l.Remove(frontNode)    // [37 1300]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, 37, l.Front().Value)
	})

	t.Run("remove back node", func(t *testing.T) {
		l := NewList()
		l.PushFront(14)      // [14]
		l.PushBack(85)       // [14 85]
		l.PushBack(73)       // [14 85 73]
		backNode := l.Back() // 73
		l.Remove(backNode)   // [14 85]
		require.Equal(t, l.Len(), 2)
		require.Equal(t, 85, l.Back().Value)
	})

	t.Run("move midlle node to front", func(t *testing.T) {
		l := NewList()
		l.PushFront(14)              // [14]
		l.PushBack(85)               // [14 85]
		l.PushBack(73)               // [14 85 73]
		midlleNode := l.Front().Next // 85
		l.MoveToFront(midlleNode)    // [85 14 73]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, 85, l.Front().Value)
		require.Equal(t, 14, l.Front().Next.Value)
		require.Equal(t, 73, l.Back().Value)
	})

	t.Run("move last node to front", func(t *testing.T) {
		l := NewList()
		l.PushFront(14)         // [14]
		l.PushBack(85)          // [14 85]
		l.PushBack(73)          // [14 85 73]
		lastNode := l.Back()    // 73
		l.MoveToFront(lastNode) // [73 14 85]
		require.Equal(t, l.Len(), 3)
		require.Equal(t, 73, l.Front().Value)
		require.Equal(t, 14, l.Front().Next.Value)
		require.Equal(t, 85, l.Back().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()
		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, l.Len(), 2)
		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]
		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
