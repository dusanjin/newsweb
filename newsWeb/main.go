package main

import (
	_ "newsWeb/routers"
	_ "newsWeb/models"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("PrePage", PrePage)
	beego.AddFuncMap("NextPage", NextPage)
	beego.Run()
}

func PrePage(pageIndex int) int {

	res := pageIndex - 1
	if res < 1 {
		res = 1
	}
	return res
}
func NextPage(pageIndex int, page float64) int {
	if pageIndex+1 > int(page) {
		return pageIndex
	}

	return pageIndex + 1
}
