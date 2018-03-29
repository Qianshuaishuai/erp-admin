package controllers

type QuestionController struct {
	BaseController
}


func (self *QuestionController) Search(){
	self.Data["pageTitle"] = "搜索试题"
	self.Data["ApiCss"] = true

	self.display()
}