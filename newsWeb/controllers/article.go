package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/orm"
	"path"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"time"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

//ArtrollersGet 列表页面
func (this *ArticleController) ArtrollersGet() {
	errmsg := this.GetString("errnsg")
	if errmsg != "" {
		this.Data["errmsg"] = errmsg
	}
	//判断是否有登陆，抓取从服务区是否有传送username
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	//传送 过来的为接口类型，转换为string类型
	this.Data["userName"] = userName.(string)
	o := orm.NewOrm()
	var article []models.Article
	//讲数据库中areicle列表
	qs := o.QueryTable("article")
	//qs.All(&article)  这种是将所有文件信息存入到切片中
	pag := 2 //判定呈现的数据个数
	//判定传送过来的id数据
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	start := pag * (int(pageIndex) - 1)
	//是否有区分列表，当使用get请求访问是select为空
	//首次请求取到的select为空，在此通过更改列表再次访问会传送一个select的值
	selects := this.GetString("select")
	this.Data["selects"] = selects
	beego.Info(selects)
	var pages float64
	if selects == ""||selects=="请选择"{
		count, _ := qs.Count()
		pages = float64(count) / float64(pag)
		pages = math.Ceil(pages)
		//relateasel 一对多的关系表查询中，用来制定另外一张表的函数，关联后在写入到哦article切片中
		qs.Limit(pag, start).RelatedSel("ArticleType").All(&article)
		this.Data["count"] = count
	}else {
		count, _ := qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", selects).Count()
		//relateasel 一对多的关系表查询中，用来制定另外一张表的函数
		qs.Limit(pag, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", selects).All(&article)
		pages = float64(count) / float64(pag)
		pages = math.Ceil(pages)
		this.Data["count"] = count
	}
	this.Data["page"] = pages

	//用来让循环两张表循环关联起来
	//var cf []models.Article
	//for _,v:=range article{
	//	id:=v.ArticleType.Id
	//	var cd  models.ArticleType
	//	cd.Id=id
	//	err=o.Read(&cd)
	//	if err != nil {
	//		this.Data["errmsg"] = "id读取文件错误"
	//		this.TplName = "add.html"
	//		return
	//	}
	//	v.ArticleType=&cd
	//	cf=append(cf,v)
	//}
	qa := o.QueryTable("ArticleType")
	var ck []models.ArticleType

	qa.All(&ck)
	this.Data["ck"] = ck
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = article

	this.TplName = "index.html"
}

//添加列表页面get请求

func (this *ArticleController) AddArticleController() {
	errmsg := this.GetString("errmsg")
	this.Data["errmsg"] = errmsg
	o := orm.NewOrm()
	qa := o.QueryTable("ArticleType")
	var ck []models.ArticleType
	qa.All(&ck)
	this.Data["ck"] = ck
	this.TplName = "add.html"
}

//添加列表页面post请求
func (this *ArticleController) Showadd() {
	arename := this.GetString("articleName")
	content := this.GetString("content")
	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	if arename == "" || content == "" {
		errmsg := "文件不能为空"
		this.Redirect("/AddArticleController?errmsg="+errmsg, 302)
		return
	}
	if err != nil {
		errmsg := "文件添加失败"
		this.Redirect("/AddArticleControlle?errmsg="+errmsg, 302)
		return
	}
	if head.Size > 50000 {

		errmsg := "文件过大，请重新上传"
		this.Redirect("/AddArticleController?errmsg="+errmsg, 302)
		return
	}
	fixfname := path.Ext(head.Filename)
	if fixfname != ".jpg" && fixfname != ".png" && fixfname != ".jpeg" {
		errmsg := "文件格式错误"
		this.Redirect("/AddArticleController?errmsg="+errmsg, 302)
		return
	}

	fillname := time.Now().Format("2006-01-02-15-04-05") + fixfname
	this.SaveToFile("uploadname", "./static/image/"+fillname)
	o := orm.NewOrm()
	selects := this.GetString("select")
	var cd models.ArticleType
	cd.TypeName = selects
	err = o.Read(&cd, "TypeName")
	if err != nil {
		this.Data["errmsg"] = "id读取文件错误"
		this.TplName = "add.html"
		return
	}
	var cl models.Article
	cl.ArticleType = &cd
	cl.Title = arename
	cl.Content = content
	cl.Image = "/static/image/" + fillname
	o.Insert(&cl)
	this.Redirect("/article", 302)

}

//查看文件内容页面
func (this *ArticleController) ContentGet() {
	o := orm.NewOrm()
	var cl models.Article
	pageIndde, err := this.GetInt("pageIndde")
	if err != nil {
		this.Redirect("/article", 302)
		return
	}
	cl.Id = pageIndde
	err = o.Read(&cl)
	if err != nil {
		this.Redirect("/article", 302)
		return
	}
	var ck models.ArticleType
	id := cl.ArticleType.Id
	ck.Id = id
	err = o.Read(&ck)
	if err != nil {
		this.Redirect("/article", 302)
		return
	}
	cl.ArticleType = &ck
	this.Data["cl"] = cl
	//第一部 获取需要插入的文章个
	//获取文章一对多关系  获取需要插入的文章的那个字段
	m2m := o.QueryM2M(&cl, "User")
	//获取user对象  创建一个需要插入的对象
	var user models.User

	//获取传送过来的username
	userName := this.GetSession("userName")
	user.Username = userName.(string)
	//赋值
	o.Read(&user, "UserName")
	//插入多对多关系  插入对象
	m2m.Add(user)
	//第一种多对多查询  查询
	o.LoadRelated(&cl, "User")
	n := len(cl.User)

	var users []models.User
	//select * from user where
	o.QueryTable("User").Filter("Article__Article__Id", pageIndde).Distinct().All(&users)
	cl.ReadCount = n
	o.Update(&cl)
	this.Data["len"] = n
	this.Data["Username"] = user.Username
	this.Data["users"] = users
	this.TplName = "content.html"
}

//编辑文章内容get方法
func (this *ArticleController) Updateaget() {
	articlid, err := this.GetInt("Id")
	errmsg := this.GetString("errmsg")
	if errmsg != "" {
		this.Data["errmsg"] = errmsg

	}
	if err != nil {
		errmsg := "请求路径错误"
		beego.Error("请求路径错误")
		this.Redirect("/article?errmsg="+errmsg, 302)
	}
	o := orm.NewOrm()
	var cl models.Article
	cl.Id = articlid
	err = o.Read(&cl)
	if err != nil {
		errmsg := "请求路径错误"
		beego.Error("请求路径错误")
		this.Redirect("/article?errmsg="+errmsg, 302)
	}
	this.Data["cl"] = cl

	this.TplName = "update.html"
}

//获取图片存储u路径函数
func up(this *ArticleController, update string, err1 string) string {
	file, head, err := this.GetFile(update)
	defer file.Close()
	if err != nil {
		beego.Error(1)
		errmsg1 := "文件添加失败"
		this.TplName = err1 + errmsg1
		return ""
	}
	if head.Size > 50000 {
		beego.Error(2)
		errmsg := "文件过大，请重新上传"
		this.TplName = err1 + errmsg
		return ""
	}
	fixfname := path.Ext(head.Filename)
	if fixfname != ".jpg" && fixfname != ".png" && fixfname != ".jpeg" {
		beego.Error(3)
		errmsg2 := "文件格式错误"
		this.TplName = err1 + errmsg2
		return ""
	}

	fillname := time.Now().Format("2006-01-02-15-04-05") + fixfname
	this.SaveToFile(update, "./static/image/"+fillname)
	return "/static/image/" + fillname
}

//编辑文章内容post方法
func (this *ArticleController) ShowUp() {
	arename := this.GetString("articleName")
	content := this.GetString("content")
	id, err := this.GetInt("Id")
	articleType := this.GetString("select")
	updete := up(this, "uploadname", "/updeta?Id="+strconv.Itoa(id)+"&errmsg=")

	beego.Info(updete)

	if arename == "" || content == "" || err != nil || updete == "" {
		errmsg := "传输有无"
		this.Redirect("/updeta?Id="+strconv.Itoa(id)+"&errmsg="+errmsg, 302)
		return
	}

	o := orm.NewOrm()
	var cl models.Article
	cl.Id = id
	err = o.Read(&cl)
	if err != nil {
		errmsg := "请求失败"
		this.Redirect("/updeta?Id="+strconv.Itoa(id)+"&errmsg"+errmsg, 302)
		return
	}

	cl.ArticleType.TypeName = articleType
	cl.Image = updete
	cl.Content = content
	cl.Title = arename
	o.Update(&cl)
	this.Redirect("/article", 302)

}

//文章内容删除get方法
func (this *ArticleController) Showdelet() {
	id, err := this.GetInt("Id")
	if err != nil {
		errmsg := "获取id失败"
		this.Redirect("/article?id="+strconv.Itoa(id)+"&errmsg"+errmsg, 302)
		return
	}
	o := orm.NewOrm()
	var cl models.Article
	cl.Id = id
	err = o.Read(&cl)
	if err != nil {
		errmsg := "无此文件"
		this.Redirect("/article?id="+strconv.Itoa(id)+"&errmsg"+errmsg, 302)
		return
	}
	_, err = o.Delete(&cl)
	if err != nil {
		errmsg := "删除文件失败"
		this.Redirect("/article?id="+strconv.Itoa(id)+"&errmsg"+errmsg, 302)
		return
	}
	this.Redirect("/article", 302)
}

//添加文章类型get方法
func (this *ArticleController) Addtypeget() {
	o := orm.NewOrm()
	qs := o.QueryTable("ArticleType")
	var addtypeget []models.ArticleType
	_, err := qs.All(&addtypeget)

	if err != nil {
		this.Data["errmsg"] = "切片获取失败"
		this.TplName = "addType"
		return
	}
	this.Data["addtypeget"] = addtypeget
	//this.Data["addtypeget"]=addtypeget
	this.TplName = "addType.html"
}

//添加文章pos方法
func (this *ArticleController) Showaddtype() {
	typename := this.GetString("typeName")
	o := orm.NewOrm()
	var cl models.ArticleType
	cl.TypeName = typename
	o.Insert(&cl)
	this.Redirect("/addtype", 302)
}

//删除文章类型get方法
func (this *ArticleController) ShowDelet() {
	id, err := this.GetInt("Id")
	if err != nil {
		errmsg := "id获取失败"
		this.Redirect("/addtype?errmsg="+errmsg, 302)
		return
	}
	o := orm.NewOrm()
	var cl models.ArticleType
	cl.Id = id
	_, err = o.Delete(&cl)
	if err != nil {
		errmsg := "删除失败"
		this.Redirect("/addtype?errmsg="+errmsg, 302)
		return
	}
	this.Redirect("/addtype", 302)
}

//退出用户
func (this *ArticleController) Showuot() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}
