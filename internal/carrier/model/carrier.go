package model

import "time"

// Carrier represents a row in the carrier table and its relationships.
type Carrier struct {
	ID           int64       `gorm:"column:id;primaryKey;autoIncrement"`
	Version      int         `gorm:"column:version;not null;default:0"`
	Name         string      `gorm:"column:name;not null;unique"`
	Nation       string      `gorm:"column:nation;not null"`
	CarrierType  CarrierType `gorm:"column:carrier_type;type:carrier_type;not null"`
	Erzeugt      time.Time   `gorm:"column:erzeugt;not null;autoCreateTime"`
	Aktualisiert time.Time   `gorm:"column:aktualisiert;not null;autoUpdateTime"`

	CommandCenter CommandCenter `gorm:"foreignKey:CarrierID;references:ID;constraint:OnDelete:CASCADE"`
	Aircrafts     []Aircraft    `gorm:"foreignKey:CarrierID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName returns the existing PostgreSQL table name.
func (Carrier) TableName() string {
	return "carrier"
}
