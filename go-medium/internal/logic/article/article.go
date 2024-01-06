package article

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	v1 "go-medium/api/article/v1"
	"go-medium/internal/consts"
	"go-medium/internal/model/entity"
	"go-medium/internal/service"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	cacheCtx = gctx.New()
	cache    = gcache.New()
)

type sArticle struct {
}

func init() {
	service.RegisterArticle(New())
}

func New() *sArticle {
	return &sArticle{}
}

func (s *sArticle) ArticleGetList(ctx context.Context, req *v1.ArticleGetListReq) (res *v1.ArticleGetListRes, err error) {
	g.Log().Info(ctx, "do service ArticleGetList .")

	client := g.Client()
	if req.Host == "" {
		client.SetProxy(fmt.Sprintf("socks5://%s:%s@%s:%d", consts.PROXY_USER, consts.PROXY_PASSWD, consts.PROXY_URL, consts.PROXY_PORT))
	}

	var result *gclient.Response

	postJson := `{"operationName":"TagArchiveFeedQuery",` +
		`"query":"query TagArchiveFeedQuery($tagSlug: String!, $timeRange: TagPostsTimeRange!, $sortOrder: TagPostsSortOrder!, $first: Int!, $after: String) {\n   tagFromSlug(tagSlug: $tagSlug) {\n     id\n     sortedFeed: posts(\n       timeRange: $timeRange\n       sortOrder: $sortOrder\n       first: $first\n       after: $after\n     ) {\n       ...TagPosts_tagPostConnection\n     }\n   }\n }\n \n fragment TagPosts_tagPostConnection on TagPostConnection {\n   edges {\n     cursor\n     node {\n       id\n       clapCount\n       mediumUrl\n       title\n     }\n   }\n }",`

	if req.Year == 0 {
		postJson = postJson + `"variables":{"tagSlug": "software-engineering", "timeRange": {   "kind": "ALL_TIME" }, "sortOrder": "MOST_READ", "first": 10, "after": ""}}`
	} else {
		postJson = postJson + fmt.Sprintf(`"variables":{"tagSlug": "software-engineering", "timeRange": {   "kind": "IN_YEAR", "inYear": { "year" : %d } }, "sortOrder": "MOST_READ", "first": 10, "after": ""}}`, req.Year)
	}

	result, err = client.Post(context.Background(), consts.BASE_URI, postJson)
	if err != nil {
		g.Log().Errorf(ctx, "get article list error. %s.", err)
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "get article list error.")
	}
	defer result.Close()
	var atd *entity.ArticleTagData
	//fmt.Println(result.ReadAllString())
	if err = gjson.DecodeTo(result.ReadAll(), &atd); err != nil {
		g.Log().Errorf(ctx, "scan  decodeTo articletagdata error.%s", err.Error())
		return nil, gerror.WrapCode(gcode.CodeNotImplemented, err, "scan  decodeTo articletagdata error.")
	}
	res = &v1.ArticleGetListRes{
		Lists: make([]entity.Article, 0),
	}
	for _, edge := range atd.Data.TagFromSlug.SortedFeed.Edges {
		res.Lists = append(res.Lists, edge.Node)
		err = cache.Set(cacheCtx, edge.Node.Id, edge.Node, 0)
		if err != nil {
			return nil, gerror.WrapCode(gcode.CodeNotImplemented, err, "set cache error.")
		}
	}

	return res, err
}

func (s *sArticle) ArticleGetContent(ctx context.Context, req *v1.ArticleGetContentReq) (res *v1.ArticleGetContentRes, err error) {

	g.Log().Infof(ctx, "do service ArticleGetContent id:%s", req.Id)

	var article entity.Article

	client := g.Client()
	client.SetProxy(fmt.Sprintf("socks5://%s:%s@%s:%d", req.UserName, req.Password, req.Host, req.Port))

	if value, e := cache.Get(ctx, req.Id); e != nil {
		return nil, gerror.WrapCode(gcode.CodeNotImplemented, err, "read  cache id error.")
	} else {
		e = value.Scan(&article)
		if e != nil {
			return nil, gerror.WrapCode(gcode.CodeNotImplemented, err, "read  cache id error.")
		}
	}

	baseUri := gstr.StrTill(article.MediumUrl, "//") + "/" + gstr.StrTill(gstr.StrEx(article.MediumUrl, "//"), "/")

	var result *gclient.Response
	result, err = client.Post(context.Background(), baseUri+"_/graphql",
		fmt.Sprintf(`{"operationName":"PostViewerEdgeContentQuery","variables":{"postId":"%s","postMeteringOptions":{"referrer":""}},"query":"query PostViewerEdgeContentQuery($postId: ID!, $postMeteringOptions: PostMeteringOptions) {\n  post(id: $postId) {\n    ... on Post {\n      id\n      viewerEdge {\n        id\n        fullContent(postMeteringOptions: $postMeteringOptions) {\n          isLockedPreviewOnly\n          validatedShareKey\n          bodyModel {\n            ...PostBody_bodyModel\n            __typename\n          }\n          ...FriendLinkMeter_postContent\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment PostBody_bodyModel on RichText {\n  sections {\n    name\n    startIndex\n    textLayout\n    imageLayout\n    backgroundImage {\n      id\n      originalHeight\n      originalWidth\n      __typename\n    }\n    videoLayout\n    backgroundVideo {\n      videoId\n      originalHeight\n      originalWidth\n      previewImageId\n      __typename\n    }\n    __typename\n  }\n  paragraphs {\n    id\n    ...PostBodySection_paragraph\n    __typename\n  }\n  ...normalizedBodyModel_richText\n  __typename\n}\n\nfragment PostBodySection_paragraph on Paragraph {\n  name\n  ...PostBodyParagraph_paragraph\n  __typename\n  id\n}\n\nfragment PostBodyParagraph_paragraph on Paragraph {\n  name\n  type\n  ...ImageParagraph_paragraph\n  ...TextParagraph_paragraph\n  ...IframeParagraph_paragraph\n  ...MixtapeParagraph_paragraph\n  ...CodeBlockParagraph_paragraph\n  __typename\n  id\n}\n\nfragment ImageParagraph_paragraph on Paragraph {\n  href\n  layout\n  metadata {\n    id\n    originalHeight\n    originalWidth\n    focusPercentX\n    focusPercentY\n    alt\n    __typename\n  }\n  ...Markups_paragraph\n  ...ParagraphRefsMapContext_paragraph\n  ...PostAnnotationsMarker_paragraph\n  __typename\n  id\n}\n\nfragment Markups_paragraph on Paragraph {\n  name\n  text\n  hasDropCap\n  dropCapImage {\n    ...MarkupNode_data_dropCapImage\n    __typename\n    id\n  }\n  markups {\n    ...Markups_markup\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment MarkupNode_data_dropCapImage on ImageMetadata {\n  ...DropCap_image\n  __typename\n  id\n}\n\nfragment DropCap_image on ImageMetadata {\n  id\n  originalHeight\n  originalWidth\n  __typename\n}\n\nfragment Markups_markup on Markup {\n  type\n  start\n  end\n  href\n  anchorType\n  userId\n  linkMetadata {\n    httpStatus\n    __typename\n  }\n  __typename\n}\n\nfragment ParagraphRefsMapContext_paragraph on Paragraph {\n  id\n  name\n  text\n  __typename\n}\n\nfragment PostAnnotationsMarker_paragraph on Paragraph {\n  ...PostViewNoteCard_paragraph\n  __typename\n  id\n}\n\nfragment PostViewNoteCard_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment TextParagraph_paragraph on Paragraph {\n  type\n  hasDropCap\n  codeBlockMetadata {\n    mode\n    lang\n    __typename\n  }\n  ...Markups_paragraph\n  ...ParagraphRefsMapContext_paragraph\n  __typename\n  id\n}\n\nfragment IframeParagraph_paragraph on Paragraph {\n  type\n  iframe {\n    mediaResource {\n      id\n      iframeSrc\n      iframeHeight\n      iframeWidth\n      title\n      __typename\n    }\n    __typename\n  }\n  layout\n  ...Markups_paragraph\n  __typename\n  id\n}\n\nfragment MixtapeParagraph_paragraph on Paragraph {\n  type\n  mixtapeMetadata {\n    href\n    mediaResource {\n      mediumCatalog {\n        id\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  ...GenericMixtapeParagraph_paragraph\n  __typename\n  id\n}\n\nfragment GenericMixtapeParagraph_paragraph on Paragraph {\n  text\n  mixtapeMetadata {\n    href\n    thumbnailImageId\n    __typename\n  }\n  markups {\n    start\n    end\n    type\n    href\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment CodeBlockParagraph_paragraph on Paragraph {\n  codeBlockMetadata {\n    lang\n    mode\n    __typename\n  }\n  __typename\n  id\n}\n\nfragment normalizedBodyModel_richText on RichText {\n  paragraphs {\n    ...normalizedBodyModel_richText_paragraphs\n    __typename\n  }\n  sections {\n    startIndex\n    ...getSectionEndIndex_section\n    __typename\n  }\n  ...getParagraphStyles_richText\n  ...getParagraphSpaces_richText\n  __typename\n}\n\nfragment normalizedBodyModel_richText_paragraphs on Paragraph {\n  markups {\n    ...normalizedBodyModel_richText_paragraphs_markups\n    __typename\n  }\n  codeBlockMetadata {\n    lang\n    mode\n    __typename\n  }\n  ...getParagraphHighlights_paragraph\n  ...getParagraphPrivateNotes_paragraph\n  __typename\n  id\n}\n\nfragment normalizedBodyModel_richText_paragraphs_markups on Markup {\n  type\n  __typename\n}\n\nfragment getParagraphHighlights_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment getParagraphPrivateNotes_paragraph on Paragraph {\n  name\n  __typename\n  id\n}\n\nfragment getSectionEndIndex_section on Section {\n  startIndex\n  __typename\n}\n\nfragment getParagraphStyles_richText on RichText {\n  paragraphs {\n    text\n    type\n    __typename\n  }\n  sections {\n    ...getSectionEndIndex_section\n    __typename\n  }\n  __typename\n}\n\nfragment getParagraphSpaces_richText on RichText {\n  paragraphs {\n    layout\n    metadata {\n      originalHeight\n      originalWidth\n      id\n      __typename\n    }\n    type\n    ...paragraphExtendsImageGrid_paragraph\n    __typename\n  }\n  ...getSeriesParagraphTopSpacings_richText\n  ...getPostParagraphTopSpacings_richText\n  __typename\n}\n\nfragment paragraphExtendsImageGrid_paragraph on Paragraph {\n  layout\n  type\n  __typename\n  id\n}\n\nfragment getSeriesParagraphTopSpacings_richText on RichText {\n  paragraphs {\n    id\n    __typename\n  }\n  sections {\n    ...getSectionEndIndex_section\n    __typename\n  }\n  __typename\n}\n\nfragment getPostParagraphTopSpacings_richText on RichText {\n  paragraphs {\n    type\n    layout\n    text\n    codeBlockMetadata {\n      lang\n      mode\n      __typename\n    }\n    __typename\n  }\n  sections {\n    ...getSectionEndIndex_section\n    __typename\n  }\n  __typename\n}\n\nfragment FriendLinkMeter_postContent on PostContent {\n  validatedShareKey\n  shareKeyCreator {\n    ...FriendLinkSharer_user\n    __typename\n    id\n  }\n  __typename\n}\n\nfragment FriendLinkSharer_user on User {\n  id\n  name\n  ...userUrl_user\n  __typename\n}\n\nfragment userUrl_user on User {\n  __typename\n  id\n  customDomainState {\n    live {\n      domain\n      __typename\n    }\n    __typename\n  }\n  hasSubdomain\n  username\n}\n"}`, req.Id))
	if err != nil {
		g.Log().Errorf(ctx, "get article content error. %s.", err)
		return nil, gerror.WrapCode(gcode.CodeInternalError, err, "get article content error.")
	}

	defer result.Close()
	//fmt.Println(result.ReadAllString())

	var acd *entity.ArticleContentData

	if err = gjson.DecodeTo(result.ReadAll(), &acd); err != nil {
		g.Log().Errorf(ctx, "scan  decodeTo ArticleContentData error.%s", err.Error())
		return nil, gerror.WrapCode(gcode.CodeNotImplemented, err, "scan  decodeTo ArticleContentData error.")
	}

	if ok, e := cache.Contains(cacheCtx, fmt.Sprintf("%s__data_%d", req.Id, req.Translate)); e == nil && ok {
		if rs, ee := cache.Get(cacheCtx, fmt.Sprintf("%s__data_%d", req.Id, req.Translate)); ee == nil {
			if ee = rs.Scan(&res); ee == nil && res != nil {
				return res, nil
			}
		}
	}

	res = &v1.ArticleGetContentRes{}
	res.ArticleContent = entity.ArticleContent{
		Id:         req.Id,
		Contents:   make([]entity.ArticleContentSub, 0),
		ZhContents: make([]entity.ArticleContentSub, 0),
	}

	for _, p := range acd.Data.Post.ViewerEdge.FullContent.BodyModel.Paragraphs {
		if p.Text != "" {
			res.ArticleContent.Contents = append(res.ArticleContent.Contents, entity.ArticleContentSub{
				Text:    p.Text,
				HtmlTag: p.Type,
			})
			if req.Translate == 1 {
				r, _ := s.Translate(ctx, p.Text)
				res.ArticleContent.ZhContents = append(res.ArticleContent.ZhContents, entity.ArticleContentSub{
					Text:    r,
					HtmlTag: p.Type,
				})
			}
		}
	}

	if ok, e := cache.Contains(cacheCtx, fmt.Sprintf("%s__data_%d", req.Id, req.Translate)); e == nil && !ok {
		cache.Set(cacheCtx, fmt.Sprintf("%s__data_%d", req.Id, req.Translate), res, time.Hour*24)
	}

	return res, err
}

func (s *sArticle) Translate(ctx context.Context, str string) (string, error) {
	url := "https://aip.baidubce.com/rpc/2.0/mt/texttrans/v1?access_token=" + GetAccessToken()
	payload := strings.NewReader(fmt.Sprintf(`{"from":"auto","to":"zh","q":"%s"}`, str))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		g.Log().Errorf(ctx, "Translate file error.%s", err.Error())
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		g.Log().Errorf(ctx, "Translate file error.%s", err.Error())
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		g.Log().Errorf(ctx, "Translate file error.%s", err.Error())
		return "", err
	}

	var dst entity.ArticleTranslate
	if err = gjson.DecodeTo(body, &dst); err != nil {
		g.Log().Errorf(ctx, "scan  decodeTo ArticleContentData error.%s", err.Error())
		return "", gerror.WrapCode(gcode.CodeNotImplemented, err, "scan Articletranslate error.")
	}

	ss := ""

	for _, d := range dst.Result.TransResult {
		ss = fmt.Sprintf(`%s%s`, ss, d.Dst)
	}

	return ss, nil
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", consts.API_KEY, consts.SECRET_KEY)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessTokenObj := map[string]string{}
	json.Unmarshal([]byte(body), &accessTokenObj)
	return accessTokenObj["access_token"]
}
