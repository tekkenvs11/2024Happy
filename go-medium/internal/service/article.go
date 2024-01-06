// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "go-medium/api/article/v1"
)

type (
	IArticle interface {
		ArticleGetList(ctx context.Context, req *v1.ArticleGetListReq) (res *v1.ArticleGetListRes, err error)
		ArticleGetContent(ctx context.Context, req *v1.ArticleGetContentReq) (res *v1.ArticleGetContentRes, err error)
	}
)

var (
	localArticle IArticle
)

func Article() IArticle {
	if localArticle == nil {
		panic("implement not found for interface IArticle, forgot register?")
	}
	return localArticle
}

func RegisterArticle(i IArticle) {
	localArticle = i
}
