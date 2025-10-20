package main

import (
	"fmt"

	"github.com/wonderfulsuccess/go-web-app/back/utils/crawler"
)

const cURL = `
curl 'https://sf.taobao.com/list/50025969.htm?spm=a213w.7398504.pagination.1.181e3a1bhiX6ck&auction_source=0&st_param=-1&auction_start_seg=-1&page=2' \
  -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
  -H 'accept-language: zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7' \
  -b 'arms_uid=97f1e648-da22-4b43-86de-6f943318ee71; cna=3B8PHd/7vXMCAXhVqzUAc7cs; tracknick=%5Cu5B89%5Cu59AE%5Cu4F60%5Cu6765%5Cu4E00%5Cu4E0B; _hvn_lgc_=0; wk_unb=VvlxdApZR0sq; lgc=%5Cu5B89%5Cu59AE%5Cu4F60%5Cu6765%5Cu4E00%5Cu4E0B; wk_cookie2=1639f5aa54e0c7bca29df7864429cd0c; dnk=%5Cu5B89%5Cu59AE%5Cu4F60%5Cu6765%5Cu4E00%5Cu4E0B; thw=cn; havana_lgc2_0=eyJoaWQiOjU3NjUyMjkxMCwic2ciOiI4ZjE3YzBhM2VkZjdiMDhhNjM2MjlhM2MzNDNiMWFiNCIsInNpdGUiOjAsInRva2VuIjoiMUhEVi16NEpjaW9HVkFhVVJ2ZjQ3UmcifQ; hng=CN%7Czh-CN%7CCNY%7C156; t=10591c9b6612f080a4a7ecc4636a2d1a; cnaui=576522910; aui=576522910; sn=; cancelledSubSites=empty; _tb_token_=ee05e3141f37f; cookie2=1951109d88bcecf9ae9b139d232ff39a; miid=1237445372379309768; useNativeIM=false; wwUserTip=false; _uetvid=9da5c3a084a011efba454fa816dbe63b; xlly_s=1; mtop_partitioned_detect=1; _m_h5_tk=4319d5aade90cc6d99789658cf6502ea_1760723863228; _m_h5_tk_enc=af636290c5477734353d52791589f003; sca=96e55ae6; _samesite_flag_=true; unb=576522910; uc3=lg2=V32FPkk%2Fw0dUvg%3D%3D&nk2=0%2FJGX6oixYstTzcZ&vt3=F8dD2ky%2FIqPPXzg6Vug%3D&id2=VvlxdApZR0sq; csg=862f07f1; cookie17=VvlxdApZR0sq; skt=895f50d793dafeda; existShop=MTc2MDcxMzc4NA%3D%3D; uc4=nk4=0%400Un2aljs0y6kvH%2FKs10JRusyuv2829A%3D&id4=0%40VH95fxVO5RZkLrICdMa1Jc44X10%3D; _cc_=W5iHLLyFfA%3D%3D; _l_g_=Ug%3D%3D; sg=%E4%B8%8B0e; _nk_=%5Cu5B89%5Cu59AE%5Cu4F60%5Cu6765%5Cu4E00%5Cu4E0B; cookie1=VqoLGUhPY5%2BZbqVpddi7JnXkcQgYfQEaN6uZX6KXUGs%3D; sgcookie=E100KkYt8XLze%2BaCf9HQJyTX4Qd9jSMrDOswuK%2FLvEnFWtUqTQvCfsXEF8%2Fmx0AglzxwpJIqrW28JehpbnOAst3LBqllL2nKrcfj4aKEKsOzdq3U%2Bam9hpNCh7%2BSp4uvhDjv; havana_lgc_exp=1791817784821; sdkSilent=1760742584821; havana_sdkSilent=1760742584821; isg=BE1NmA70Gd4nMbHYhP-mv9ghXG_HKoH8MqLwkY_Sj-RThm04V3jfzOwH8hrgRpm0; uc1=cookie14=UoYY4%2FxQx2pGTg%3D%3D&cookie15=WqG3DMC9VAQiUQ%3D%3D&pas=0&cookie21=UIHiLt3xTIkz&existShop=false&cookie16=WqG3DMC9UpAPBHGz5QBErFxlCA%3D%3D; tfstk=gqQSXI0n49QV5Pj9VePqchl9-MLQgSzaFXOds63r9ULJpD9f67vzxuXCc9BD4LIyL2_CDtCrzp5KOBB2yHjyLLJBR6X_gRza7_fk-UyaQPzV0L0wS03-yvJxMBYBw0MP0l1k-eeqg2Ea5_XoTyxHwBdYGBRE2bBK2qUvnBR-e9pJDxd6sepdp9pvkBdE2bdpeqLvnBLpp9KpkmpDOeddpeFfMKAKOoZkgX9cN5dUzqBrFF1Jh23dyyxWCPvviI7j73p1watJVZwMVdCJh2e-amn9HC_YKYJhlsQJ_tUqPe1vOOTRc-35d_IGKCB8lYd5211wcw2tRIjAEnR5cWnve1pR6U7EH79hzT_BqNeZoLIfowTP08gwKHIVrhbzH4TOj_8PvtFKJpsv9guE7dOFQDGXj2OXQSNjxDvI-LBiSJqj03dDgXFbGcuk2IAXJSNjxDxJiIraGSiw2' \
  -H 'priority: u=0, i' \
  -H 'referer: https://sf.taobao.com/item_list.htm?spm=a213w.3064813.a214dqe.3.7e8e3fe7dwuO0m&category=50025969' \
  -H 'sec-ch-ua: "Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  -H 'sec-fetch-dest: document' \
  -H 'sec-fetch-mode: navigate' \
  -H 'sec-fetch-site: same-origin' \
  -H 'sec-fetch-user: ?1' \
  -H 'upgrade-insecure-requests: 1' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36'
`

func main() {
	result, err := crawler.RunCURL(cURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(result)
}
