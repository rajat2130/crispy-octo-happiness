package models

type Asset struct {
	AssetId          int64  `json:"assetId" gorm:"PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	AssetDescription string `json:"assetDescription"`
}
