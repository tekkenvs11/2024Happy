package entity

type ArticleTagData struct {
	Data struct {
		TagFromSlug struct {
			Id         string `json:"id"`
			SortedFeed struct {
				Edges []struct {
					Cursor string  `json:"cursor"`
					Node   Article `json:"node"`
				} `json:"edges"`
			} `json:"sortedFeed"`
		} `json:"tagFromSlug"`
	} `json:"data"`
}
