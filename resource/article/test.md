# 案例：数据表格的增删查改

## 1 web页面

先实现一个基本的网页，能够实现表格的增删查改，然后在将其转化为模板；

这里由于没有服务器因此每次刷新都会到最初的状态：

起初将`<button>`标签置于`<form>`中，在这种情况下，即使重写了按钮的点击事件，也会触发刷新：

> 当 `<button>` 标签被 `<form>` 标签包裹时，点击按钮会触发表单的提交行为。默认情况下，点击按钮会导致页面刷新，以提交表单数据。为了防止页面刷新，可以：

> 1. **使用 `type="button"`**：确保按钮的类型为 `button`，而不是 `submit`。例如：
>   ```html
>   <button type="button" onclick="myFunction()">Click me</button>
>   ```
>
> 2. **阻止表单提交**：如果按钮的类型是 `submit`，可以使用 JavaScript 阻止表单提交：
>   ```javascript
>   document.querySelector('button').addEventListener('click', function(event) {
>       event.preventDefault();
>       // 处理按钮点击事件
>   });
>   ```

> 这些方法可以防止在点击按钮时触发表单提交，从而避免页面刷新。

这里我是暂时将`<form>`替换为`<div>`，由于通过类名选择，因此也还会保持样式。

另外，要注意添加新的单元格时也设置其类名(通过`className`属性)，使其可以加载预设的样式。

## 2 服务器响应主页

### scp

首先要把windows上写好的网页上传到Linux上，这里可以使用`scp`命令：

`scp`*Secure Copy Protocol*是一种用于在两台计算机之间安全地传输文件的命令行工具。

`scp`基于SSH协议，因此在传输过程中提供了加密保护。

* 基本语法
```shell
scp [选项] [源路径] [目标路径]
```

由于是网络传输，因此指定目标路径时类似`ssh`，还需要用户名和ip

* 例：上传本地某个目录下所有文件到指定路径
```shell
scp * ty@192.168.18.128:/home/ty/resource
```

* 选项

`-r`: 递归复制目录及其内容

`-P`: 指定SSH端口


### 层次设计

一个完整的 Web 应用层次结构通常包括以下几个层次：

1. **Client (前端)**：
   - **作用**：用户界面，呈现数据并处理用户交互。使用 HTML, CSS, JavaScript, 前端框架（如 React、Vue）等技术。
   - **功能**：显示数据、收集用户输入、发起请求到服务器。

2. **Controller (控制器) 层**：
   - **作用**：处理用户请求，协调 Model 层和 View 层的交互，控制应用程序的流程。
   - **功能**：接收和解析请求、调用 Model 层处理数据、选择合适的视图进行响应。

3. **Model (模型) 层**：
   - **作用**：管理应用程序的业务逻辑和数据访问。与数据库交互，执行数据验证和处理业务规则。
   - **功能**：定义数据结构、执行数据库操作、封装业务逻辑。

4. **View (视图) 层**：
   - **作用**：生成用户界面。负责将数据呈现给用户，通常由模板引擎生成 HTML。
   - **功能**：根据控制器传递的数据生成前端界面。

5. **Middleware (中间件) 层**：
   - **作用**：处理请求和响应的通用功能，如认证、授权、日志记录、请求解析等。
   - **功能**：拦截请求、执行预处理、增强请求和响应处理流程。

6. **Database (数据库) 层**：
   - **作用**：存储和管理应用程序的数据。可以是关系型数据库（如 MySQL、PostgreSQL）或非关系型数据库（如 MongoDB）。
   - **功能**：提供数据存储、检索和管理服务。

7. **Service (服务) 层**（可选）：
   - **作用**：封装业务逻辑，提供跨多个控制器和模型的功能，促进代码的复用。
   - **功能**：实现业务规则、处理复杂的操作和协调不同的组件。

这种分层结构有助于解耦不同的应用程序组件，提高代码的可维护性和可扩展性。

### 响应主页

#### model

1. 实现数据结构定义 `type UserInformation struct {}`
2. 实现一个查询当前所有用户的函数 (这里还没有引入数据库，可以暂时返回一个空切片)

#### view

1. 修改静态html为模板文件 (主要时数据行的部分)
2. 加载模板并对外提供加载好的模板对象

这里要注意`tempalte.New()`指定的模板名要与模板文件中的定义一致，如果没有定义那么默认是文件名。

#### cotroller

1. 注册*home-handler*并对外提供注册的函数
2. 注册`FileServer`响应要加载的css/js文件

## 3 持久化存储

考虑了一下，还是放到数据库里，如果放文件里感觉会有点麻烦。

### 创建表

仅仅是一个简单的demo，一张表即可，也没有什么设计可言，这里直接采用自增主键。

```sql
create table user_informaions (id int unsigned auto_increment primary key, name varchar(15), gender varchar(10), introduction varchar(25), password varchar(15));
```

### model

1. 提供一个导出的数据库连接池句柄
2. 在`init()`中初始化连接池句柄
3. 提供一个导出的查询存储表所有记录的方法

在出现*比较严重*的错误时，`log`包提供了`panic`和`fatal`两种机制：

* `panic`
  * 引发一个运行时错误，可被`recover()`捕获
  * 进行栈回溯，执行`defer`
* `fatal`
  * 立即终止程序`os.Exit(1)`

通常fatal(致命的)比panic(恐慌)更加严重，前者处理配置错误，初始化错误等，后者用于不可恢复错误。

### 依赖处理

需要下载mysql的驱动，然后在`go.mod`中添加依赖：
```shell
go get -u ...
go mod tidy
```

### 测试

添加两条记录测试获取全部用户信息的功能：
```shell
insert into user_informations(name, gender, introduction, password) values('ty', 'male', 'I am ty.', '123456');

insert into user_informations(name, gender, introduction, password) values('uu', 'female', 'uu!', '369369');
```

***

结果时没有问题的，但现在网页中的按钮没有提交表单的提供，只是自己操作自己。

到这里也是大体也是比较完整了，接下来搞一下各种请求处理就好了

## 4 路由处理

这里以后基本就不会很麻烦了，都是增删改查之类的相对容易的操作。

### 4.1 新增一条记录

1. 写一个注册用户的页面`add_user.html`(也可以在主页使用提示窗)
2. 注册主页请求注册新用户的路由，即返回上述页面

不难发现，两个页面可能会有相同的部分，比如导航栏、备案信息等等，因此这里需要进行*组合模板*

***

我直接一手小天才设计：

* 目标：将所有的模板文件加载都同一个`template.Template`对象中：
  1. 要替换的块不能都命名为`content` --> 传入参数表示要加载的块+条件判断
  2. 参数还要有一些块需要的信息 --> 将加载的块名字和信息封装到一个结构中
  3. 傻逼3.5

***

* 返回状态

每次操作后如果只是返回状态码，那么浏览器显示一片空白，效果非常不好。

因此设计一个状态页面`response_status.html`，返回一些状态描述。


### 删除一条记录

1. 点击按钮后获取该行的id，向服务器发送请求
2. 服务器操作后返回一个状态

#### XMLHTTPRequest

#### fetch api

#### CROS策略

Q:CROS策略是怎样的，都是针对同一服务器的get请求，如何通过表单实现发送请求，就不需要设置，如果通过`fetch()`发送请求，就要设置？

CORS（跨源资源共享）策略用于限制从一个源（协议、域名、端口）发起的 Web 请求对不同源的资源的访问。它主要用于确保客户端安全，防止恶意网站访问用户数据。

##### CORS 策略：

1. **同源策略**：浏览器默认的安全策略，只允许同源（相同协议、域名、端口）之间的请求。这就是为什么你直接通过表单提交或从同源的 `<script>` 标签发起请求时，通常不需要处理 CORS。

2. **跨源请求**：如果请求涉及不同的源，比如通过 JavaScript 的 `fetch()`、`XMLHttpRequest` 或 AJAX 发送请求，服务器必须设置适当的 CORS 头以允许请求。这通常是因为浏览器需要额外的安全检查来确认是否允许这些请求。

##### 表单提交 vs. `fetch()` 请求：

- **表单提交**：当使用表单（`<form>` 标签）提交数据时，浏览器的默认行为是进行同步请求，并且不会受到 CORS 的限制。表单提交的请求通常是简单请求（simple request），例如 POST 请求时不需要 CORS 头部，因为这类请求被认为是安全的。

- **`fetch()` 请求**：当使用 `fetch()` 发起请求时，特别是当涉及到跨源的 POST、PUT、DELETE 等请求，浏览器会发送预检请求（OPTIONS 请求）。这就需要服务器正确设置 CORS 头部来允许这些请求。

###### 如何设置 CORS 头部：

在服务器端，你需要设置以下 CORS 头部来允许跨源请求：

```go
w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有源
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // 允许的方法
w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // 允许的头部
```

对于 `fetch()` 请求，确保：

1. **配置正确的 CORS 头**：服务器必须设置 `Access-Control-Allow-Origin` 等头部。
2. **处理预检请求**：对于非简单请求，服务器需要正确处理 OPTIONS 请求并返回必要的 CORS 头。

这样配置之后，跨源请求会被允许，通过 `fetch()` 和其他 AJAX 请求可以顺利访问不同源的资源。

***

我操这些东西！？

### 修改一条记录

1. 为单元格添加点击事件，接收用户数据
2. 考虑简单实现就是直接通过输入框或跳转表单，也可以内联编辑
3. 每次单次编辑不会发出请求，点击按钮后发送

* JavaScript 处理单元格点击事件代码的详细解释：

```javascript
document.addEventListener('DOMContentLoaded', () => {
  const table = document.getElementById('editable-table');

  table.addEventListener('click', function(event) {
    const target = event.target;

    if (target.classList.contains('editable')) {
      const currentValue = target.innerText;
      const input = document.createElement('input');
      input.type = 'text';
      input.value = currentValue;
      input.classList.add('editing');

      target.innerHTML = ''; // Clear cell content
      target.appendChild(input);

      input.focus();

      input.addEventListener('blur', () => {
        const newValue = input.value.trim();
        if (newValue !== currentValue) {
          target.innerText = newValue;
          // Here, you can send an AJAX request to save the new value
          console.log('Updated value:', newValue);
        } else {
          target.innerText = currentValue;
        }
      });
    }
  });
});
```

#### 代码解释

1. **`document.addEventListener('DOMContentLoaded', () => { ... });`**
   - 确保 DOM 内容完全加载后再执行代码。这避免了在页面元素尚未加载时尝试访问它们。

2. **`const table = document.getElementById('editable-table');`**
   - 获取表格元素的引用。

3. **`table.addEventListener('click', function(event) { ... });`**
   - 给表格添加一个点击事件监听器，当表格中的任何单元格被点击时，都会触发这个事件。

4. **`const target = event.target;`**
   - 获取点击事件的目标元素，即被点击的单元格。

5. **`if (target.classList.contains('editable')) { ... }`**
   - 检查被点击的单元格是否具有 `editable` 类，确保只有可编辑的单元格会被处理。

6. **`const currentValue = target.innerText;`**
   - 获取当前单元格中的文本值，供后续比较使用。

7. **`const input = document.createElement('input');`**
   - 创建一个新的文本输入框。

8. **`input.type = 'text';`**
   - 将输入框的类型设置为文本。

9. **`input.value = currentValue;`**
   - 将输入框的初始值设置为单元格的当前值。

10. **`input.classList.add('editing');`**
    - 给输入框添加 `editing` 类，以便可以应用特定的样式（如背景色变化）。

11. **`target.innerHTML = '';`**
    - 清空单元格的内容，为输入框腾出空间。

12. **`target.appendChild(input);`**
    - 将输入框添加到单元格中，使其成为可编辑内容。

13. **`input.focus();`**
    - 聚焦到输入框中，以便用户可以立即开始编辑。

14. **`input.addEventListener('blur', () => { ... });`**
    - 添加一个失去焦点事件监听器，处理用户完成编辑后的操作。

15. **`const newValue = input.value.trim();`**
    - 获取并修剪输入框中的新值。

16. **`if (newValue !== currentValue) { ... }`**
    - 如果新值与当前值不同，则更新单元格内容，并可以在此处添加 AJAX 请求以保存更改。

17. **`target.innerText = newValue;`**
    - 将单元格的内容更新为新值。

18. **`console.log('Updated value:', newValue);`**
    - 在控制台输出更新后的值（此处可以用 AJAX 请求代替）。

19. **`else { target.innerText = currentValue; }`**
    - 如果新值与当前值相同，则恢复单元格的原始内容。

这个代码实现了一个简单的内联编辑功能，让用户可以直接在表格中编辑数据。

***

这里注意到，调用`fetch()`需要处理CROS，对于预检请求*preflight request*，可以放到中间件中去处理：

1. 定义一个`CrosMiddleware`，实现`http.Handler`
2. 检测当前请求是否涉及到cros，如果是，则设置相应的头部字段
3. 检测当前请求是否是预检请求，如果是，响应OK
4. 如果不是预检请求，则下放`controller`层

这里将js文件移植到服务器上发现请求时使用缓存的，禁用一下缓存就可以了。

### 查询若干记录

这部分就没有什么特别的了，其实但从前端也是也可以完成的。

这里就简单实现按名称查找和按id查找。

1. 类似`/home`的请求，不过需要进行筛选
2. 服务端实现起来也基本是相同的，主要集中于web端的实现

刚折腾破网页了：搞了半天最后还是用了原先添加用户的页面，没办法越弄越丑，然后又学了个单选框。

好，差不多了，基本的都能跑了。

## 5 TODO

考虑一些拓展，增强功能。