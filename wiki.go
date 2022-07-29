package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"text/template"
)

type Page struct {
	Title string //标题
	Body  []byte //正文 类型为[]byte类型 而不是string，为了方便存储
}

//保存Page到一个文本文件
func (p *Page) save() error {
	filename := p.Title + ".txt"
	//保存到下一级data目录
	dist := path.Join("./data", filename)
	return os.WriteFile(dist, p.Body, 0600) //0600让当前用户拥有读写权限
}

//将保存的Page读取出来
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	//从data目录中读取txt文件
	dat := path.Join("./data", filename)
	body, err := os.ReadFile(dat)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//查看方法
func viewHandler(w http.ResponseWriter, r *http.Request) {
	//记录访问的链接
	log.Println("viewHandler:", r.URL.Path)
	//获取了/view/后边的内容作为标题
	//title := r.URL.Path[len("/view/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		//没有该内容，则重定向到edit界面编辑该内容
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	/*----原代码
	//读取view模板
	t, err2 := template.ParseFiles("./tmpl/view.html")
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	//解析view模版
	err2 = t.Execute(w, p)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
	*/
	renderTemplate(w, "./tmpl/view", p)
}

//编辑方法
func editHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/edit/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//t, err3 := template.ParseFiles("edit.html")
	//if err3 != nil {
	//	http.Error(w, err3.Error(), http.StatusInternalServerError)
	//	return
	//}
	//err3 = t.Execute(w, p)
	//if err3 != nil {
	//	http.Error(w, err3.Error(), http.StatusInternalServerError)
	//	return
	//}
	renderTemplate(w, "./tmpl/edit", p)
}

//viewHandler和editHandler有很多相同的代码，定义一个函数重构它们，优化代码
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//保存方法
func saveHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/save/"):]
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "Success，已将 %s 保存!", title)
	//保存好后重定向到view界面
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//正则表达式
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//获取title，验证输入的内容
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // 提取title  在下标2 (localhost:8000/xxx/title)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
