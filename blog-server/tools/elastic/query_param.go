package elastic

import "blog-server/common"

/**
@author: axiszql
@date: 2022-5-13 23:24
@desc: es查询参数定义
*/

// Param 统一高亮查询参数
type Param struct {
	Query     query     `json:"query,omitempty"`
	Highlight highlight `json:"highlight,omitempty"`
}

type query struct {
	Bool Bool `json:"bool"`
}

type highlight struct {
	Fields   interface{} `json:"fields"`    // 自定义要高亮查询的字段
	PreTags  string      `json:"pre_tags"`  // 自定义高亮左符
	PostTags string      `json:"post_tags"` // 自定义高亮右符号
}

type multiMatch struct {
	Query  string   `json:"query"`  // 要全文高亮搜索的字符串
	Fields []string `json:"fields"` // 要进行对目标字符串进行查找的字段
}

type should struct {
	MultiMatch multiMatch `json:"multi_match"`
}

type Bool struct {
	Should             []should `json:"should"`
	MinimumShouldMatch int      `json:"minimum_should_match"` // 最小匹配词长度
}

// QueryResp 统一查询响应
type QueryResp struct {
	Error struct {
		RootCause []struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"root_cause"`
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
	Hits   struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64     `json:"max_score"`
		Hits     interface{} `json:"hits"`
	} `json:"hits"`
}

// 统一高亮查询返回数据前缀
type hist struct {
	Index     string      `json:"_index"`
	Type      string      `json:"_type"`
	ID        string      `json:"_id"`
	Score     float64     `json:"_score"`
	Ignored   []string    `json:"_ignore"`
	Highlight interface{} `json:"highlight"`
}

// ArticleHighlight 文章高亮查询响应数据
type ArticleHighlight struct {
	hist
	Source    common.TArticle `json:"_source"`
	Highlight interface{}     `json:"highlight"`
}
