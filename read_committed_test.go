package mysql_transaction

import (
	"context"
	"database/sql"
	"log"
	"testing"
)

/*
LevelReadCommitted
没有脏读

*/
func TestReadCommittedNoDirtyRead(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `123`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	log.Println(`tx2 begin`)

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	//////////////////////////////////////////////////////////c//////////////////////////////////////////////////////

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	log.Println(`tx1 begin`)

	var tran Transaction1
	tran.ID = 1

	m := map[string]interface{}{
		`username`: `333`,
	}

	tx1.Model(&tran).Updates(m)
	log.Println(`tx1 update `)

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	tran2 = Transaction1{}
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

}

/*
LevelReadCommitted
没有脏写
如果第一个事务写第一个key ， 还没提交或回滚， 第二个事务 不让写， 防止脏写
*/
func TestReadCommittedNoDirtyWrite(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `0`})
	DB.Create(&Transaction1{ID: 2, Username: `0`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	log.Println(`tx2 begin`)

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	log.Println(`tx1 begin`)

	////////////////////////////////////////////////////////

	tx1.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `111`})
	log.Println(`tx1 update 1`)

	tx2.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `222`})
	log.Println(`tx2 update 1`)
	tx2.Model(&Transaction1{ID: 2}).Updates(map[string]interface{}{`username`: `222`})
	log.Println(`tx2 update 2`)

	tx2.Commit()

	tx1.Model(&Transaction1{ID: 2}).Updates(map[string]interface{}{`username`: `111`})
	log.Println(`tx1 update 2`)

	tx1.Commit()

	////////////////////////////////////////////////////////

	var beans []Transaction1
	DB.Find(&beans)
	log.Println(`beans `, beans)

}

/*
LevelReadCommitted

读偏斜， 不可重复读。
*/
func TestReadCommittedReadSkew(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `0`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	log.Println(`tx1 begin`)

	tx1.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `111`})
	log.Println(`tx1 update 1`)

	tx1.Commit()
	log.Println(`tx1 Commit `)

	tran2 = Transaction1{}
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

}
