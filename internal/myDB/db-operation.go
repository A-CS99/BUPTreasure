package myDB

import (
	"BUPTreasure/internal/ApiDTO"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

type SignInfo = ApiDTO.SignInfo
type SignInfoJson = ApiDTO.SignInfoJson

var Db *sql.DB

var createTableSql = `create table if not exists User
    	(
			id       int primary key auto_increment,
			name     varchar(20) not null,
			avatarUrl text not null,
			prize    varchar(10) default '未中奖'
		);`

func InitDB() (err error) {
	dsn := "root:3a2b6c5d@tcp(mysql-container:3306)/BUPTreasure?charset=utf8"
	//dsn := "root:3a2b6c5d@tcp(127.0.0.1:3306)/lantu_lottery?charset=utf8" //本地连接
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("连接失败: ")
		fmt.Println(err)
		return err
	}
	err = Db.Ping()
	if err != nil {
		fmt.Println("Ping连接失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("数据库连接成功")
	Db.SetMaxIdleConns(10)
	Db.SetMaxOpenConns(10)
	// 创建User表
	sqlStr := createTableSql
	_, err = Db.Exec(sqlStr)
	if err != nil {
		fmt.Println("创建User表失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("创建User表成功")
	return nil
}

func FlushTable() (err error) {
	sqlStr := "drop table User"
	_, err = Db.Exec(sqlStr)
	if err != nil {
		fmt.Println("删除User表失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("删除User表成功")
	// 创建User表
	sqlStr = createTableSql
	_, err = Db.Exec(sqlStr)
	if err != nil {
		fmt.Println("创建User表失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("创建User表成功")
	fileName := "/app/buffer/lottery-config-users.json"
	rawArr := []SignInfoJson{}
	jsonArr, err := json.MarshalIndent(rawArr, "", "    ")
	err = ioutil.WriteFile(fileName, jsonArr, 0644)
	if err != nil {
		fmt.Println("写入失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("写入成功")
	return nil
}

func ShowTables() (err error) {
	rows, err := Db.Query("show tables")
	if err != nil {
		fmt.Print("查询失败: ")
		fmt.Println(err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			fmt.Print("获取数据失败: ")
			fmt.Println(err)
			return err
		}
		fmt.Println(tableName)
	}
	return nil
}

func Insert(data SignInfo) (err error) {
	sqlStr := "insert into User(name, avatarUrl, prize) values (?, ?, ?)"
	_, err = Db.Exec(sqlStr, data.Name, data.AvatarUrl, "未中奖")
	if err != nil {
		fmt.Println("插入失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("插入成功")
	sqlStr = "select id, name, avatarUrl from User"
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("查询失败: ")
		fmt.Println(err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	fileName := "/app/buffer/lottery-config-users.json"
	rawArr := []SignInfoJson{}
	for rows.Next() {
		var id int
		var name string
		var avatarUrl string
		err = rows.Scan(&id, &name, &avatarUrl)
		if err != nil {
			fmt.Print("获取数据失败: ")
			fmt.Println(err)
			return err
		}
		rawArr = append(rawArr, SignInfoJson{ID: id, Name: name, Avatar: avatarUrl})
	}
	jsonArr, err := json.MarshalIndent(rawArr, "", "    ")
	err = ioutil.WriteFile(fileName, jsonArr, 0644)
	if err != nil {
		fmt.Println("写入失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("写入成功")
	// 读出文件, 检查是否写入成功
	jsonArr, err = ioutil.ReadFile(fileName)
	fmt.Println(string(jsonArr))
	if err != nil {
		fmt.Println("读取失败: ")
		fmt.Println(err)
		return err
	}
	return nil
}

func RandomQuery(num int, awardType string) (picked []SignInfo, err error) {
	sqlStr := fmt.Sprintf("select id, name, avatarUrl from User where prize = '%s' order by rand() limit ?", "未中奖")
	rows, err := Db.Query(sqlStr, num)
	if err != nil {
		fmt.Println("抽取失败: ")
		fmt.Println(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	for rows.Next() {
		var id int
		var name string
		var avatarUrl string
		err = rows.Scan(&id, &name, &avatarUrl)
		if err != nil {
			fmt.Print("获取数据失败: ")
			fmt.Println(err)
			return nil, err
		}
		_, err2 := Db.Exec("update User set prize = ? where id = ?", awardType, id)
		if err2 != nil {
			fmt.Println("更新失败: ")
			fmt.Println(err2)
			return nil, err2
		}
		picked = append(picked, SignInfo{ID: id, Name: name, AvatarUrl: avatarUrl})
		fmt.Println(id, name, avatarUrl)
	}
	return picked, nil
}

func QueryAll() (users []SignInfo, err error) {
	sqlStr := "select id, name, avatarUrl from User"
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("查询失败: ")
		fmt.Println(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	for rows.Next() {
		var id int
		var name string
		var avatarUrl string
		err = rows.Scan(&id, &name, &avatarUrl)
		if err != nil {
			fmt.Print("获取数据失败: ")
			fmt.Println(err)
			return nil, err
		}
		users = append(users, SignInfo{ID: id, Name: name, AvatarUrl: avatarUrl})
		fmt.Println(id, name, avatarUrl)
	}
	return users, nil
}

func FlushAllAward() (err error) {
	sqlStr := "update User set prize = '未中奖'"
	_, err = Db.Exec(sqlStr)
	if err != nil {
		fmt.Println("更新失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("更新成功")
	return nil
}

func RangeQuery(from int, to int) (avatars []string, err error) {
	sqlStr := "select avatarUrl from User where id >= ? and id < ?"
	rows, err := Db.Query(sqlStr, from, to)
	if err != nil {
		fmt.Println("查询失败: ")
		fmt.Println(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	for rows.Next() {
		var avatarUrl string
		err = rows.Scan(&avatarUrl)
		if err != nil {
			fmt.Print("获取数据失败: ")
			fmt.Println(err)
			return nil, err
		}
		avatars = append(avatars, avatarUrl)
		fmt.Println(avatarUrl)
	}
	return avatars, nil
}

func QueryByName(qname string) (res SignInfo, err error) {
	sqlStr := "select id, name, avatarUrl from User where name = ?"
	rows, err := Db.Query(sqlStr, qname)
	if err != nil {
		fmt.Println("查询失败: ")
		fmt.Println(err)
		return SignInfo{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	rows.Next()
	var id int
	var name string
	var avatarUrl string
	err = rows.Scan(&id, &name, &avatarUrl)
	if err != nil {
		fmt.Print("获取数据失败: ")
		fmt.Println(err)
		return SignInfo{}, err
	}
	res = SignInfo{ID: id, Name: name, AvatarUrl: avatarUrl}
	fmt.Println(id, name, avatarUrl)
	return res, nil
}

func TableSize() (size int, err error) {
	sqlStr := "select count(*) from User"
	rows, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println("查询失败: ")
		fmt.Println(err)
		return 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Print("关闭失败: ")
			fmt.Println(err)
			return
		}
	}(rows)
	for rows.Next() {
		err = rows.Scan(&size)
		if err != nil {
			fmt.Print("获取数据库失败: ")
			fmt.Println(err)
			return 0, err
		}
	}
	return size, nil
}

func Disconnect() (err error) {
	err = Db.Close()
	if err != nil {
		fmt.Println("关闭失败: ")
		fmt.Println(err)
		return err
	}
	fmt.Println("关闭成功")
	return nil
}
