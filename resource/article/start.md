# go-web入门

* 处理请求
* 模板
* 中间件
* 存储数据
* https/http2
* 测试
* 部署

## go mod...

之前的go做的东西都用的GOPATH模式，现代go语言更推荐go module，这里复习一下把它用起来。

```shell
# 查看help
go mod help

# 开启go mod
go env -w GO111MODULE=on

# 设置代理
go env -w GOPROXY
```

```shell
# 初始化gomod模块
go mod init [模块名]

go mod init github.com/ty/mod_test

# 下载第三方库 默认最新版本
go get [url]
```

## 1 net.http编写第一个demo

```go
package main

import "net/http"

func main() {
    // 注册一个函数到路由器
    http.HandleFunc("/", func(w http.responseWriter, r *http.request) {
        w.Write([]byte("net.http demo"))
    })

    // 启动服务，使用默认路由器
    http.ListenAndServe("localhost:8080", nil)
}
```

## 2 处理(Handle)请求

### 创建Web Server

* `http.ListenAndServe()`
  * 参数1表示网络地址，如果是空，则表示所有网络接口的80端口
  * 参数2是handler，如果是`nil`，就表示`DefaultServeMux`

* `DefaultServeMux`是一个`multiplexer`(多路复用器，类似路由器)

`http.ListenAndServe()`函数体中创建了一个`http.Server`并启动:

* `http.Server`是一个`struct`
  * Addr，Handler对应上面函数的两个参数
  * ListenAndServe()成员函数

```go
// 创建web server的另一种方法
s := http.Server {
    Addr: "192.168.18.128",
    Handler: nil
}
s.ListenAndServe()
```

### Handler

* 源码定义
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

### DefaultServeMux

`DefaultServeMux`实现了`Handler`接口，作用是*http路由*;

即解析http请求的uri，然后将请求转发至响应的handler。

据此，可以自定义实现一个Handler：

```go
type myHandler struct {}

func (m *myHandler) Handler(w http.ResponseWriter, r *Request) {

}
```

### 多个Handler

* 不指定`http.Server`中的`Handler`字段，使用默认的多路复用器
* 使用`http.Handle()`将某个`Handler`注册到`DefaultServeMux`上

```go
// 注册Handler方法
// 1
func Handle(pattern string, handler Handler)

// 2
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
/*
使用 http.HandlerFunc() 将handler函数转化为一个实现http.Handler接口的对象
然后在调用函数1注册路由, 相当于 Handler(pattern, HandlerFunc(handler))
*/
```

注意，`HandlerFunc`并不是一个函数，而是一个自定义类型，该类型实现了`Handler`接口：

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHttp(w ResponseWriter, r *Request) {
    f(w, r)
}
```

## 3 内置的Handler

### `http.NotFoundHandler`

```go
func NotFoundHandler() Handler
```

返回一个handler，它总是响应"404, page not found"

### `http.RedirectHandler`

```go
func RedirectHandler(url string, code int) Handler
```

返回一个hander，它用给定状态码跳转至指定URL

### `http.StripPrefix`

```go
func http.StripPrefix(prefix string, h Handler) Handler
```

返回一个handler，它从请求URL中取出指定前缀，并调用指定handler；

如果请求url中与提供前缀不符合，那么响应404；

略像中间件，修饰了另一个handler

### `http.TimeoutHandler`

```go
func TimeoutHandler(h Handler, dt time.Duration, msg string)
```

返回一个handler，它用来在指定时间内执行指定的Handler；

相当于一个修饰器；

msg表示如果超时，会将msg作为响应正文返回。

### `http.FileServer`

```go
func FileServer(root FileSystem) Handler
```

返回一个handler，使用基于root的文件系统来响应请求，其中：

```go
type FileSystem interface {
    open(name string) (File, error)
}
```

使用时需要用到操作系统的文件系统，所以还需要委托给：
```go
type Dir string
func (d Dir) Open(name string) (File, error)
```

***

```go

```

## 4 请求Request

* HTTP请求
* Request
* URL
* Header
* Body

### http消息(Request/Response)

### 请求Request

Request是个`struct`，代表客户端发送的HTTP请求消息

* 重要字段：
  * URL
  * Header
  * Body
  * Form、PostFrom、MultipartForm

也可以通过其中的方法获取请求中的Cookie、URL、User Agent等信息

#### URL

* 通用形式：`scheme://[userinfo@]host/path[?query][fragment]`
* 不以斜杠开头的URL：`scheme:opaque[?query][fragment]`

* **URL Query**

例URL: http://example.org/post?id=123&thread_id=456

`r.URL.RawQuery`可以获取实际查询的原始字符串，实例中即`id=123&thread_id=456`

`r.URL.Query()`会提共查询字符串对应的`map[string][]string`
```go
query := r.URL.Query() // map

id := query["id"]
name := query.Get("name")
```

* **URL fragment**

如果是从浏览器发出的请求，那么无法提取出`fragment`字段(浏览器发送请求时会将该部分去掉)

但不是所有请求都是从浏览器发出的。

#### Header

是一个map，key是`string`，value是`[]string`

```go
// 获取map
r.Header

// 通过键获取值
var vals []string = r.Header["accept-encoding"] 

var str = r.Header.Get("accept-encoding") // vals[0]
```

#### Body

请求和响应的bodies都是使用`Body`字段来表示的；

`Body`是一个`io.ReadCloser`接口，即实现了`Reader`接口和`Closer`接口：
```go
type ReadCloser interface {
	Reader // Read([]byte) (int, error)
    Closer // Close() error
}
```

```go
// 读取消息体内容 r *http.Request
bodybuf := make([]byte, r.ContentLength)
r.Body.Read(bodybuf)
fmt.Println(string(bodybuf))
```

## 5 Form

* 通过表单发送请求
* Form字段
* PostForm字段
* MutipartForm字段
* FromValue&PostFormValue方法
* 文件Files
* POST JSON

### 来自表单的post请求

```html
<form action = "/process" method = "post">
    <input type = "text", name = "name"/>
</form>
```

示例html表单里面的数据以name-value对的形式，通过post请求发送出去；

它的数据会放在post请求的Body中；

***

通过post请求发送的name-value数据对的格式可以通过表单的Content Type来指定，即`enctype`：

### 表单的`enctype`属性

默认值：`application/x-www-form-urlencoded`

浏览器被要求至少要支持：`application/x-www-form-urlencoded` `mutipart/form-data`

HTML5要需要支持`text/plain`

* 如果`enctype`是`application/x-www-form-urlencoded`，那么浏览器会将表单的数据编码到查询字符串中；(简单文本)
* 如果`enctype`是`mutipart/form-data` (大量数据)：
  * 每一个name-value对都会被转换为一个MIME消息部分
  * 每一个部分都有自己的Content Type和Content Disposition

### 表单的GET

通过表单的method属性，可以设置post或get；

get请求没有Body，所有的数据都通过URL的name-value对发送

### Form字段

* `http.Request`上的函数允许我们从URL或Body中提取数据，用到的字段：
  * `Form`
  * `PostForm`
  * `MultipartForm`

* `Form`中的数据是*key-value*对
* 通常的用法是：
  * 先调用`ParseForm()`或`ParseMultipartForm`来解析相应字段
  * 然后通过字段值访问即可

```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintln(w, r.Form)
})
```

### PostForm字段

如果表单和URL中有同样的key，那么它们的的值会被放到同一个slice中，表单的值靠前，URL的值靠后

如果只想要表单的key-value对，可以使用`PostForm`字段

```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    // 只会回写表单的数据
    fmt.Fprintln(w, r.PostForm)
})
```

***

Form和PostForm只支持`application/x-www-form-urlencoded`

想要得到multipart key-value对，必须使用`MultipartForm`字段

### MultipartForm字段

* 想要使用这个字段，首先要调用`ParseMultipartForm()`
  * 该方法在必要时调用`ParseForm()`
  * 参数是需要读取的数据长度

* `MultipartForm`只包含表单的*key-value*对
* 返回类型是一个`struct`，包含两个`map`
  1. key->string, value->[]string
  2. key->string, value->file(用于上传文件)


```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(1024)
    fmt.Fprintln(w, r.MultipartForm)
})
```

### FormValue() 和 PostFormValue()

* `FromValue()`会返回Form字段中指定的key对应的第一个value
  * 无需调用`ParseForm()`等，会自动调用

```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, r.FormValue("name"))
})
```

* `PostFormValue()`也类似，但只会读取`PostForm`
  * 上述两个方法都会调用`ParseMultipartForm()`方法

### 上传文件

1. 首先调用`ParseMutipartForm()`方法
2. 通过`MultipartForm`字段的`File`成员获取文件句柄
3. 打开文件读写

```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(1024)
    
    fileHeader := r.MultipartForm.File["uploaded"][0]
    file, err := fileHeader.Open()
    if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
            fmt.Fprintln(w, string(data))
        }
    }
})
```

上述函数可使用`FormFile()`方法简化：

```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    // r.ParseMultipartForm(1024)
    
    // fileHeader := r.MultipartForm.File["uploaded"][0]
    // file, err := fileHeader.Open()

    file, fileHeader, err := r.FormFile("uploaded")
    if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
            fmt.Fprintln(w, string(data))
        }
    }
})
```

如果只上传了一个文件，这种方式会快一些

### POST请求 - Json Body

* 不是所有的Post请求都来自Form
* 客户端框架会以不同的方式对Post请求进行编码
  * jQuery通常使用`application/x-www-form-urlencoded`
  * Angular是`application/json`
* `ParseForm()`方法无法处理`application/json`

### MutipartReader()方法

```go
func (r *request) MultipartReader() (*mutltipart.Reader, error)
```

如果请求是multipart/form-data POST请求，MultipartReader返回一个multipart.Reader接口，否则返回nil和一个错误。

使用本函数代替ParseMultipartForm，可以将r.Body作为流stream处理。


## 6 ResponseWriter

* 从服务器向客户端返回相应要使用`ResponseWriter`
* `ResponseWriter`是一个接口，handler用它返回响应
* 真正支持`ResponseWriter`的`struct`是非导出的`http.response`

```go
type ResponseWriter interface {
    // Header返回一个Header类型值，该值会被WriteHeader方法发送。
    // 在调用WriteHeader或Write方法后再改变该对象是没有意义的。
    Header() Header
    
    // WriteHeader该方法发送HTTP回复的头域和状态码。
    // 如果没有被显式调用，第一次调用Write时会触发隐式调用WriteHeader(http.StatusOK)
    // WriterHeader的显式调用主要用于发送错误码。
    WriteHeader(int)
    
    // Write向连接中写入作为HTTP的一部分回复的数据。
    // 如果被调用时还未调用WriteHeader，本方法会先调用WriteHeader(http.StatusOK)
    // 如果Header中没有"Content-Type"键，
    // 本方法会使用包函数DetectContentType检查数据的前512字节，将返回值作为该键的值。
    Write([]byte) (int, error)
}
```

### 内置Response响应

* `NotFound()` 包装一个404状态码和一个额外信息
* `ServeFile()` 从文件系统提供文件
* `ServeContent()` 提供内容，支持分片
* `Redirect()` 重定向至另一个url


## 7 模板

### 模板

Web模板就是预先设计好的HTML页面，它可以被*模板引擎template engine*反复使用，来产生HTML页面；

Go标准库提供了`text/template`和`html/template`两个模板库；

大多数Go的web框架都使用这些库作为默认的模板引擎

### 模板与模板引擎

模板引擎可以合并模板和上下文数据，产生最终的HTML。

* **无逻辑模板引擎**
  * 通过占位符，动态数据被替换到模板中
  * 不做任何逻辑处理，只做字符串替换
  * 处理完全由handler完成
  * 目标是展示层和逻辑的完全分离

* **逻辑嵌入模板引擎**
  * 编程语言被嵌入到模板中
  * 在运行时由模板引擎来执行，也包含替换功能
  * 功能强大
  * 逻辑代码遍布handler和模板，难以维护

### Go的模板引擎

* 主要使用的是`text/template`，html相关部分使用了`html/template`
* 模板可以完全无逻辑，但有足够的嵌入特性
* 和大多数模板引擎一样，Go Web的模板位于无逻辑和嵌入逻辑之间

* **工作原理**
  * 在web应用中，通常是由handler触发模板引擎
  * handler调用模板引擎，并将使用的模板(通常是一组模板文件和动态数据)传递给引擎
  * 模板引擎生成HTML，写入到`ResponseWriter`中


### 关于模板

* 模板必须是可读的*文本文件*，*扩展名任意*；对web应用通常是html
  * 里面会内嵌一些命令(称为*action*)
* `text/template`是通用模板引擎，`html/template`是html模板引擎
* action位于双层花括号之间：`{{.}}`
  * 这里的`.`就是一个action
  * 它可以命令模板引擎将其替换成一个值


```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    t, _ := temlate.ParseFiles("tmpl.html")
    t.Execute(w, "HelloWorld")
})
```

### 解析模板

* **`ParseFiles()`**

```go
func ParseFiles(filenames ...string) (*Template, error)
```

解析模板文件，并返回一个解析好的模板`struct`

实际上是调用`Tempalte.ParseFiles()`方法：
```go
t := template.New("tmpl.html")
t, _ = t.ParseFlies("tmpl.html")
```

`ParseFiles()`参数数量可变，但只返回一个模板，当解析多个文件是，第一个文件作为返回的模板，其余作为`map`，供后续执行使用

* **`ParseGlob()`**

使用模式匹配解析特定文件：
```go
t, _ := tempalte.ParseGlob("*.html")
```

* **`Parse()`**

解析字符串模板，其它方式最终都会调用`Parse()`

***

* `Lookup()`: 通过模板名来寻找模板

* `Must()`: 包裹一个函数，返回一个模板指针和错误(不为`nil`是`panic()`)

### 执行模板

* **`Execute()`**
  * 参数是ResponseWriter 数据
  * 适用于单模板，模板集只用第一个模板

* **`ExecuteTemplate()`**
  * 参数是ResponseWriter 模板名 数据
  * 适用于模板集


```go
http.Handlefunc("/process", func(w http.ResponseWriter, r *http.Request) {
    t, _ := temlate.ParseFiles("tmpl.html")
    t.Execute(w, "HelloWorld")

    ts, _ := template.ParseFiles("t1.html", "t2.html")
    ts.ExecuteTemplate(w, "t2.html", "response for t2")
})
```

### Action

* Action就是go模板中嵌入的命令，位于两组花括号之间`{{xxx}}`

* `.`就是一个Action，而且是最重要的一个，代表传入模板中的数据

* Action可分为以下5类

#### 条件类Action

```go
{{if arg}}
  some content
{{end}}

{{if arg}}
  some content
{{else}}
  other content
{{end}}
```

`arg`可选为`.`，表示程序中传入的上下文数据如果为`true`，则执行if后的内容

#### 遍历/迭代 Action

```go
{{range array}}
  Dot is set to the element {{.}}
{{end}}
```

`array`可以是数组，`slice`，`map`或`channel`

此处的`.`表示每次迭代的元素，当`array = .`时，两处的`.`含义不同

提供回落机制，当`array`为空时触发：

```go
{{range array}}
  Dot is set to the element {{.}}
{{else}}
  Empty array.
{{end}}
```

#### 设置 Action

```go
{{with arg}}

{{end}}
```

在指定范围内，让`.`表示其它指定的值`arg`， 也有回落机制。

#### 包含 Action

```go
{{template "name"}}

{{template "name" arg}}
```

允许在模板中包含其它的模板，`arg`表示传递给子模板的参数

#### 定义 Action


### 函数与管道

#### 参数(argument)

* *参数*是模板里面用到的值
  * 可以是整数，`bool`，`string`等
  * 也可以是`struct`，字段，数组的key等
* 参数可以是变量、方法(返回单个值+一个可选错误)或函数
* 参数可以是一个`.`，表示传入模板引擎中的值


```go
{{if arg}}
  some content
{{end}}
```

此处的`arg`就是参数

#### 在Action中设置变量

* 可以在action中设置变量，以`$`开头：
  * `$variable := value`

```
{{ range $key, $value := . }}
 {{$key}} -> {{$value}}
{{ end }}
```

#### 管道*pipeline*

管道是按顺序连接到一起的参数、函数或方法，类似Linux中的管道

允许我们把前一个参数的输出发给下一个参数：`p1 | p2 | p3`

#### 函数

Go模板引擎提供了一些基本的内置函数，功能比较有限。

开发者也可以自定义函数，但要求可以接收任意数量的参数，而返回一个值或一个值和一个错误

* **内置函数**
  * define、template、block 组合模板相关
  * html、js、urlquery 对字符串进行转义
    * 如果是web模板，那么不会经常使用这些函数
  * index 访问指定下标
  * print printf println
  * len  
  * with  设置参数

* **自定义函数**

```go
func (t *Template) Funccs(funcMap FuncMap) *Template
type FuncMap map[string]interface{}
```

1. 创建一个`FuncMap`变量
  * key是(模板中的)函数名，value是函数
2. 将`FuncMap`附加到模板

### 组合模板

#### layout 模板

Layout模板就是网页中固定的部分，可以被多个网页重复使用；

例如导航栏，页脚信息等

***

通常，创建一个`layout.html`，该文件定义了一个公有的部分

然后通过`{{template content}}`包含其它具体模板

为避免模板名的问题，可以通过`{{define content}}`自定义模板名

在代码中实现时，`ParseFlies()`转换所有所需的文件，执行时通过`ExecuteTemplate()`调用`layout`对应的模板即可

***

这样在实现每个网页的特定功能时，都可以避免重复代码，虽然它们都是同样的模板名，但是在http路由时的路径不同，加载的文件也就不同。

#### block 定义默认模板

```
{{ block arg}}
default-template
{{end}}
```

如果可以找到`arg`对应的模板文件就使用这个模板，否则使用自定义的默认模板。

### 逻辑运算符

* eq/ne
* lt/gt
* le/ge
* and
* or
* not

```
{{ if eq . "hello world"}}
  true
{{else}}
  false
{{end}}
```

## 8 数据库

这里主要用到关系型数据库，关于这部分内容前面已经学习过，此处仅作茶漏。

```go
func (*DB) PingContext(context.Context)
```

接收一个上下文参数，验证连接数据库:

```go
ctx := context.Background()
if err := db.PingContext(ctx); err != nil {
  panic()
} 
// conn sql ok
```

Context类型可以携带截至时间、取消信号和其它请求范围的值，并且可以横跨API边界和进程。

上例中，通过`context.Background()`创建一个非`nil`的`Context`对象；

它不会被取消，没有值，没有截至时间。

> 在Go语言中，`context`包提供了在Go程序中跟踪请求的上下文信息、控制请求的生命周期以及取消操作的机制。主要用途是在处理多个goroutine并发执行的场景下，有效地管理请求的上下文数据。

> 1. **传递请求相关的值**: `context.Context`类型可以用来传递请求的元数据，比如请求的唯一标识、认证信息等。这些信息可以通过`context.WithValue`方法添加到context中，然后在goroutine之间传递。

> 2. **控制请求的超时和取消**: 通过`context.WithTimeout`和`context.WithDeadline`方法，可以设定一个时间点或者超时时长，当超过这个时间时，相关的goroutine可以安全地退出。这种机制在避免因某个请求的处理时间过长而导致整个程序性能下降十分有效。

> 3. **处理请求的取消**: 通过`context.WithCancel`方法，可以创建一个可取消的context。当需要提前结束一个请求时，可以调用生成的`cancel`函数，通知所有使用该context的goroutine退出。

> 4. **跟踪请求的调用链**: `context.WithValue`方法可以用来传递请求的元数据，比如请求的唯一标识、认证信息等。


### 数据库操作

* 前面使用过`Qeury()`和`QueryRow()`两个方法，对应的还有`QueryContext()`和`QueryRowContext()`

* 对于`Exec()`，对应有`ExecContext()`

## 9 路由

路由：根据不同的请求路径pattern，触发不同的handler。

### Controller的角色

前面设置路由和启动服务都放在`main()`，显然不是很合适的：

* `main()` 设置类工作
* controller
  * 静态资源
  * 把不同的请求送到不同的controller进行处理


因此需要在项目下创建一个`package controller`，该包实现各种handler并提供导出的注册函数；

此外，不宜将所有handler实现在同一go文件中

### 路由参数

* 静态路由：一个路由对应一个页面
  * /home
  * /about

* 带参数的路由：根据路由参数，创建出一族不同的页面
  * /companied/123
  * /companies/google

### 第三方路由器

* gorilla/mux: 功能强大，性能相对较差
* httprouter: 功能简单，注重性能

## 10 Json

* Go结构体标签映射Json键：

```go
type Company struct {
  ID        int     `json:"id"`
  Name      string  `json:"name"`
  Country   string  `json:"country"`
}
```

* 类型映射
  * Go bool --- Json boolean
  * Go float64 --- Json 数值
  * Go string --- Json strings
  * Go nil --- Json null

* 未知结构的Json
  * `map[string]interface{}`可以存储任意Json对象
  * `[]interface{}`可以存储任意的Json数组

* Json读写
  * 针对`string`或`bytes`
    * `Marshal => string`
    * `Unmarshal <= string`
  * 针对*stream*
    * `Encode => stream`
    * `Decode <= stream`

## 11 中间件-middleware

* 示意
```shell
request   ->  |              |   ->   |
              |  middleware  |        |  handler
response  <-  |              |   <-   |
```

类似`preHandle()`和`postHandle()`

* 创建中间件

可以在创建一个Handler，这个Handler会调用注册好的handler，但在调用之前和调用之后就可以做一些操作：

```go
type Middleware struct {
  // 下一层，如果只有一层中间件那么则为defaultServeMux
  Next http.Handler
}

func (m *Middlware) ServeHttp(w http.ResponseWriter, r *http.Request) {
  preHandle()
  Next.ServeHttp(w, r)
  postHandle()
}
```

* 中间件的用途
  * Logging
  * 安全
  * 请求超时
  * 响应压缩

## 12 请求上下文

在handler处理各种请求时，可能访问数据库、web-service或者文件系统，这期间可能需要中间件或其它地方传递的信息，比如用户信息，操作时间等等，这就需要用到上下文对象：

* `request context`
```go
// 返回当前请求的上下文
func (*Request) Context() context.Context

// 基于给定的上下文，设置当前请求上下文对象
func (*Request) WithContext(ctx context.Context) context.Context
```

* `context.Context`
```go
type Context interface {
  Deadline() (deadline time.Time, ok bool)
  Done() <-chan struct{}
  Err() error
  Value(key interface{}) interface{}
}
```

## 13 HTTPS

传统的http协议在数据传输时是以明文传输的，会缺乏安全性。

https添加了tls保证了传输层的数据安全。

```go
func http.ListenAndServeTLS(addr string, certFile string, keyFile string handler http.Handler)
```

其中`certFile`和`keyFile`是与数据加密的相关文件。

go提供了生成证书文件的方法：

```shell
# 查看帮助
go run ${GOROOT}/src/crypto/tls/generate_cert.go -h

# 生成证书文件 cert.pen key.pen
go run ${GOROOT}/src/crypto/tls/generate_cert.go -host localhost
```

当将http改为https时，对应使用的版本自动从http1.1升级到2.0

## 14 HTTP/2

http1.1中header与body捆绑严重，而且header不能被压缩；

http2仍然采用tcp传输，会在两者之间建立起stream(类似通道)，数据以frame的形式发送

* http/2
  * 请求多路复用
  * Header压缩
  * 默认安全
    * HTTP
    * HTTPS
  * Serve Push

* **Server Push**

当客户端接收一个响应进行解析时，发现这个文件还有其他的依赖，因此会再次发送请求获取相关文件。

如果提供ServerPush支持，服务器可以在响应目标文件时先响应相关依赖文件，最后响应请求文件，这样客户端就不必多次发起请求。

```go
func handleHome(w http.ResponseWriter, r *http.Request) {
  // 如果支持server push
  if pusher, ok := w.(http.Pusher); ok {
    pusher.Push("/css/app.css"， &http.PushOptions{
      Header: http.Header {"Content-Type": []string{"text/css"}},
    })
  }

  // 响应请求文件...
}
```

## 15 测试

`net/http/httptest`提供了http测试的常用函数：

```go
// 获取一个Request
func NewRequest(method, url string, body io.Reader) (*Request, error)
```

```go
// 响应记录器
type ResponseRecorder {
  Code int // 状态码
  HeaderMap http.Header
  Body *bytes.Buffer 
  Fulshed bool
}
```

```go
func TestHanleCompany(t *testing.T) {
  r := httptest.NewRequest(http.MethodGet, "/company", nil)
  w := httptest.NewRecoder()

  handlerCompany(r, w)

  result, _ := ioutil(w.Result().Body)
  c := company{}
  json.Unmarshal(result, &c)
  if c.Id != 123 {
    t.Errorf("failed test")
  }
}
```

## 16 profiling 性能分析

* Go的性能分析提供如下功能
  * 内存消耗
  * CPU使用
  * 阻塞的goroutine
  * 执行追踪
  * 提供一个web界面：应用的实时数据

* 使用

```go
import (
  _ "net/http/pprof"
)
```

匿名引用相关的包，会为web应用绑定几个性能分析相关的handler；

要避免性能分析的hanler受到业务handler的影响，可以使用一个go承载：
```go
go http.ListenAndServe("localhost:8080", nil)
```

* 访问
  * 直接在web页面查看
  * 在命令行界面中使用`go tool`访问


## 17 总结

实现一个支持增删改查的web界面。