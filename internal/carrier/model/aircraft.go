package model

// Aircraft represents a row in the aircraft table.
type Aircraft struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Model        string `gorm:"column:model;not null"`
	Manufacturer string `gorm:"column:manufacturer;not null"`
	CarrierID    int64  `gorm:"column:carrier_id;not null;index:aircraft_carrier_id_idx"`
}

// TableName returns the existing PostgreSQL table name.
func (Aircraft) TableName() string {
	return "aircraft"
}
