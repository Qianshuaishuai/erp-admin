## 试卷后台管理系统

试卷后台管理系统使用Layui+Beego，后台程序员借助Layui也可以搭建出不错的前端界面，需要对go的html模板有一定的了解。

Layui文档：
http://www.layui.com/doc/element/layout.html

Beego模板文档：
https://beego.me/docs/quickstart/view.md


### 地址
http://120.77.238.241:7730/home

### 运行
```
go build main.go
./dreamEbagPaperAdmin
```

### 目录介绍

> 此项目结构类似dreamEbagPapers，只是多了一些前端的东西。

- static
项目用到的各类图片、Js、Lib、Css等

- views
项目的模板网页，每个controller都有一个文件夹对应。

### 项目简介

这个项目不会发布到正式环境，只在测试服务器上运行，公司内部使用。在测试数据库中，为此项目建了几张表，下面进行介绍：

|表名|说明|
|---|---|
|t_users|用户表|
|t_add_paper_temps|临时存放新试卷
|t_chapter_temps|临时存放新试卷对应的章节、题目|
|t_check_data|审核表，记录着系统里的所有重要操作|

|帐号|密码|
|---|---|
|Super|Yufeng123|

Super帐号拥有所有的权限，并且可以创建删除数据员、审核员。
> 关于权限问题，里面主页会显示系统使用说明，必读！

### 项目进度

|模块|进度|
|---|---|
|查看、编辑试卷|开发完毕，可以使用|
|试题模块|开发完毕，可以使用|
|审核机制|开发完毕，可以使用|
|上报问题处理|开发完毕，可以使用|
|权限管理|开发完毕，可以使用|
|添加试卷|目前开发到可以添加、删除新试卷；添加、删除、修改、排序新章节。但是不能添加题目；不能在测试环境下查看新试卷；也不能把新试卷发送到正式环境|

### 工作提醒

- 这个管理系统维护着真题试卷的试卷、试题数据库，我的前端基础很差，尽量做到能用、好用，但可能前端代码工程性不高，有一些重复代码有待完善。
- 这是同步试卷的试卷管理后台可以感受一下：http://192.168.20.235:8000/sjlr/index.php/Paper/index.html
- 同步试卷和真题试卷是分属不同的数据库，各有一套数据和格式，各有一套管理系统，数据整合工作可谓路漫漫其修远兮。
- 请明确测试环境和正式环境的区别，[测试环境](http://dreamtest.strongwind.cn:7350/login.html)是我们测试看到的数据，[正式环境](https://teacher.ebag.readboy.com/)是老师用的数据，真题试卷管理系统的所有操作都针对测试环境，只有审核通过，才会发送新数据到正式环境，这时老师就可以看到数据了。所以测试环境数据和正式环境数据**可能不同步**，要注意这一点区别。