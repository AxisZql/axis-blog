package elastic

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {

	query := Param{
		Query: query{
			Bool: Bool{
				Should: []should{
					{
						MultiMatch: multiMatch{
							Query:  "网络",
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
			PreTags:  `<span style=\"color: red\">`,
			PostTags: `</span>`,
		},
	}

	//buf, _ := json.MarshalIndent(query, " ", " ")
	//fmt.Println(string(buf))

	var res QueryResp
	res.Hits.Hits = new([]ArticleHighlight)

	err := Query("myblog", query, &res)
	fmt.Println((*res.Hits.Hits.(*[]ArticleHighlight))[0].Highlight)

	if err != nil {
		t.Fatalf("Test es hightlight search failure：%+v", err)
	}
}
