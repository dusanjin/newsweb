package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/articles/*",beego.BeforeExec, funcfil)
    beego.Router("/register", &controllers.UserController{},"get:Get;post:ShowHandle")
	beego.Router("/login", &controllers.UserController{},"get:LoginGet;post:ShowLogin")
	beego.Router("/article", &controllers.ArticleController{},"get:ArtrollersGet")
	beego.Router("/AddArticleController", &controllers.ArticleController{},"get:AddArticleController;post:Showadd")
	beego.Router("/ContentcleController", &controllers.ArticleController{},"get:ContentGet")
	beego.Router("/updeta",&controllers.ArticleController{},"get:Updateaget;post:ShowUp")
	beego.Router("/delet",&controllers.ArticleController{},"get:Showdelet")
	beego.Router("/addtype",&controllers.ArticleController{},"get:Addtypeget;post:Showaddtype")
	beego.Router("/delettype",&controllers.ArticleController{},"get:ShowDelet")
	beego.Router("/igout",&controllers.ArticleController{},"get:Showuot")
}
func funcfil(ctx *context.Context)  {
	UserName:=ctx.Input.Session("UserName")
	if UserName.(string)=="" {
		ctx.Redirect(302,"/login")
		return
	}
}