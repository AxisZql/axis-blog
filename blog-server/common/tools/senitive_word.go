package tools

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

/*
@author AxisZql
@date 2022-4-23 11:20 AM
@desc 前缀树进行敏感词检测
*/

type treeNode struct {
	Value      rune               // 单个敏感字符
	IsEnd      bool               // 是否是叶子节点
	ExistChild map[rune]*treeNode // 记录唯一的孩子节点
}

// Fores 记录所有敏感树的森林
type Forest struct {
	Tree map[rune]*treeNode //存放敏感树树根和数根值之间的映射关系
}

func readSenitiveWordFile() ([]string, error) {
	data, err := ioutil.ReadFile("./sensitive-words.txt")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	mp := make(map[string]struct{})
	dataLine := strings.Split(string(data), "\n")
	senitiveWord := make([]string, 0)
	// 过滤重复的词
	for _, val := range dataLine {
		if _, ok := mp[val]; !ok && val != "" {
			mp[val] = struct{}{}
			senitiveWord = append(senitiveWord, val)
		}
	}
	return senitiveWord, nil
}

func initSenitiveForest() *Forest {
	senitiveWordList, err := readSenitiveWordFile()
	if err != nil {
		panic(err)
	}
	forest := Forest{
		Tree: make(map[rune]*treeNode),
	}
	for _, word := range senitiveWordList {
		wo := []rune(word)
		forest.generateNode(wo)
	}
	return &forest
}

// 根据当前敏感词构建敏感树并把其加入森林
func (forest *Forest) generateNode(wo []rune) {
	if len(wo) == 0 {
		return
	}
	node := forest.Tree[wo[0]]
	//如果没有相同开头的敏感树
	if node == nil {
		node = &treeNode{
			Value:      wo[0],
			ExistChild: make(map[rune]*treeNode),
		}
		forest.Tree[wo[0]] = node
	}
	for i := 1; i < len(wo); i++ {
		if _, ok := node.ExistChild[wo[i]]; !ok {
			newNode := &treeNode{
				Value:      wo[i],
				ExistChild: make(map[rune]*treeNode),
			}
			if i == len(wo)-1 {
				newNode.IsEnd = true
			}
			node.ExistChild[wo[i]] = newNode
			node = newNode
		} else {
			node = node.ExistChild[wo[i]]
		}
	}
}

// GetSenitiveWord 获取文本中的敏感词
func (forest *Forest) GetSenitiveWord(text string) []string {
	senitiveWordList := make([]string, 0)
	_text := []rune(text)
	for i := 0; i < len(_text); i++ {
		length := forest.getSentiveWordLength(_text[i:])
		if length != 0 {
			senitiveWordList = append(senitiveWordList, string(_text[i:i+length]))
			// 跳过检测处理的关键词（最短匹配策略）
			i = i + length - 1
		}
	}
	return senitiveWordList
}

// 从子文本中提取敏感词长度
func (forest *Forest) getSentiveWordLength(subText []rune) int {
	var count int
	if len(subText) == 0 {
		return count
	}
	node := forest.Tree[subText[0]]
	// 不是敏感词
	if node == nil {
		return count
	}
	count++
	flag := false
	if node.IsEnd {
		flag = true
	}
	for i := 1; i < len(subText); i++ {
		node = node.ExistChild[subText[i]]
		if node != nil && node.IsEnd {
			count++
			break
		} else if node != nil {
			if i == len(subText)-1 {
				count = 0
				break
			}
			count++
		} else {
			count = 0
			break
		}
	}
	if !flag && count == 1 {
		count = 0
	}
	return count
}

var senitiveForest *Forest
var once sync.Once

// 全局初始化敏感词汇森林（单例模式）
func GetSenitiveForest() *Forest {
	once.Do(func() {
		senitiveForest = initSenitiveForest()
	})
	return senitiveForest
}
