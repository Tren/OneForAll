package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetCurrentPath 返回当前工作目录
func GetCurrentPath() (string, error) {
	return os.Getwd()
}

// GetSubdomains 调用 ksubdomain 工具并返回其输出
func GetSubdomains(domain string) (string, error) {
	// 获取当前路径
	currentDir, err := GetCurrentPath()
	if err != nil {
		return "", err
	}

	// 构建 ksubdomain 的路径
	var ksubdomainPath string
	if runtime.GOOS == "windows" {
		ksubdomainPath = filepath.Join(currentDir, "lib", "ksubdomain.exe") // Windows
	} else {
		ksubdomainPath = filepath.Join(currentDir, "lib", "ksubdomain") // Linux / MacOS
		cmd := exec.Command("chmod", "777", ksubdomainPath)
		fmt.Println(cmd)
	}

	Dir_path := fmt.Sprintf("%s/test.txt", currentDir)

	// 构建命令并添加参数
	cmd := exec.Command(ksubdomainPath, "enum", "-d", domain, "-o", Dir_path)
	fmt.Println(cmd)

	// 捕获输出和错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 输出具体的错误信息
		return "", fmt.Errorf("执行 ksubdomain 时出错: %s, 错误输出: %s", err, string(output))
	}
	return string(output), nil
}

func Ksubdomain_ExtractAndSaveResults() error {
	inputFilePath := "test.txt"
	outputFilePath := "result.txt"

	// 读取 test.txt 的内容
	data, err := os.ReadFile(inputFilePath) // 使用 os.ReadFile
	if err != nil {
		return fmt.Errorf("读取 %s 时出错: %s", inputFilePath, err)
	}

	lines := strings.Split(string(data), "\n")
	var results []string

	// 处理每一行
	for _, line := range lines {
		if strings.Contains(line, "=>") {
			parts := strings.Split(line, "=>")
			// 提取 => 前面的部分
			results = append(results, strings.TrimSpace(parts[0]))
		}
	}

	// 以追加模式写入 result.txt
	f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开 %s 时出错: %s", outputFilePath, err)
	}
	defer f.Close()

	if len(results) > 0 {
		if _, err := f.WriteString(strings.Join(results, "\n") + "\n"); err != nil {
			return fmt.Errorf("写入 %s 时出错: %s", outputFilePath, err)
		}
	}

	// 删除 test.txt
	if err := os.Remove(inputFilePath); err != nil {
		return fmt.Errorf("删除 %s 时出错: %s", inputFilePath, err)
	}

	return nil
}

func RemoveDuplicates(filename string) error {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 map 存储唯一的链接
	uniqueLinks := make(map[string]struct{})

	// 读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		uniqueLinks[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 创建一个新文件以写入去重后的内容
	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 写入去重后的链接
	for link := range uniqueLinks {
		if _, err := outputFile.WriteString(link + "\n"); err != nil {
			return err
		}
	}

	return nil
}
