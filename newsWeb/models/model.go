package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Username string `orm；"unique"`
	Pwd      string
	Article []*Article`orm:"rel(m2m)"`
}
type Article struct {
	Id        int       `orm:"pk:aotu"`
	Title     string    `orm:"size(100)"`
	Content   string    `orm:"size(500)"`
	Time      time.Time `orm:"type(datetime);auto_now"`
	ReadCount int       `orm:"defailt（0)"`
	Image     string    `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	User []*User `orm:"reverse(many)"`
}
type ArticleType struct {
		Id int
		TypeName string `orm:"size(100)"`
		Article []*Article `orm:"reverse(many)"`
}

func init() {
	//生成三步走
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/newsweb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default", false, true)

}
