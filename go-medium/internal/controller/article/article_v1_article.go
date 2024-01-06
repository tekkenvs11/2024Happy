package article

import (
	"context"
	"go-medium/internal/service"

	"go-medium/api/article/v1"
)

func (c *ControllerV1) ArticleGetList(ctx context.Context, req *v1.ArticleGetListReq) (res *v1.ArticleGetListRes, err error) {
	return service.Article().ArticleGetList(ctx, req)
}
func (c *ControllerV1) ArticleGetContent(ctx context.Context, req *v1.ArticleGetContentReq) (res *v1.ArticleGetContentRes, err error) {
	return service.Article().ArticleGetContent(ctx, req)
}
