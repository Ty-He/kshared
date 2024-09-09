# fetch 

`Fetch API`是现代`JavaScript`提供的一种用于发送 HTTP 请求和处理响应的接口。它比传统的 XMLHttpRequest 更加简洁和强大，支持异步操作，使得处理网络请求更为高效和易于管理。

## 使用fetch发送基本get请求

* 如果`fetch()`只接收一个`url`字符串，表示默认发送一个`get`请求，会返回一个*Promise对象*
* 如果需要设置`get`的参数，直接拼接在字符串上即可

```js
fetch(url)
.then(response => { return response.json(); /* 也是异步操作 */ })
.then(json => { console.log(json) })
.catch(err => { console.log(err) })
```

上述的异步操作如果用回调表示：
```C++
async_fetch(url, [](response res) {
    // ...
})
```

而异步操作可以改写为使用协程的类似的同步代码：
```js
async function getResponse() {
    let res = await fetch(url);
    let json = await res.json();
}
```

异常捕获也是同步的思路：
```js
async function getResponse() {
    try {
        let res = await fetch(url);
        let json = await res.json();
    } catch (err) {
        console.log(err)
    }
}
```

* 在`get`请求中指定查询字符串
```js
async function getResponse() {
    query = getQueryString(); // id=2&thread_id=3
    let res = await fetch(url + "?" + query);
    let json = await res.json();
}
```

然后在go的服务器端可以通过请求对象的`URL`字段或`Form`字段获取相应数据。

## response对象

通过调用`fetch()`可以获取到服务器的http响应，即`response`对象(下面用`res`表示)

|属性|含义|
|---|----|
|`res.ok`|返回一个`bool`值，表示请求是否成功|
|`res.status`|返回一个数字，即http响应的状态码|
|`res.res.statusText`|返回状态码对应的文本描述|
|`res.url`|返回请求的url，即`fetch()`传入的url|

* `response`实现了Body接口，因此有以下读取数据的方法：
  * `json()`
  * `text()`
  * `blob()`
  * `formData()`
  * `arrayBuffer()`

<https://developer.mozilla.org/zh-CN/docs/Web/API/Response>


## fetch配置参数

`fetch()`还可以接收第二个参数，作为配置对象，可以自定义发出http请求。

```js
fetch(url,
    {
        // 请求方法
        method: 'post',
        headers: {
            // 设置头部字段
        },
        // 请求数据体
        body: ''
    }
)
```

* 示例：发送一个post请求，数据体为json

```js
async function sendPost() {
    let jsonObj = {
        key1: 'val1',
        key2: 'val2'
    };

    let res = await fetch(url, 
        {
            method: 'post',
            headers: {
                // 对应表单的enctype属性
                'Content-Type':'application/json'
            },
            body: JSON.stringify(jsonObj)
        }
    );

    let json = await res.json();
    
    console.log(json);
}
```

## fetch函数封装

直接使用`fetch()`尽管简单很多，但还是有繁琐的地方，因此可以封装一下常用操作。

```js
// 针对get delete
async function http(obj) {
    // 结构赋值，类似结构化绑定
    let {method, url, params, data} = obj;
    // 处理params 如果有需要附加到url中
    if (params) {
        let str = new URLSearchparams(params).toString();
        url += '?' + str;
    }

    // data 如果有需要写headers
    var res;
    if (data) {
        res = await fetch(url, {
            method: method,
            headers: {
                'Content-Type':'application/json'
            },
            body: JSON.stringify(data)
        });
    } else {
        res = await fetch(url)
    }
    return res.json();
}
```

