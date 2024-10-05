package main

import (
	"OneForAll/lib"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	data := flag.String("d", "", "输入内容")
	proxy := flag.String("p", "", "代理地址 (格式: http://your_proxy_address:port)")
	flag.Parse()

	// chinaz
	fmt.Println("chinaz: ")
	Chinaz_htmlContent, err := lib.Chinaz_FetchContent(*data, *proxy)
	if err != nil {
		log.Fatalf("获取内容时出错: %v", err)
	}

	links, err := lib.Chinaz_ExtractAllLinks(Chinaz_htmlContent)
	if err != nil {
		log.Printf("错误: %v", err)
		links = []string{}
	}

	for _, link := range links {
		fmt.Println(link)

		chinaz_resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开 result.txt 时出错: %v", err)
		}
		defer chinaz_resultFile.Close()

		if _, err := chinaz_resultFile.WriteString(link + "\n"); err != nil {
			log.Fatalf("写入 result.txt 时出错: %v", err)
		}
	}

	// rapiddns

	fmt.Println("rapiddns: ")
	Rapiddns_htmlContent, err := lib.Rapiddns_PostRequest(*data, *proxy)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	results, err := lib.Rapiddns_ExtractIPFromHTML(Rapiddns_htmlContent)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)

		rapiddns_resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开 result.txt 时出错: %v", err)
		}
		defer rapiddns_resultFile.Close()

		if _, err := rapiddns_resultFile.WriteString(result + "\n"); err != nil {
			log.Fatalf("写入 result.txt 时出错: %v", err)
		}
	}

	// crt.sh

	fmt.Println("crt: ")
	Crt_nameValues, err := lib.Crt_FetchCertificates(*data, *proxy)
	if err != nil {
		log.Fatalf("错误: %v", err)
	}

	// 打印提取的 name_value
	for _, name := range Crt_nameValues {
		fmt.Println(name)

		crt_resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开 result.txt 时出错: %v", err)
		}
		defer crt_resultFile.Close()

		if _, err := crt_resultFile.WriteString(name + "\n"); err != nil {
			log.Fatalf("写入 result.txt 时出错: %v", err)
		}
	}

	// hackertarget

	fmt.Println("hackertarget: ")
	Hackertarget_hosts, err := lib.Hackertarget_FetchHostInfo(*data, *proxy)
	if err != nil {
		log.Fatalf("错误: %v", err)
	}

	for _, host := range Hackertarget_hosts {
		fmt.Println(host)

		hackertarget_resultFile, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("打开 result.txt 时出错: %v", err)
		}
		defer hackertarget_resultFile.Close()

		if _, err := hackertarget_resultFile.WriteString(host + "\n"); err != nil {
			log.Fatalf("写入 result.txt 时出错: %v", err)
		}

	}

	//ksubdomain

	ksubdomain_subdomain, err := lib.GetSubdomains(*data)
	if err != nil {
		fmt.Printf("执行 ksubdomain 时出错: %s\n", err)
		return
	}

	err = lib.Ksubdomain_ExtractAndSaveResults()
	if err != nil {
		fmt.Printf("提取结果时出错: %s\n", err)
		return
	}

	fmt.Println("结果已成功保存到 result.txt，test.txt 已删除。", ksubdomain_subdomain)

	result_err := lib.RemoveDuplicates("result.txt")
	if result_err != nil {
		log.Fatalf("去重时出错: %v", err)
	}

	log.Println("去重完成，结果已写入 result.txt", result_err)
}
