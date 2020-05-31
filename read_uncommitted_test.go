package mysql_transaction

import (
	"context"
	"database/sql"
	"log"
	"testing"
)

/*
LevelReadUncommitted
没有脏写
*/
func TestReadUncommittedNoDirtyWrite(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `0`})
	DB.Create(&Transaction1{ID: 2, Username: `0`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	log.Println(`tx2 begin`)

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	log.Println(`tx1 begin`)

	////////////////////////////////////////////////////////
	log.Println(`tx1 update 1 start`)
	tx1.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `111`})
	log.Println(`tx1 update 1 end`)

	log.Println(`tx2 update 1 start`)
	tx2.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `222`})
	log.Println(`tx2 update 1 end`)

	log.Println(`tx2 update 2 start`)
	tx2.Model(&Transaction1{ID: 2}).Updates(map[string]interface{}{`username`: `222`})
	log.Println(`tx2 update 2 end`)

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
sql.LevelReadUncommitted
脏读
*/
func TestReadUncommittedDirtyRead1(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `123`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})

	log.Println(`tx2 begin`)

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	//////////////////////////////////////////////////////////c//////////////////////////////////////////////////////

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
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
sql.LevelReadUncommitted
脏读
tx2 读了两次，值不一样
tx1 修改了数据， 还没有提交或者中止。  不管tx1是先开始的，还是后开始的， 都会有脏读。
*/
func TestReadUncommittedDirtyRead2(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `123`})

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	log.Println(`tx1 begin`)

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})

	log.Println(`tx2 begin`)

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	//////////////////////////////////////////////////////////c//////////////////////////////////////////////////////

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
