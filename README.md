# 2024Happy
2024新春IT技能提升挑战   亿字节邵欣提交

## 目标：

​        medium是很多it博主分享知识的收费Blog，本次挑战目标是建立一个web应用，实现两个功能，点击"Create pdf files"按钮，抓取https://medium.com/?tag=software-engineering的top 10 pages并逐段翻译成完整中英文对照文件，保存成对应pdf文件。生成文件后，点击"Download"按钮，下载生成的10个pdf文件的zip包到本地。

## 解题思路

1. https://medium.com/?tag=software-engineering 

带了tag: software-engineering，所以需要找到带了”software-engineering “的tag的阅读量排名前十的文章

2. https://medium.com/tag/software-engineering/archive/2023

此网页界面带了很多api请求，其中有个api请求TagArchiveFeedQuery为:

```json
[
    {
        "operationName": "TagArchiveFeedQuery",
        "variables": {
            "tagSlug": "software-engineering",
            "timeRange": {
                "kind": "IN_YEAR",
                "inYear": {
                    "year": 2023
                }
            },
            "sortOrder": "NEWEST",
            "first": 10,
            "after": ""
        },
        "query": "..." (省略)
    }
]
```

此请求返回中参数 sortOrder修改为"MOST_READ" 即为查找阅读量最高的前10文章列表

```json
{
    "tagSlug": "software-engineering",
    "timeRange": {
        "kind": "IN_YEAR",
                "inYear": {
                    "year": 2023
                }
    },
    "sortOrder": "MOST_READ",
    "first": 10,
    "after": ""
}
```

3. 找到获得文章内容得api, PostViewerEdgeContentQuery 此文章id和host需要根据前面获取得文章列表内容单独解析。
4. 由于需要web应用，所以需要前端。

## 实际项目

### 演示地址： http://101.36.113.108:8081/

### 前端-React  

 前端采用react编写单页面展示。提供了安年进行“时间筛选” 按钮 ， 展示当前时间段的阅读量排名前十的文章，并且提供翻译和PDF导出功能。



### 后端 go-medium

采用goframe框架，编译方式查看 

[https://goframe.org/display/gf#all-updates]: 

提供了已经编译好的两个可执行文件在bin目录下 go-medium.exe (win) go-medium(linux) 

如需编译其它版本请修改”hack/config.yml“中的system属性，然后执行”gf build“命令

```yaml
gfcli:
  build:
    arch: "amd64"
    system: "windows"
   #system: "linux"
   #system: "darwin"
    mod: "none"
    packSrc: "manifest"
    version: "v1.0.0"
    output: "./bin/go-medium.exe"
```

实现了两个接口

1. article/list  返回top10阅读量的文章列表。可选参数year=2021。

   ```json
   // 列表中，每个文章返回的示例
   {
       "id": "c043a964e463", // 文章ID
       "clapCount": 5527,    // 点赞数
       "mediumUrl": "https://betterprogramming.pub.....",  // 链接
       "title": "“If software engineering is in demand, why is it so hard to get a software engineering job?”"  // 标题
   },
   ```

2. article/content 返回单个文章内容。必填参数id(文章ID) ,可选参数translate(0 不翻译 1翻译)

3. 增加缓存  减少请求次数 目前缓存1天

