package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// 定义请求的负载

func Rapiddns_PostRequest(input string, proxyURL string) (string, error) {
	// 创建请求体

	var transport *http.Transport

	// 根据传入的代理 URL 设置 Transport
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return "", fmt.Errorf("解析代理地址失败: %w", err)
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

	var payload = map[string]string{"full": "1"}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	url := fmt.Sprintf("http://rapiddns.io/subdomain/%s?full=1#result", input)

	// 创建新的 POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求错误:", err)
		return "", nil // 不退出程序，返回空字符串
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码: %s\n", resp.Status)
		return "", nil
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func Rapiddns_ExtractIPFromHTML(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var totalResults []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var subdomain string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					if c.FirstChild != nil && c.FirstChild.Type == html.TextNode {
						if subdomain == "" {
							subdomain = strings.TrimSpace(c.FirstChild.Data)
						}
					}
				}
			}

			// 处理子域名，去掉冒号后的部分
			if strings.Contains(subdomain, ":") {
				subdomain = strings.Split(subdomain, ":")[0]
			}

			totalResults = append(totalResults, subdomain)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return totalResults, nil
}
