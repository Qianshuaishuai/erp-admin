package controllers

import "erp-admin/models"

type ProductController struct {
	BaseController
}

func (self *ProductController) List() {
	self.Data["pageTitle"] = "产品管理"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *ProductController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	// q 查询条件
	q := self.GetString("q")

	id, _ := self.GetInt("id")
	if id > 0 {
		models.DeletePersonTag(id)
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetProductList(q, limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name_sku"] = v.Name + " / " + v.SKU
		row["brand"] = v.Brand
		row["developer"] = v.Developer
		row["number"] = v.Number
		row["purchase_cost"] = "￥" + v.PurchaseCost
		row["preferred_supplier"] = v.PreferredSupplier
		row["classify"] = v.Classify
		row["create_time"] = v.CreateTime.Format("2006-01-02 15:04:05")

		if v.SaleStatus == 0 {
			row["sale_status"] = "在售"
		} else if v.SaleStatus == 1 {
			row["sale_status"] = "停售"
		} else if v.SaleStatus == -1 {
			row["sale_status"] = "未知"
		}

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}
