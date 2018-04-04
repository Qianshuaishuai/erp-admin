package controllers

var (
	//(1-题干错别字，2-答案错误，3-解析错误，4-题目超钢，5-其它)
	INCORRECT_TYPE = map[int]string{
		1:"题干错别字",
		2:"答案错误",
		3:"解析错误",
		4:"题目超钢",
		5:"其它",
	}

	CORRECT_STATUS_FLAG = [3]string{
		"<span class='layui-badge layui-bg-orange'>未处理</span>",
		"<span class='layui-badge layui-bg-gray'>修改中</span>",
		"<span class='layui-badge layui-bg-green'>已发布</span>",
	}
)

type CorrectController struct{
	BaseController
}

