package main

import "strings"

// 파일명의 상위 dir을 구한다.
func GetParentDir(path string, sep string) string {
	slice := strings.Split(path, sep)
	return strings.Join(slice[:len(slice)-1], sep) + sep
}

// 공통 dir을 구한다.
func GetSameDir(a string, b string) string {

	// 시작, 끝의 디렉토리 제거
	a = strings.Trim(a, "/")
	b = strings.Trim(b, "/")

	aa := strings.Split(a, "/")
	bb := strings.Split(b, "/")

	if len(aa) > len(bb) {
		aa, bb = bb, aa
	}

	common := ""
	for i := range aa {
		if aa[i] == bb[i] {
			common += "/" + aa[i]
		}
	}

	return common + "/"
}
