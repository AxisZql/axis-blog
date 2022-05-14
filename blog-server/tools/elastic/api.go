package elastic

import (
	"fmt"
	"github.com/pkg/errors"
)

func HighlightQueryArticle(word string) (resp *[]ArticleHighlight, err error) {
	query := Param{
		Query: query{
			Bool: Bool{
				Should: []should{
					{
						MultiMatch: multiMatch{
							Query:  word,
							Fields: []string{"article_content", "article_title"},
						},
					},
				},
			},
		},
		Highlight: highlight{
			Fields: struct {
				ArticleContent struct{} `json:"article_content"`
				ArticleTitle   struct{} `json:"article_title"`
			}{struct{}{}, struct{}{}},
			PreTags:  `<span style='color: red'>`,
			PostTags: `</span>`,
		},
	}
	var res QueryResp
	res.Hits.Hits = new([]ArticleHighlight)
	err = Query("myblog", query, &res)
	if err != nil {
		return
	}
	if res.Status != 0 {
		err = errors.New(fmt.Sprintf("Status: %d Error: %+v", res.Status, res.Error))
		return
	}

	// 因为Hits.Hits是interface所以要进行类型的强制转换
	resp = res.Hits.Hits.(*[]ArticleHighlight)
	return
}
