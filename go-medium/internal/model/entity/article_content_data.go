package entity

type ArticleContentData struct {
	Data struct {
		Post struct {
			Id         string `json:"id"`
			ViewerEdge struct {
				Id          string `json:"id"`
				FullContent struct {
					IsLockedPreviewOnly bool   `json:"isLockedPreviewOnly"`
					ValidatedShareKey   string `json:"validatedShareKey"`
					BodyModel           struct {
						Sections []struct {
							Name            interface{} `json:"name"`
							StartIndex      int         `json:"startIndex"`
							TextLayout      interface{} `json:"textLayout"`
							ImageLayout     interface{} `json:"imageLayout"`
							BackgroundImage interface{} `json:"backgroundImage"`
							VideoLayout     interface{} `json:"videoLayout"`
							BackgroundVideo interface{} `json:"backgroundVideo"`
						} `json:"sections"`
						Paragraphs []struct {
							Id       string      `json:"id"`
							Name     string      `json:"name"`
							Type     string      `json:"type"`
							Href     interface{} `json:"href"`
							Layout   *string     `json:"layout"`
							Metadata *struct {
								Id             string      `json:"id"`
								OriginalHeight int         `json:"originalHeight"`
								OriginalWidth  int         `json:"originalWidth"`
								FocusPercentX  interface{} `json:"focusPercentX"`
								FocusPercentY  interface{} `json:"focusPercentY"`
								Alt            interface{} `json:"alt"`
								Typename       string      `json:"__typename"`
							} `json:"metadata"`
							Text         string      `json:"text"`
							HasDropCap   interface{} `json:"hasDropCap"`
							DropCapImage interface{} `json:"dropCapImage"`
							Markups      []struct {
								Type         string      `json:"type"`
								Start        int         `json:"start"`
								End          int         `json:"end"`
								Href         interface{} `json:"href"`
								AnchorType   interface{} `json:"anchorType"`
								UserId       interface{} `json:"userId"`
								LinkMetadata interface{} `json:"linkMetadata"`
								Typename     string      `json:"__typename"`
							} `json:"markups"`
							Typename          string      `json:"__typename"`
							CodeBlockMetadata interface{} `json:"codeBlockMetadata"`
							Iframe            interface{} `json:"iframe"`
							MixtapeMetadata   interface{} `json:"mixtapeMetadata"`
						} `json:"paragraphs"`
					} `json:"bodyModel"`
					ShareKeyCreator interface{} `json:"shareKeyCreator"`
					Typename        string      `json:"__typename"`
				} `json:"fullContent"`
			} `json:"viewerEdge"`
		} `json:"post"`
	} `json:"data"`
}
