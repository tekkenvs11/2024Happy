package packed

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBgYGDg5FOLYEAC/AycDLmJeZlpqcUl+sn5eWmZ6aEhrAyMyRtnpAZ4s3Mgq4WZwoFhijSmKVBKrzIxNwduoklPIBeTIY/L/+IuRfYLdV5+D0WZeF973Zi2JHq9vmj5vrSKM97dmYosGWsTZzaJTnPXbOXMuMv39nfwce8D5W9/v537NkLDfleumPgvf7VP205N2uVbzBFxt3564G1LFmU/ZaVfbAwM//+D3B2/kq8vjYGBoR7J3QwMLa6o7uZAcjfYkTeEZmLxNiOTCDPu4IMAAYb/jiAad2DCTIEEH7IDpeGmMDBsa3TDNAVrYCK7C5vHEO76h2QikjdZ2UAKmBmYGa4zMDDYMoJ4gAAAAP//yzKACyoCAAA="); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
