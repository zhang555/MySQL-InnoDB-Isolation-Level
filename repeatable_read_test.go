package mysql_transaction

import (
	"context"
	"database/sql"
	"log"
	"testing"
)

/*

mysql innodb
可重复读 隔离级别
*/

/*
LevelRepeatableRead
没有 读偏斜
*/
func TestLevelRepeatableReadNoReadSkew(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1, Username: `0`})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	tx1.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`username`: `111`})
	log.Println(`tx1 update 1`)

	tx1.Commit()

	tran2 = Transaction1{}
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	var tran3 Transaction1
	DB.First(&tran3, 1)
	log.Println(`tx3 read.  tran3: `, tran3)

	log.Println(`mysql innodb LevelRepeatableRead 没有 读偏斜 `)

}

/*
LevelRepeatableRead
丢失更新 lost update
*/
func TestLevelRepeatableReadLostUpdate(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////
	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	var tran1 Transaction1
	tx1.First(&tran1, 1)
	log.Println(`tx1 read.  tran1: `, tran1)

	tx1.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`age`: tran1.Age + 100})

	log.Println(`tx1 update +100`)

	////////////////////////////////////////////////////////

	var tran2 Transaction1
	tx2.First(&tran2, 1)
	log.Println(`tx2 read.  tran2: `, tran2)

	tx1.Commit()

	tx2.Model(&Transaction1{ID: 1}).Updates(map[string]interface{}{`age`: tran1.Age + 100})

	log.Println(`tx2 update +100 `)

	tx2.Commit()

	var tran3 Transaction1
	DB.First(&tran3, 1)
	log.Println(`tx3 read.  tran3: `, tran3)

	log.Println(`mysql innodb LevelRepeatableRead 丢失更新 lost update`)

}

/*
LevelRepeatableRead
使用特殊语法 防止 更新丢失

*/
func TestLevelRepeatableReadPreventLostUpdate(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////
	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	tx1.Exec(`update transaction1 set age = age + 100 where id = 1 `)

	log.Println(`tx1 update +100`)

	tx1.Commit()

	////////////////////////////////////////////////////////

	tx2.Exec(`update transaction1 set age = age + 100 where id = 1 `)
	log.Println(`tx2 update +100`)

	tx2.Commit()

	var tran3 Transaction1
	DB.First(&tran3, 1)
	log.Println(`tx3 read.  tran3: `, tran3)

	log.Println(`mysql innodb LevelRepeatableRead 使用特殊语法 防止 更新丢失 `)

}

/*
LevelRepeatableRead
写偏斜：读多个对象，更新其中一个对象

*/
func TestLevelRepeatableReadWriteSkew(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1})
	DB.Create(&Transaction1{ID: 2})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////
	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	sql1 := `
select count(*) c
from transaction1
where age = 0;
`

	type C struct {
		C int
	}
	var c C
	tx1.Raw(sql1).Scan(&c)
	log.Println(`tx1 c `, c)
	if c.C == 2 {
		tx1.Exec(`update transaction1 set age = 1 where id = 1 `)
		log.Println(`tx1 set `)

	}

	c = C{}
	tx2.Raw(sql1).Scan(&c)
	log.Println(`tx2 c `, c)

	tx1.Commit()

	if c.C == 2 {
		tx2.Exec(`update transaction1 set age = 1 where id = 2 `)
		log.Println(`tx2 set `)

	}

	tx2.Commit()

	////////////////////////////////////////////////////////

	trans := make([]Transaction1, 0)
	DB.Find(&trans)
	log.Println(`trans `, trans)

	log.Println(`mysql innodb LevelRepeatableRead 写偏斜 `)

}

/*
LevelRepeatableRead
写偏斜：读多个对象，更新其中一个对象
for update 防止写偏斜 ,
*/
func TestLevelRepeatableReadPreventWriteSkew(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1})
	DB.Create(&Transaction1{ID: 2})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////
	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	sql1 := `
select count(*) c
from transaction1
where age = 0 for
update;
`

	type C struct {
		C int
	}
	var c C
	tx1.Raw(sql1).Scan(&c)
	log.Println(`tx1 c `, c)

	if c.C == 2 {
		tx1.Exec(`update transaction1 set age = 1 where id = 1 `)
		log.Println(`tx1 set end`)
	}

	log.Println(`tx2 start read . tx2 read wouldn't success `)

	c = C{}
	tx2.Raw(sql1).Scan(&c)
	log.Println(`tx2 c `, c)

	if c.C == 2 {
		tx2.Exec(`update transaction1 set age = 1 where id = 2 `)
		log.Println(`tx2 set end`)
	}

	tx1.Commit()

	tx2.Commit()

	////////////////////////////////////////////////////////

	trans := make([]Transaction1, 0)
	DB.Find(&trans)
	log.Println(`trans `, trans)

	log.Println(`mysql innodb LevelRepeatableRead for update 防止写偏斜  `)

}

/*
LevelRepeatableRead
写偏斜：读多个对象，更新其中一个对象
select * for update 防止写偏斜


*/
func TestLevelRepeatableReadPreventWriteSkew2(t *testing.T) {

	DB.Delete(&Transaction1{})
	DB.Create(&Transaction1{ID: 1})
	DB.Create(&Transaction1{ID: 2})

	////////////////////////////////////////////////////////
	tx2 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx2 begin`)

	////////////////////////////////////////////////////////
	tx1 := DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	log.Println(`tx1 begin`)

	sql1 := `
select *
from transaction1
where age = 0 for
update;
`

	trans1 := make([]Transaction1, 0)

	tx1.Raw(sql1).Scan(&trans1)
	log.Println(`tx1 trans1 `, trans1)

	if len(trans1) >= 2 {
		tx1.Exec(`update transaction1 set age = 1 where id = 1 `)
		log.Println(`tx1 set `)
	}

	trans2 := make([]Transaction1, 0)
	log.Println(`tx2 read start  `)
	tx2.Raw(sql1).Scan(&trans2)
	log.Println(`tx2 read end . trans2: `, trans2)

	tx1.Commit()

	if len(trans2) >= 2 {
		tx2.Exec(`update transaction1 set age = 1 where id = 2 `)
		log.Println(`tx2 set `)
	}

	tx2.Commit()

	////////////////////////////////////////////////////////

	trans := make([]Transaction1, 0)
	DB.Find(&trans)
	log.Println(`trans `, trans)

}
