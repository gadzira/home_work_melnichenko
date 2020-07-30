package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

func SplitString(s string) map[string]int {
	resultDict := make(map[string]int)
	str := strings.Fields(s)
	for _, val := range str {
		resultDict[val]++
	}
	return resultDict
}

type kv struct {
	Key   string
	Value int
}

func Top10(text string) []string {
	if len(text) == 0 {
		return nil
	}
	unSortedMap := SplitString(text)
	ss := make([]kv, len(unSortedMap))

	for k, v := range unSortedMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	keyList := make([]string, 0, len(ss))
	for _, v := range ss {
		keyList = append(keyList, v.Key)
	}

	if len(keyList) >= 10 {
		return keyList[0:10]
	}
	return keyList
}
