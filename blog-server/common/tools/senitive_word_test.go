package tools

import (
	"fmt"
	"testing"
)

func TestReadSenitiveWordFile(t *testing.T) {
	readSenitiveWordFile()
}

func TestGenerateSenitiveForest(t *testing.T) {
	forest := GetSenitiveForest()
	// all, _ := readSenitiveWordFile()
	// str := strings.Join(all, "")
	ans := forest.GetSenitiveWord("sb傻逼，adult")
	fmt.Println(ans)

}
