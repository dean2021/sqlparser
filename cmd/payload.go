package main

import (
	"bufio"
	"os"
	"path/filepath"
)

func Load(dirname string) []string {

	var payloads []string
	// 使用filepath.Walk遍历目录
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为.txt文件
		if info.IsDir() == false && filepath.Ext(path) == ".txt" {
			// 打开文件
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// 创建bufio.Reader以逐行读取
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// 如果行以"#"开头，则跳过
				if len(line) > 0 && line[0] == '#' {
					continue
				}

				// 打印非注释行
				payloads = append(payloads, line)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return payloads
}
