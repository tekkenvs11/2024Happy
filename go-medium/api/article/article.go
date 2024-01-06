// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package article

import (
	"context"

	"go-medium/api/article/v1"
)

type IArticleV1 interface {
	ArticleGetList(ctx context.Context, req *v1.ArticleGetListReq) (res *v1.ArticleGetListRes, err error)
	ArticleGetContent(ctx context.Context, req *v1.ArticleGetContentReq) (res *v1.ArticleGetContentRes, err error)
}
