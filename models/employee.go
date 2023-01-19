package models

type Employee struct {
	EmpId   int    `json:"empId" binding:"Required"`
	EmpName string `json:"empName" binding:"Required"`
	Address string `json:"address" binding:"Required"`
	Phone   int64  `json:"phone" binding:"Required"`
	AssetId int64  `json:"assetId" binding:"Required"`
}
