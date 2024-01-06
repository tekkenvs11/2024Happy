package entity

type Article struct {
	Id        string `json:"id"`
	ClapCount int    `json:"clapCount"`
	MediumUrl string `json:"mediumUrl"`
	Title     string `json:"title"`
}

type ArticleContentSub struct {
	HtmlTag string `json:"htmlTag"`
	Text    string `json:"text"`
}

type ArticleContent struct {
	Id         string              `json:"id"`
	ZhContents []ArticleContentSub `json:"zhContents"`
	Contents   []ArticleContentSub `json:"contents"`
}

type ArticleTranslate struct {
	Result struct {
		From        string `json:"from"`
		TransResult []struct {
			Dst string `json:"dst"`
			Src string `json:"src"`
		} `json:"trans_result"`
		To string `json:"to"`
	} `json:"result"`
	LogId int64 `json:"log_id"`
}
