package controllers

import (
	"github.com/astaxie/beego"
	"newsWeb/models"
	"github.com/astaxie/beego/orm"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	this.TplName = "register.html"
}
func (this *UserController) LoginGet() {
	//判断是否有缓存后台缓存的客户账号
	bas := this.Ctx.GetCookie("userName")
	userName,err:=base64.StdEncoding.DecodeString(bas)
	if err!=nil {
		this.Data["errmsg"]=err
		this.TplName = "login.html"
		return
	}
	if string(userName)!="" {
		this.Data["userName"]=string(userName)
		this.Data["checked"]="checked"
	}else {
		this.Data["userName"]=""
		this.Data["checked"]=""
	}
	this.TplName = "login.html"
}

func (this *UserController) ShowHandle() {
	username := this.GetString("userName")
	pwd := this.GetString("password")
	if username == "" || pwd == "" {
		beego.Error("账号密码不能为空")
		this.TplName = "resister.html"
		return
	}
	o := orm.NewOrm()
	var cl models.User
	cl.Username = username
	cl.Pwd = pwd
	o.Insert(&cl)
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login", 302)
}
func (this *UserController) ShowLogin() {
	username := this.GetString("userName")
	pwd := this.GetString("password")
	if username == "" || pwd == "" {
		this.Data["errMsg"] = "用户密码不能为空"
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var cl models.User
	cl.Username = username
	err := o.Read(&cl, "username")
	if err != nil {
		this.Data["errMsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}
	if pwd != cl.Pwd {
		this.Data["errMsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}
	remember := this.GetString("remember")
	//判断是由有勾选
	if remember == "on" {
		bas:=base64.StdEncoding.EncodeToString([]byte(username))
		//判定结果后传送给客户端，缓存在客户端。
		this.Ctx.SetCookie("userName", bas, 36000*1)
	}else {
		this.Ctx.SetCookie("userName","",-1)
	}
	//this.Ctx.WriteString("登陆成功")
	//登陆成功后向列表页面冲宋服务器username。由后面的列表页面进行判定
	this.SetSession("userName",username)
	this.Redirect("/article", 302)
}
