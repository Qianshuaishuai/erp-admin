package models

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type Product struct {
	ID                       int       `gorm:"id" json:"id"`
	SKU                      string    `gorm:"sku" json:"sku"`
	Name                     string    `gorm:"name" json:"name"`
	SaleStatus               int       `gorm:"sale_status" json:"sale_status"`
	Unit                     string    `gorm:"unit" json:"unit"`
	Classify                 int       `gorm:"classify" json:"classify"`
	Brand                    int       `gorm:"brand" json:"brand"`
	Number                   string    `gorm:"number" json:"number"`
	Developer                string    `gorm:"developer" json:"developer"`
	Leader                   string    `gorm:"leader" json:"leader"`
	Desc                     string    `gorm:"desc" json:"desc"`
	PicURL                   string    `gorm:"pic_url" json:"pic_url"`
	Purchaser                string    `gorm:"purchaser" json:"purchaser"`
	PurchaseCost             string    `gorm:"purchase_cost" json:"purchase_cost"`
	DeliveryTime             string    `gorm:"delivery_time" json:"delivery_time"`
	MaterialQuality          string    `gorm:"material_quality" json:"material_quality"`
	SingleLength             float64   `gorm:"single_length" json:"single_length"`
	SingleWidth              float64   `gorm:"single_width" json:"single_width"`
	SingleHeight             float64   `gorm:"single_height" json:"single_height"`
	SingleWeight             float64   `gorm:"single_weight" json:"single_weight"`
	PackWidth                float64   `gorm:"pack_width" json:"pack_width"`
	PackHeight               float64   `gorm:"pack_height" json:"pack_height"`
	PackLength               float64   `gorm:"pack_length" json:"pack_length"`
	SingleRoughWeight        float64   `gorm:"single_rough_weight" json:"single_rough_weight"`
	BoxLength                float64   `gorm:"box_length" json:"box_length"`
	BoxWidth                 float64   `gorm:"box_width" json:"box_width"`
	BoxHeight                float64   `gorm:"box_height" json:"box_height"`
	BoxCount                 int       `gorm:"box_count" json:"box_count"`
	ContactName              string    `gorm:"contact_name" json:"contact_name"`
	ContactPhone             string    `gorm:"contact_phone" json:"contact_phone"`
	SupplierName             string    `gorm:"supplier_name" json:"supplier_name"`
	MinimumPurchaseCount     int       `gorm:"minimum_purchase_count" json:"minimum_purchase_count"`
	PreferredSupplier        string    `gorm:"preferred_supplier" json:"preferred_supplier"`
	SupplierPurchasePriceCNY string    `gorm:"supplier_purchase_price_cny" json:"supplier_purchase_price_cny"`
	SupplierPurchasePriceUSD string    `gorm:"supplier_purchase_price_usd" json:"supplier_purchase_price_usd"`
	PurchaseURL              string    `gorm:"purchase_url" json:"purchase_url"`
	CreateTime               time.Time `gorm:"create_time" json:"create_time"`
}

type Brand struct {
	ID   int       `gorm:"id" json:"id"`
	Name string    `gorm:"name" json:"name"`
	Time time.Time `gorm:"time" json:"time"`
}

type Classify struct {
	ID   int       `gorm:"id" json:"id"`
	Name string    `gorm:"name" json:"name"`
	Time time.Time `gorm:"time" json:"time"`
}

func GetProductList(q string, limit int, page int, sort int) (list []Product, count int64) {
	db := GetEliteDb().Table("products")
	queryParams := make(map[string]interface{})
	list = make([]Product, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	// 将搜索字符串按空格拆分
	q = strings.TrimSpace(q)
	var qstring string
	if len(q) > 0 {
		qs := strings.Fields(q)
		for _, v := range qs {
			qstring += "%" + v
		}
		qstring += "%"
	}

	if len(qstring) > 0 {
		db = db.Where("name LIKE ?", qstring)
	}

	// var sortStr = "DESC" // 默认时间 降序
	// if sort == 1 {
	// 	sortStr = "ASC"
	// }

	db = db.Model(UserInfoSimple{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		// Order("register " + sortStr).
		Scan(&list)

	return
}

func GetBrands() (datas []Brand) {
	GetEliteDb().Table("brands").Find(&datas)
	return
}

func GetClassifys() (datas []Classify) {
	GetEliteDb().Table("classifys").Find(&datas)
	return
}

func GetClassifySingle(name string, datas []Classify) (id int) {
	for d := range datas {
		if datas[d].Name == name {
			return datas[d].ID
		}
	}

	return 0
}

func GetBrandSingle(name string, datas []Brand) (id int) {
	for d := range datas {
		if datas[d].Name == name {
			return datas[d].ID
		}
	}

	return 0
}

func TranslateMoreProduct(fileName string, file multipart.File) (err error) {
	fileDir := "./" + "excel-import" + "/"

	os.RemoveAll(fileDir)
	os.Mkdir(fileDir, 0766)

	f, err := os.OpenFile(fileDir+fileName, os.O_CREATE|os.O_RDWR, 0766)
	io.Copy(f, file)

	xlsxfile2, err := xlsx.OpenFile(fileDir + fileName)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return
	}

	brands := GetBrands()
	classifys := GetClassifys()

	tx := GetEliteDb().Begin()
	var snow *MSnowflakCurl

	for _, sheet := range xlsxfile2.Sheets {
		rowCount := sheet.MaxRow
		for r := 0; r < rowCount; r++ {
			if r != 0 {
				row, _ := sheet.Row(r)
				sku := row.GetCell(0).String()
				name := row.GetCell(1).String()
				currentBrand := row.GetCell(5).String()
				currentClassify := row.GetCell(4).String()
				classifyID := GetClassifySingle(currentClassify, classifys)
				brandID := GetBrandSingle(currentBrand, brands)
				unit := row.GetCell(3).String()
				saleStatus := row.GetCell(3).String()
				saleStatusInt := 0
				number := row.GetCell(6).String()
				developer := row.GetCell(7).String()
				leader := row.GetCell(8).String()
				desc := row.GetCell(9).String()
				picURL := row.GetCell(10).String()
				purchaser := row.GetCell(11).String()
				purchaseCost := row.GetCell(12).String()
				DeliveryTime := row.GetCell(13).String()
				materialQuality := row.GetCell(14).String()
				singleLength, _ := row.GetCell(15).Float()
				singleWidth, _ := row.GetCell(16).Float()
				singleHeight, _ := row.GetCell(17).Float()
				singleWeight, _ := row.GetCell(18).Float()
				packLength, _ := row.GetCell(19).Float()
				packWidth, _ := row.GetCell(20).Float()
				packHeight, _ := row.GetCell(22).Float()
				singleRoughWeight, _ := row.GetCell(23).Float()
				boxLength, _ := row.GetCell(24).Float()
				boxWidth, _ := row.GetCell(25).Float()
				boxHeight, _ := row.GetCell(26).Float()
				boxCount, _ := row.GetCell(27).Int()
				contactName := row.GetCell(28).String()
				contactPhone := row.GetCell(29).String()
				supplierName := row.GetCell(30).String()
				minimumPurchaseCount, _ := row.GetCell(31).Int()
				preferredSupplier := row.GetCell(32).String()
				supplierPurchasePriceCNY := row.GetCell(33).String()
				supplierPurchasePriceUSD := row.GetCell(34).String()
				purchaseURL := row.GetCell(35).String()

				if saleStatus == "在售" {
					saleStatusInt = 0
				} else if saleStatus == "停售" {
					saleStatusInt = 1
				} else {
					saleStatusInt = -1
				}

				if classifyID == 0 {
					return errors.New("有产品分类未找到")
				}

				if brandID == 0 {
					return errors.New("有产品品牌未找到")
				}

				var product Product
				product.ID = snow.GetIntId(false)
				product.SKU = sku
				product.Name = name
				product.Brand = brandID
				product.Classify = classifyID
				product.Unit = unit
				product.SaleStatus = saleStatusInt
				product.Number = number
				product.Developer = developer
				product.Leader = leader
				product.Desc = desc
				product.PicURL = picURL
				product.Purchaser = purchaser
				product.PurchaseCost = purchaseCost
				product.DeliveryTime = DeliveryTime
				product.MaterialQuality = materialQuality
				product.SingleLength = singleLength
				product.SingleWeight = singleWeight
				product.SingleWidth = singleWidth
				product.SingleHeight = singleHeight
				product.PackHeight = packHeight
				product.PackLength = packLength
				product.PackWidth = packWidth
				product.SingleRoughWeight = singleRoughWeight
				product.BoxLength = boxLength
				product.BoxHeight = boxHeight
				product.BoxWidth = boxWidth
				product.BoxCount = boxCount
				product.ContactName = contactName
				product.ContactPhone = contactPhone
				product.SupplierName = supplierName
				product.MinimumPurchaseCount = minimumPurchaseCount
				product.PreferredSupplier = preferredSupplier
				product.SupplierPurchasePriceCNY = supplierPurchasePriceCNY
				product.SupplierPurchasePriceUSD = supplierPurchasePriceUSD
				product.PurchaseURL = purchaseURL

				err = tx.Table("products").Create(&product).Error

				if err != nil {
					return errors.New("创建失败")
				}
			}

		}
	}

	if err != nil {
		return errors.New("创建失败")
	}

	return nil
}
