# 1 前端初步构建

网页由三部分构成：*模板部分，特定组件，内容*；

这里先写一组静态网页，然后考虑模板替换；

前后端交互的话先前crub示例已经接触过，再做起来就比较简单了。

## note基本显示

包括以下两个函数：

* **`getfile()`**
  * 用`fetch()`像服务端请求md文件，接收后使用`marked`进行渲染；

* **`generate_toc`**
  * 根据转换后的结果，寻找h(1, 2, 3)标签生成TOC(Table Of Content)

## 索引项

对于article存储，采用分开存储的方案：内容使用文件系统存储，相关信息使用数据库存储：
  * 数据库存储结构化数据，便于搜索和管理。
  * 文件系统存储实际内容，易于编辑和备份。
考虑每一个记录具体包含的字段
  * title 标题
  * ahthor_id 作者
  * update_time 更新时间
  * type 分类
  * label 标签

这里考虑了一下，将分类与标签分离，分类相对更宽泛，标签描述更细致

比如: asio.md 
  * type: C++ network
  * label: (C++ network) aysnc callback co-routinue tcp concurrent event_look epoll iocp等等

用户表：
  * id
  * name

评论表
  * id
  * user_id
  * article_id
  * content
  * time

### 小结

终于差不多把前端页面的框架搭好了，主要包括两个页面：

  1. 主页，包含banner、缩略图、最新内容
  2. article详情页，包含article内容、对应目录

接下来写前后端交互的article查询。

# 前端请求article

## 请求article列表

进入主页时，会请求最新的5个article，应**在<a>的herf字段中使用模板替换查询字符串**，这样，在点击链接时会发送响应的请求。

进入分类页(当前未实现)，应按请求某一分类所有的article列表，这里可以考虑做分页，但不会是现在。

因为是动态网站，所以要考虑**自定义添加分类**，前面的分类都是预定义的，这里考虑在导航栏中添加一个分类，进入后显示所有分类的前若干篇文章，不宜在首页显示过多的分类。


