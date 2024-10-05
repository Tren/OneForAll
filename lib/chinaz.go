package lib

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

func Chinaz_FetchContent(input string, proxyURL string) (string, error) {
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

	url := fmt.Sprintf("https://chaziyu.com/%s", input)

	// 创建新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 添加用户代理信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败: %s", resp.Status)
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Chinaz_ExtractAllLinks(htmlContent string) ([]string, error) {
	re := regexp.MustCompile(`rel="nofollow">(.*?)</a></td>`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("未找到匹配内容")
	}

	return results, nil
}
