package mysql_transaction

import "time"

type DbParam struct {
	User          string
	Password      string
	Host          string
	Port          string
	Schema        string
	TableNameLike string
}

type Columns struct {
	TableName      string `gorm:"column:TABLE_NAME" json:"tableName"`
	ColumnName     string `gorm:"column:COLUMN_NAME" `
	DataType       string `gorm:"column:DATA_TYPE" `
	Ordinal        int    `gorm:"column:ORDINAL_POSITION" `
	COLUMN_COMMENT string `gorm:"column:COLUMN_COMMENT" `
	COLUMN_TYPE    string `gorm:"column:COLUMN_TYPE" `
}

type BeanColumn struct {
	Name           string
	Type           string
	Ordinal        int
	COLUMN_COMMENT string
	COLUMN_TYPE    string
}

type INNODBTRX struct {
	TrxId                   string    `json:"trxId"`                   //
	TrxState                string    `json:"trxState"`                //
	TrxStarted              time.Time `json:"trxStarted"`              //
	TrxRequestedLockId      string    `json:"trxRequestedLockId"`      //
	TrxWaitStarted          time.Time `json:"trxWaitStarted"`          //
	TrxWeight               int       `json:"trxWeight"`               //
	TrxMysqlThreadId        int       `json:"trxMysqlThreadId"`        //
	TrxQuery                string    `json:"trxQuery"`                //
	TrxOperationState       string    `json:"trxOperationState"`       //
	TrxTablesInUse          int       `json:"trxTablesInUse"`          //
	TrxTablesLocked         int       `json:"trxTablesLocked"`         //
	TrxLockStructs          int       `json:"trxLockStructs"`          //
	TrxLockMemoryBytes      int       `json:"trxLockMemoryBytes"`      //
	TrxRowsLocked           int       `json:"trxRowsLocked"`           //
	TrxRowsModified         int       `json:"trxRowsModified"`         //
	TrxConcurrencyTickets   int       `json:"trxConcurrencyTickets"`   //
	TrxIsolationLevel       string    `json:"trxIsolationLevel"`       //
	TrxUniqueChecks         int       `json:"trxUniqueChecks"`         //
	TrxForeignKeyChecks     int       `json:"trxForeignKeyChecks"`     //
	TrxLastForeignKeyError  string    `json:"trxLastForeignKeyError"`  //
	TrxAdaptiveHashLatched  int       `json:"trxAdaptiveHashLatched"`  //
	TrxAdaptiveHashTimeout  int       `json:"trxAdaptiveHashTimeout"`  //
	TrxIsReadOnly           int       `json:"trxIsReadOnly"`           //
	TrxAutocommitNonLocking int       `json:"trxAutocommitNonLocking"` //
}

func (INNODBTRX) TableName() string {
	return "INNODB_TRX"
}
