package mysqlapi

type MysqlApi struct {
	// Address of MySQL-API
	Address string `json:"address"`
	// Error messages during the health request
	Message string `json:"message"`
	// Response of the health request
	Response *HealthView `json:"response"`
}
