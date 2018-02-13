package utils

import (
	"github.com/deckarep/golang-set"
)

func StringArrUion(strs1, strs2 []string) (strs []string) {
	set1 := mapset.NewSet()
	for _, str := range strs1 {
		set1.Add(str)
	}
	set2 := mapset.NewSet()
	for _, str := range strs2 {
		set2.Add(str)
	}
	set := set1.Union(set2)
	for _, str := range set.ToSlice() {
		strs = append(strs, str.(string))
	}
	return
}
