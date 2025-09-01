package main

import (
	"creeper/internal/parser"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("🧪 测试 TXT 文件分类和关键字解析功能")

	// 创建解析器
	p := parser.New()

	// 测试文件列表
	testFiles := []string{
		"../novels/TXT示例小说.txt",
		"../novels/TXT科幻小说.txt",
		"../novels/TXT多文件示例",
	}

	for _, file := range testFiles {
		fmt.Printf("\n📖 解析文件: %s\n", file)
		fmt.Println(strings.Repeat("=", 50))

		novel, err := p.ParseNovel(file)
		if err != nil {
			fmt.Printf("❌ 解析失败: %v\n", err)
			continue
		}

		// 显示解析结果
		fmt.Printf("📚 标题: %s\n", novel.Title)
		fmt.Printf("👤 作者: %s\n", novel.Author)
		fmt.Printf("📂 分类: %s\n", novel.Category)
		fmt.Printf("🏷️  标签: %v\n", novel.Tags)
		fmt.Printf("📝 简介: %s\n", novel.Description)
		fmt.Printf("📊 章节数: %d\n", len(novel.Chapters))

		// 显示前几个章节
		if len(novel.Chapters) > 0 {
			fmt.Println("\n📖 章节列表:")
			for i, chapter := range novel.Chapters {
				if i >= 3 { // 只显示前3章
					break
				}
				fmt.Printf("  %d. %s (%d字)\n", chapter.ID, chapter.Title, chapter.WordCount)
			}
			if len(novel.Chapters) > 3 {
				fmt.Printf("  ... 还有 %d 章\n", len(novel.Chapters)-3)
			}
		}
	}

	fmt.Println("\n✅ 测试完成！")
}
