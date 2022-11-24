package model

/******sql******
CREATE TABLE `comm_config` (
  `configId` varchar(200) NOT NULL,
  `configValue` varchar(1024) DEFAULT NULL,
  `description` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`configId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3
******sql******/
// CommConfig [...]
type CommConfig struct {
	ConfigID    string `gorm:"primaryKey;column:configId;type:varchar(200);not null" json:"-"`
	ConfigValue string `gorm:"column:configValue;type:varchar(1024)" json:"config_value"`
	Description string `gorm:"column:description;type:varchar(2000)" json:"description"`
}

// TableName get sql table name.获取数据库表名
func (m *CommConfig) TableName() string {
	return "comm_config"
}

// CommConfigColumns get sql column name.获取数据库列名
var CommConfigColumns = struct {
	ConfigID    string
	ConfigValue string
	Description string
}{
	ConfigID:    "configId",
	ConfigValue: "configValue",
	Description: "description",
}
