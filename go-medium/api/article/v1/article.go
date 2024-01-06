package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-medium/internal/model/entity"
)

type ArticleGetListReq struct {
	g.Meta   `path:"/article/list" tags:"article" method:"get" summary:"List"`
	Year     int    `json:"year" dc:"year"`
	Host     string `json:"host" dc:"socks5 host"`
	Port     int    `json:"port" dc:"socks5 port"`
	UserName string `json:"UserName" dc:"socks5 userName"`
	Password string `json:"password" dc:"socks5 password"`
}

type ArticleGetListRes struct {
	Lists []entity.Article `json:"data" dc:"article list"`
}

type ArticleGetContentReq struct {
	g.Meta    `path:"/article/content" tags:"article" method:"get" summary:"Content"`
	Id        string `json:"id" dc:"article id" v:"required"`
	Translate int    `json:"translate" dc:"translate id"`
	Host      string `json:"host" dc:"socks5 host" d:"101.36.113.108"`
	Port      int    `json:"port" dc:"socks5 port" d:"10801"`
	UserName  string `json:"UserName" dc:"socks5 userName" d:"appproxy"`
	Password  string `json:"password" dc:"socks5 password" d:"662F41D99defM377"`
}

type ArticleGetContentRes struct {
	entity.ArticleContent `json:"content" dc:"content"`
}
