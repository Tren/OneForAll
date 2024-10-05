package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FetchHostInfo 发送 GET 请求并返回主机信息
func Hackertarget_FetchHostInfo(input string, proxyURL string) ([]string, error) {
	var transport *http.Transport

	// 根据传入的代理 URL 来设置 Transport
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL) // 这里使用 url.Parse
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

	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", input)

	// 发送 GET 请求
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("请求失败: %w", err)
		return nil, nil
	}
	defer resp.Body.Close() // 确保请求结束后关闭响应体

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("收到非 200 响应: %d", resp.StatusCode)
		return nil, nil
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败: %w", err)
		return nil, nil
	}

	// 将响应体内容按行拆分
	lines := strings.Split(string(body), "\n")

	// 提取主机名
	var hosts []string
	for _, line := range lines {
		if line != "" {
			// 拆分以逗号为分隔符，并取第一个元素
			parts := strings.Split(line, ",")
			hosts = append(hosts, parts[0]) // 只保留主机名
		}
	}

	return hosts, nil
}
