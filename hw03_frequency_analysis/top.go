package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

func SplitString(s string) map[string]int {
	resultDict := make(map[string]int)
	str := strings.Fields(s)
	for _, val := range str {
		if _, found := resultDict[val]; found {
			resultDict[val]++
		} else {
			resultDict[val] = 1
		}
	}
	return resultDict
}

func Top10(text string) []string {
	if len(text) > 0 {
		type kv struct {
			Key   string
			Value int
		}

		var ss []kv
		unSortedMap := SplitString(text)

		for k, v := range unSortedMap {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		keyList := []string{}

		for _, v := range ss {
			keyList = append(keyList, v.Key)
		}

		top10 := keyList[0:10]
		return top10
	}
	return nil
}
