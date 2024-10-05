package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Certificate 结构体用于解析 JSON 数据
type Certificate struct {
	NameValue string `json:"name_value"`
}

// FetchCertificates 发送 GET 请求并返回证书的 name_value 列表
func Crt_FetchCertificates(input string, proxyURL string) ([]string, error) {
	var transport *http.Transport

	// 根据传入的代理 URL 设置 Transport
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("解析代理地址失败: %w", err)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", input)

	// 发送 GET 请求
	resp, err := client.Get(url) // 使用自定义的 client
	if err != nil {
		fmt.Println("请求失败:", err)
		return nil, nil // 返回 nil，继续运行
	}
	defer resp.Body.Close() // 确保请求结束后关闭响应体

	resp.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("收到非 200 响应: %d\n", resp.StatusCode)
		return nil, nil // 返回 nil，继续运行
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 检查响应体是否为空
	if len(body) == 0 {
		return nil, fmt.Errorf("响应体为空")
	}

	// 定义切片来存储解析后的数据
	var certificates []Certificate

	// 解析 JSON 数据
	err = json.Unmarshal(body, &certificates)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 提取 name_value
	var nameValues []string
	for _, cert := range certificates {
		nameValues = append(nameValues, cert.NameValue)
	}

	return nameValues, nil
}
