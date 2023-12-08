package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	// 数值需要业务具体情况来确定
	//db.SetConnMaxLifetime(time.Second*10)
	db.SetMaxOpenConns(100) // 最大连接数
	db.SetMaxIdleConns(10)  // 最大空闲连接数
	return nil
}

type user struct {
	id   int
	age  int
	name string
}

func queryRowDemo() {
	sqlStr := "select id,name,age from user where id=?"
	var u user
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err %v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?, ?)"
	ret, err := db.Exec(sqlStr, "zhangsan", 18)
	if err != nil {
		fmt.Printf("insert failed:%v\n", err)
	}
	theId, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get last insertId failed:%v\n", err)
	}
	fmt.Printf("insert success, the id is %d.\n", theId)
}

func queryRowsDemo() {
	sqlStr := "select id,name,age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed:%v\n", err)
		return
	}
	defer rows.Close() // 关闭rows释放持有的数据库链接
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d, name: %s, age: %d\n", u.id, u.name, u.age)
	}
}

func updateDemo() {
	sqlStr := "update user set age=? where id=?"
	ret, err := db.Exec(sqlStr, 19, 1)
	if err != nil {
		fmt.Printf("update failed, err: %v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

func deleteDemo() {
	sqlStr := "delete form user where id =?"
	ret, err := db.Exec(sqlStr, 1)
	if err != nil {
		fmt.Printf("delete failed: %v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err: %v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err : %v\n", err)
	}
	_, err = stmt.Exec("张三", 20)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	_, err = stmt.Exec("李四", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}

func transactionDemo() {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback()
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "update user set age=25 where id=?"
	ret2, err := tx.Exec(sqlStr2, 1)
	if err != nil {
		tx.Rollback()
		fmt.Printf("exec sql2 faild, err: %v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("exec ret2.RowsAffected() failed, err:%v\n", err)
	}
	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}
	fmt.Println("exec trans success!")
}
func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	defer db.Close()
	//queryRowDemo()
	//insertRowDemo()
	//queryRowsDemo()
	//updateDemo()
	//prepareQueryDemo()
	//prepareInsertDemo()
	transactionDemo()
}
