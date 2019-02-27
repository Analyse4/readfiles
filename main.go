package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"readfiles/config"
	"readfiles/dao"
	"sync"
)

// 根据实际情况调整一个合适的goroutine数量，假设是4个
const (
	NumWorkers = 4
	GROUP1     = 0
	GROUP2     = 1
	GROUP3     = 2
	GROUP4     = 3
)

func main() {

	err := loadconfig()
	if err != nil {
		log.Println(err)
	}
	dao.InitDataBase()

	group := []int{GROUP1, GROUP2, GROUP3, GROUP4}

	wgRead := &sync.WaitGroup{}

	wgWrite := &sync.WaitGroup{}

	fileChan := make(chan []byte, NumWorkers)

	for i := 0; i < NumWorkers; i++ {
		wgRead.Add(1)
		go readFile(fileChan, group[i], wgRead)
	}

	for i := 0; i < NumWorkers; i++ {
		wgWrite.Add(1)
		go handleFile(fileChan, wgWrite)
	}
	go func() {
		wgRead.Wait()
		close(fileChan)
	}()

	wgWrite.Wait()
	fmt.Println("handle completely")
}

// 这里根据实际的文件排布分为多组，每组包含一部分文件，每个goroutine去处理一组文件，组对应的文件信息可以在数据库中事先映射好
func readFile(fileChan chan<- []byte, groupNum int, wg *sync.WaitGroup) {
	defer wg.Done()
	var rows *sql.Rows
	var err error
	switch groupNum {
	case GROUP1:
		rows, err = dao.DBGroup.Query("SELECT filename FROM files WHERE groups=1")
	case GROUP2:
		rows, err = dao.DBGroup.Query("SELECT filename FROM files WHERE groups=2")
	case GROUP3:
		rows, err = dao.DBGroup.Query("SELECT filename FROM files WHERE groups=3")
	case GROUP4:
		rows, err = dao.DBGroup.Query("SELECT filename FROM files WHERE groups=4")
	default:
		log.Println("undefined group")
		return
	}
	if err != nil {
		log.Println(err)
	}
	var name string
	var files []string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			log.Println(err)
		}
		files = append(files, name)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	// 适用于单文件不大的情况, 单文件过大可以分行读
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
		}
		fileChan <- data
	}
}

func handleFile(fileChan <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case data, ok := <-fileChan:
			if !ok {
				return
			}
			// 具体的处理函数
			handleLogic(data)
		}
	}
}

func handleLogic(data []byte) {
	// 具体的处理逻辑
}

func loadconfig() error{
	err := config.Load(config.DBHost, &dao.Dsn)
	if err != nil {
		return err
	}
	return nil
}
