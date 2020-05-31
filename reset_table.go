package mysql_transaction

func ResetTable() {

	sql := `
drop table if exists transaction1;
`

	DB.Exec(sql)

	sql = `
create table transaction1
(
    id       int auto_increment
        primary key,
    username varchar(100) default '' not null,
    age      int          default 0  not null
);
`

	DB.Exec(sql)

}
