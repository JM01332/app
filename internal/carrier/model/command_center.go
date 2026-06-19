package model

// CommandCenter represents a row in the command_center table.
type CommandCenter struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CodeName      string `gorm:"column:code_name;not null"`
	SecurityLevel int    `gorm:"column:security_level;not null"`
	CarrierID     int64  `gorm:"column:carrier_id;not null;unique"`
}

// TableName returns the existing PostgreSQL table name.
func (CommandCenter) TableName() string {
	return "command_center"
}
