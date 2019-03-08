/**
* Created by GoLand.
* User: shuchenchen
* Date: 2019-03-07
* Time: 13:23
*/

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// 测试方法: 输入: set name 1000
func main() {
	address := "127.0.0.1:6379"

	tcp, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		fmt.Println("", err)

		return
	}

	conn, err := net.DialTCP("tcp", nil, tcp)
	if err != nil {
		fmt.Println("", err)

		return
	}

	defer conn.Close()

	//req := SetReq("set name xiaoming1")
	//writeLen, err := conn.Write(req)
	//fmt.Println(writeLen, err)
	//
	//buffer, _ := readLine(conn)
	////conn.Read(buffer)
	////fmt.Println("date", conn, string(buffer), len(buffer), cap(buffer))
	//parsing(buffer)

	br := bufio.NewReader(conn)

	for true {
		var (
			str string
		)

		fmt.Println(address + "> ")

		inputReader := bufio.NewReader(os.Stdin)
		str, err := inputReader.ReadString('\n')

		//fmt.Scanln(&str)
		str = strings.Replace(str, "\n", "", -1)
		fmt.Println("输入:", str, err)

		req := SetReq(str)
		_, err = conn.Write(req)
		//fmt.Println(writeLen, err)

		buffer, _ := readLine(br)
		parsing(br, buffer)
	}

}

//单行回复：回复的第一个字节是 "+"
//错误信息：回复的第一个字节是 "-"
//整形数字：回复的第一个字节是 ":"
//多行字符串：回复的第一个字节是 "$" 第二个是字节长度
//数组：回复的第一个字节是 "*"
// *3\r\n$3\r\nset\r\n$4\r\nname\r\n$8\r\nxiaoming\r\n
func SetReq(strs string) (data []byte) {

	list := strings.Split(strs, " ")

	data = make([]byte, 0)
	data = append(data, []byte(fmt.Sprintf("*%d\r\n", len(list)))...)

	for _, v := range list {
		node := fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
		data = append(data, []byte(node)...)
	}

	//fmt.Printf("SetReq:%v \n", string(bytes.Replace(data, []byte("\r\n"), []byte(" "), -1)))

	return
}

func readLine(br *bufio.Reader) (data []byte, err error) {

	var (
		line = make([]byte, 0)
	)

	for {
		p, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Println("readLine", err)

			return data, err
		}

		n := len(p) - 2
		if n < 0 {
			// 数据都是以 \r\n 结束的
			return data, errors.New("数据不合法")
		}

		// 解决输入字符串中有换行符
		if p[n] != '\r' {
			line = append(line, p[:]...)
			continue
		}

		line = append(line, p[:]...)

		return line, nil
	}
}

func parsing(br *bufio.Reader, buffer []byte) {

	//fmt.Println("parsing", string(buffer))
	list := bytes.Split(buffer, []byte("\r\n"))

	if len(list) == 0 {
		return
	}
	buffer = list[0]

	//fmt.Println("输出结果:", string(buffer))

	switch buffer[0] {
	case '-':
		fmt.Println("错误:", string(buffer[1:]))
	case ':':
		fmt.Println("数字:", string(buffer[1:]))
	case '+':
		fmt.Println("单行:", string(buffer[1:]))
	case '$':
		// 后面的字符串长度
		n, _ := strconv.Atoi(string(buffer[1:]))
		//fmt.Println("多行:", string(buffer[:]), n)

		p := make([]byte, n)
		_, err := io.ReadFull(br, p)
		if err != nil {

			return
		}

		line, _ := readLine(br)
		fmt.Println("多行:", string(line), string(p))

	case '*':
		fmt.Println("数组:", string(buffer[1:]))
		//n, _ := strconv.Atoi(string(buffer[1:]))
		//r := make([]interface{}, 0)
		//for i := range r {
		//
		//
		//}

	default:
		fmt.Println("返回错误:", string(buffer[1:]))
	}

}
