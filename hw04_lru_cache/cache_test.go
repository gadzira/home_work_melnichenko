package hw04_lru_cache //nolint:golint,stylecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("set/get", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("push-out test", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300)
		_ = c.Set("ddd", 400)

		aaa, ok := c.Get("aaa")
		require.Nil(t, aaa)
		require.False(t, ok)

		ddd, ok := c.Get("ddd")
		require.Equal(t, 400, ddd)
		require.True(t, ok)

	})

	t.Run("rare elements", func(t *testing.T) {
		c := NewCache(3)
		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300)
		fmt.Println(c.getLen())
		_, _ = c.Get("aaa")
		_, _ = c.Get("aaa")
		_ = c.Set("aaa", 567)
		_, _ = c.Get("aaa")
		fmt.Println(c.getLen())
		_ = c.Set("bbb", 1200)
		_, _ = c.Get("bbb")

		_ = c.Set("ddd", 400)

		// fmt.Println(c.getLen())

		ccc, ok := c.Get("ccc")
		require.Nil(t, ccc)
		require.False(t, ok)

	})

}

// func TestCacheMultithreading(t *testing.T) {
// t.Skip() // Remove if task with asterisk completed
//
// c := NewCache(10)
// wg := &sync.WaitGroup{}
// wg.Add(2)
//
// go func() {
// defer wg.Done()
// for i := 0; i < 1_000_000; i++ {
// c.Set(Key(strconv.Itoa(i)), i)
// }
// }()
//
// go func() {
// defer wg.Done()
// for i := 0; i < 1_000_000; i++ {
// c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
// }
// }()
//
// wg.Wait()
// }
