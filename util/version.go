package util

import "strings"

func Compare(destVersion string, sourceVersion string) int {

	dvs := strings.Split(destVersion, ".")
	svs := strings.Split(sourceVersion, ".")
	dLen := len(dvs)
	sLen := len(svs)
	lLen := dLen
	if sLen < lLen {
		lLen = sLen
	}
	for i := 0; i < lLen; i++ {
		vd := dvs[i]
		vs := svs[i]
		if len(vd) > len(vs) {
			return 1
		} else if len(vd) < len(vs) {
			return -1
		} else {
			c:=strings.Compare(vd,vs)
			if c!=0{
				return c
			}
		}
	}
	if dLen > sLen {
		return 1
	}
	if sLen > dLen {
		return -1
	}
	return 0
}
