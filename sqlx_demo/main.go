package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return
}

type user struct {
	ID   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
//func BatchInsertUsers2(users []interface{}) error {
//
//}

// 查询单条数据
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err: %v\n", err)
		return
	}
	fmt.Printf("id:%d name: %s age:%d\n", u.ID, u.Name, u.Age)
}

// queryMultiRowDemo 查询多条数据
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err: %v\n", err)
		return
	}
	fmt.Printf("users: %#v\n", users)
}

// insertDemo 插入数据
func insertDemo() {
	sqlStr := "insert into user(name,age) values (?, ?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err: %v\n", err)
		return
	}
	theId, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theId)
}

// updateRowDemo 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 5)
	if err != nil {
		fmt.Printf("update failed, err: %v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err: %v\n", err)
		return
	}
	fmt.Printf("update success,affected rows:%d\n", n)
}

// deleteRowDemo 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id=?"
	ret, err := db.Exec(sqlStr, 5)
	if err != nil {
		fmt.Printf("delete failed, err%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// insertUserDemo 使用 NamedExec 处理map同名字段
func insertUserDemo() (err error) {
	sqlStr := "INSERT INTO user (name, age) values (:name, :age)"
	_, err = db.NamedExec(sqlStr, map[string]interface{}{
		"name": "sunSmile",
		"age":  30,
	})
	return
}

func namedQuery() {
	sqlStr := "select * from user where name=:name"
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "sunSmile"})
	if err != nil {
		fmt.Printf("db.NameQuery failed, err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := user{
		Name: "sunSmile",
	}
	// 使用结构体命名查询，根据结构体字段的 db tag 进行映射
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.Namedquery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

// 事务操作
func transactionDemo2() (err error) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "update user set age=29 where id=?"
	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "update user set age=50 where id = ?"
	rs, err = tx.Exec(sqlStr2, 6)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return err
}

// QueryByIds 根据给定ID查询
func QueryByIds(ids []int) (users []user, err error) {
	// 填充id
	query, args, err := sqlx.In("select name, age from user where id in (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带"？" bindvar 的查询语句， 我们使用 Rebind()重新绑定它
	db.Rebind(query)
	db.Select(&users, query, args...)
	return
}

func BatchInsertUsers(users []*user) error {
	// 存放（？，？）的 slice
	valueStrings := make([]string, 0, len(users))
	// 存放 values 的 slice
	valueArgs := make([]interface{}, 0, len(users)*2)

	// 遍历users 准备的数据
	for _, u := range users {
		valueStrings = append(valueStrings, "(?,?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO user (name, age) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

// BatchInsertUser2 使用sqlx.In实现批量插入
// 需要我们的结构体实现driver.Valuer接口
//BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUser2(users []interface{}) error {
	query, args, _ := sqlx.In("insert into user (name, age) values (?),(?),(?)", users...)
	fmt.Println(query)
	fmt.Println(args)
	_, err := db.Exec(query, args...)
	return err
}

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*user) error {
	_, err := db.NamedExec("insert into user (name, age) values (:name, :age)", users)
	return err
}

// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIDs(ids []int) (users []user, err error) {
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("select name, age from user where id in (?) order by FIND_IN_SET(id,?)", ids, strings.Join(strIDs, ","))
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed, err: %v\n", err)
		return
	}
	fmt.Println("init DB success")
	defer db.Close()
	//queryRowDemo()
	//queryMultiRowDemo()
	//insertDemo()
	//updateRowDemo()
	//deleteRowDemo()
	//insertUserDemo()
	//namedQuery()
	//transactionDemo2()

	//u1 := user{Name: "sunsmile", Age: 18}
	//u2 := user{Name: "daniel", Age: 28}
	//u3 := user{Name: "小佑宝", Age: 38}
	// 方法1
	//users := []*user{&u1, &u2, &u3}
	//err = BatchInsertUsers(users)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
	//}

	// 方法2
	//user2 := []interface{}{u1, u2, u3}
	//err = BatchInsertUser2(user2)
	//if err != nil {
	//	fmt.Printf("BatchInsertUser2 failed, err:%v\n", err)
	//}
	// 方法3
	//users3 := []*user{&u1, &u2, &u3}
	//err = BatchInsertUsers3(users3)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
	//}

	//users, err := QueryByIds([]int{4, 6, 1})
	//if err != nil {
	//	fmt.Printf("QueryByIDs failed, err:%v\n", err)
	//	return
	//}
	//for _, user := range users {
	//	fmt.Printf("user:%#v\n", user)
	//}
	// 1. 用代码去做排序
	// 2. 让MySQL排序
	fmt.Println("----")
	users, err := QueryAndOrderByIDs([]int{7, 4, 6, 1})
	if err != nil {
		fmt.Printf("QueryByIDs failed, err:%v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("user:%#v\n", user)
	}
}
