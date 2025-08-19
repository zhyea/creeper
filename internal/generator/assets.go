package generator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// generateAssets 生成静态资源
func (g *Generator) generateAssets() error {
	// 生成CSS
	if err := g.generateCSS(); err != nil {
		return fmt.Errorf("生成CSS失败: %v", err)
	}

	// 生成JavaScript
	if err := g.generateJS(); err != nil {
		return fmt.Errorf("生成JavaScript失败: %v", err)
	}

	return nil
}

// generateCSS 生成CSS文件
func (g *Generator) generateCSS() error {
	css := fmt.Sprintf(`/* Creeper 小说站点样式 */
:root {
    --primary-color: %s;
    --secondary-color: %s;
    --background-color: %s;
    --text-color: %s;
    --font-family: %s;
    --font-size: %s;
    --line-height: %s;
    --border-color: #e1e5e9;
    --shadow: 0 2px 4px rgba(0,0,0,0.1);
    --shadow-hover: 0 4px 8px rgba(0,0,0,0.15);
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: var(--font-family);
    font-size: var(--font-size);
    line-height: var(--line-height);
    color: var(--text-color);
    background-color: var(--background-color);
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* 头部样式 */
.header {
    background: var(--primary-color);
    color: white;
    padding: 1rem 0;
    box-shadow: var(--shadow);
}

.header .container {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.site-title {
    font-size: 1.5rem;
    font-weight: bold;
}

.site-title a {
    color: white;
    text-decoration: none;
}

.nav {
    display: flex;
    align-items: center;
    gap: 2rem;
}

.nav-link {
    color: white;
    text-decoration: none;
    transition: opacity 0.3s;
}

.nav-link:hover {
    opacity: 0.8;
}

/* 搜索框样式 */
.search-box {
    position: relative;
}

#search-input {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 20px;
    width: 250px;
    font-size: 0.9rem;
}

.search-results {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    max-height: 300px;
    overflow-y: auto;
    z-index: 1000;
    display: none;
}

.search-result-item {
    padding: 0.75rem;
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
    transition: background-color 0.3s;
}

.search-result-item:hover {
    background-color: #f8f9fa;
}

.search-result-item:last-child {
    border-bottom: none;
}

.search-result-title {
    font-weight: bold;
    color: var(--primary-color);
}

.search-result-meta {
    font-size: 0.85rem;
    color: #666;
    margin-top: 0.25rem;
}

/* 主要内容区域 */
.main {
    min-height: calc(100vh - 200px);
    padding: 2rem 0;
}

/* 首页样式 */
.hero {
    text-align: center;
    padding: 3rem 0;
    background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
    color: white;
    border-radius: 8px;
    margin-bottom: 3rem;
}

.hero h2 {
    font-size: 2rem;
    margin-bottom: 1rem;
}

.novels-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
}

.novel-card {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    overflow: hidden;
    transition: transform 0.3s, box-shadow 0.3s;
    box-shadow: var(--shadow);
}

.novel-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-hover);
}

.novel-cover img {
    width: 100%;
    height: 200px;
    object-fit: cover;
}

.novel-info {
    padding: 1.5rem;
}

.novel-title {
    font-size: 1.25rem;
    margin-bottom: 0.5rem;
}

.novel-title a {
    color: var(--primary-color);
    text-decoration: none;
}

.novel-title a:hover {
    text-decoration: underline;
}

.novel-author {
    color: #666;
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
}

.novel-description {
    color: #555;
    font-size: 0.9rem;
    line-height: 1.5;
    margin-bottom: 1rem;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.novel-stats {
    display: flex;
    gap: 1rem;
    font-size: 0.85rem;
    color: #666;
}

/* 小说详情页样式 */
.novel-header {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 2rem;
    margin-bottom: 2rem;
    box-shadow: var(--shadow);
}

.novel-meta {
    display: flex;
    gap: 2rem;
    align-items: flex-start;
}

.novel-cover-large img {
    width: 200px;
    height: 280px;
    object-fit: cover;
    border-radius: 4px;
}

.novel-details {
    flex: 1;
}

.novel-details .novel-title {
    font-size: 2rem;
    margin-bottom: 1rem;
    color: var(--primary-color);
}

.novel-details .novel-author {
    font-size: 1rem;
    margin-bottom: 1rem;
}

.novel-details .novel-description {
    font-size: 1rem;
    line-height: 1.6;
    margin-bottom: 1.5rem;
    color: #555;
}

.novel-details .novel-stats {
    font-size: 0.9rem;
    margin-bottom: 2rem;
}

.novel-actions {
    display: flex;
    gap: 1rem;
}

/* 章节列表样式 */
.chapters-list {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 2rem;
    box-shadow: var(--shadow);
}

.chapters-list h2 {
    margin-bottom: 1.5rem;
    color: var(--primary-color);
}

.chapters-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 0.5rem;
}

.chapter-item {
    border: 1px solid var(--border-color);
    border-radius: 4px;
    overflow: hidden;
    transition: background-color 0.3s;
}

.chapter-item:hover {
    background-color: #f8f9fa;
}

.chapter-link {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    text-decoration: none;
    color: var(--text-color);
}

.chapter-title {
    flex: 1;
    font-weight: 500;
}

.chapter-stats {
    font-size: 0.85rem;
    color: #666;
}

/* 章节阅读页样式 */
.chapter-header {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 2rem;
    margin-bottom: 2rem;
    box-shadow: var(--shadow);
}

.breadcrumb {
    margin-bottom: 1rem;
    font-size: 0.9rem;
    color: #666;
}

.breadcrumb a {
    color: var(--primary-color);
    text-decoration: none;
}

.breadcrumb a:hover {
    text-decoration: underline;
}

.separator {
    margin: 0 0.5rem;
}

.chapter-title {
    font-size: 1.8rem;
    margin-bottom: 1.5rem;
    color: var(--primary-color);
}

.chapter-nav {
    display: flex;
    gap: 1rem;
    justify-content: center;
}

.chapter-content {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 3rem;
    margin-bottom: 2rem;
    box-shadow: var(--shadow);
    line-height: 2;
    font-size: 1.1rem;
}

.chapter-content p {
    margin-bottom: 1.5rem;
    text-indent: 2em;
}

.chapter-content h1,
.chapter-content h2,
.chapter-content h3,
.chapter-content h4,
.chapter-content h5,
.chapter-content h6 {
    margin: 2rem 0 1rem 0;
    color: var(--primary-color);
    text-indent: 0;
}

.chapter-footer {
    background: white;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 2rem;
    box-shadow: var(--shadow);
    text-align: center;
}

.chapter-info {
    margin-bottom: 1.5rem;
    color: #666;
    font-size: 0.9rem;
}

/* 按钮样式 */
.btn {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s;
}

.btn-primary {
    background: var(--primary-color);
    color: white;
}

.btn-primary:hover {
    background: var(--secondary-color);
}

.btn-nav {
    background: #f8f9fa;
    color: var(--text-color);
    border: 1px solid var(--border-color);
}

.btn-nav:hover {
    background: var(--primary-color);
    color: white;
}

/* 底部样式 */
.footer {
    background: var(--primary-color);
    color: white;
    text-align: center;
    padding: 2rem 0;
    margin-top: 3rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .container {
        padding: 0 15px;
    }
    
    .header .container {
        flex-direction: column;
        gap: 1rem;
    }
    
    .nav {
        flex-direction: column;
        gap: 1rem;
    }
    
    #search-input {
        width: 200px;
    }
    
    .novels-grid {
        grid-template-columns: 1fr;
    }
    
    .novel-meta {
        flex-direction: column;
        text-align: center;
    }
    
    .novel-cover-large img {
        width: 150px;
        height: 210px;
    }
    
    .chapters-grid {
        grid-template-columns: 1fr;
    }
    
    .chapter-content {
        padding: 2rem 1.5rem;
        font-size: 1rem;
    }
    
    .chapter-nav {
        flex-wrap: wrap;
        gap: 0.5rem;
    }
    
    .btn {
        font-size: 0.85rem;
        padding: 0.6rem 1.2rem;
    }
}

@media (max-width: 480px) {
    .hero {
        padding: 2rem 0;
    }
    
    .hero h2 {
        font-size: 1.5rem;
    }
    
    .novel-header,
    .chapters-list,
    .chapter-header,
    .chapter-content,
    .chapter-footer {
        padding: 1.5rem;
    }
    
    .chapter-content {
        padding: 2rem 1rem;
    }
}
`,
		g.config.Theme.PrimaryColor,
		g.config.Theme.SecondaryColor,
		g.config.Theme.BackgroundColor,
		g.config.Theme.TextColor,
		g.config.Theme.FontFamily,
		g.config.Theme.FontSize,
		g.config.Theme.LineHeight,
	)

	cssPath := filepath.Join(g.config.OutputDir, "static", "css", "style.css")
	return ioutil.WriteFile(cssPath, []byte(css), 0644)
}

// generateJS 生成JavaScript文件
func (g *Generator) generateJS() error {
	js := `// Creeper 小说站点脚本
(function() {
    'use strict';
    
    let searchData = [];
    let searchTimeout;
    
    // 初始化
    document.addEventListener('DOMContentLoaded', function() {
        initSearch();
        initKeyboardNavigation();
        initReadingProgress();
    });
    
    // 初始化搜索功能
    function initSearch() {
        const searchInput = document.getElementById('search-input');
        const searchResults = document.getElementById('search-results');
        
        if (!searchInput || !searchResults) return;
        
        // 加载搜索数据
        fetch('/static/js/search-data.json')
            .then(response => response.json())
            .then(data => {
                searchData = data;
            })
            .catch(error => {
                console.warn('搜索数据加载失败:', error);
            });
        
        // 搜索输入事件
        searchInput.addEventListener('input', function() {
            const query = this.value.trim();
            
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                if (query.length >= 2) {
                    performSearch(query);
                } else {
                    hideSearchResults();
                }
            }, 300);
        });
        
        // 点击其他地方隐藏搜索结果
        document.addEventListener('click', function(e) {
            if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
                hideSearchResults();
            }
        });
        
        // 键盘导航
        searchInput.addEventListener('keydown', function(e) {
            const items = searchResults.querySelectorAll('.search-result-item');
            const activeItem = searchResults.querySelector('.search-result-item.active');
            
            if (e.key === 'ArrowDown') {
                e.preventDefault();
                if (activeItem) {
                    activeItem.classList.remove('active');
                    const nextItem = activeItem.nextElementSibling;
                    if (nextItem) {
                        nextItem.classList.add('active');
                    } else if (items.length > 0) {
                        items[0].classList.add('active');
                    }
                } else if (items.length > 0) {
                    items[0].classList.add('active');
                }
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                if (activeItem) {
                    activeItem.classList.remove('active');
                    const prevItem = activeItem.previousElementSibling;
                    if (prevItem) {
                        prevItem.classList.add('active');
                    } else if (items.length > 0) {
                        items[items.length - 1].classList.add('active');
                    }
                } else if (items.length > 0) {
                    items[items.length - 1].classList.add('active');
                }
            } else if (e.key === 'Enter') {
                e.preventDefault();
                if (activeItem) {
                    activeItem.click();
                }
            } else if (e.key === 'Escape') {
                hideSearchResults();
                searchInput.blur();
            }
        });
    }
    
    // 执行搜索
    function performSearch(query) {
        const results = searchData.filter(item => {
            const searchText = (item.title + ' ' + (item.author || '') + ' ' + (item.novel || '')).toLowerCase();
            return searchText.includes(query.toLowerCase());
        }).slice(0, 10);
        
        displaySearchResults(results);
    }
    
    // 显示搜索结果
    function displaySearchResults(results) {
        const searchResults = document.getElementById('search-results');
        if (!searchResults) return;
        
        if (results.length === 0) {
            searchResults.innerHTML = '<div class="search-result-item">没有找到相关结果</div>';
        } else {
            searchResults.innerHTML = results.map(item => {
                const typeText = item.type === 'novel' ? '小说' : '章节';
                const metaText = item.type === 'novel' 
                    ? (item.author ? '作者：' + item.author : '')
                    : (item.novel ? '来自：' + item.novel : '');
                
                return ` + "`" + `<div class="search-result-item" onclick="location.href='${item.url}'">
                    <div class="search-result-title">[${typeText}] ${item.title}</div>
                    ${metaText ? ` + "`" + `<div class="search-result-meta">${metaText}</div>` + "`" + ` : ''}
                </div>` + "`" + `;
            }).join('');
        }
        
        searchResults.style.display = 'block';
    }
    
    // 隐藏搜索结果
    function hideSearchResults() {
        const searchResults = document.getElementById('search-results');
        if (searchResults) {
            searchResults.style.display = 'none';
        }
    }
    
    // 初始化键盘导航
    function initKeyboardNavigation() {
        document.addEventListener('keydown', function(e) {
            // 章节页面的键盘导航
            if (document.querySelector('.chapter-content')) {
                if (e.key === 'ArrowLeft' && e.ctrlKey) {
                    e.preventDefault();
                    const prevLink = document.querySelector('a[href*="chapter-"]:nth-of-type(1)');
                    if (prevLink && prevLink.textContent.includes('上一章')) {
                        location.href = prevLink.href;
                    }
                } else if (e.key === 'ArrowRight' && e.ctrlKey) {
                    e.preventDefault();
                    const nextLink = document.querySelector('a[href*="chapter-"]:last-of-type');
                    if (nextLink && nextLink.textContent.includes('下一章')) {
                        location.href = nextLink.href;
                    }
                } else if (e.key === 'ArrowUp' && e.ctrlKey) {
                    e.preventDefault();
                    const tocLink = document.querySelector('a[href="./index.html"]');
                    if (tocLink) {
                        location.href = tocLink.href;
                    }
                }
            }
        });
    }
    
    // 初始化阅读进度
    function initReadingProgress() {
        const chapterContent = document.querySelector('.chapter-content');
        if (!chapterContent) return;
        
        // 创建进度条
        const progressBar = document.createElement('div');
        progressBar.style.cssText = ` + "`" + `
            position: fixed;
            top: 0;
            left: 0;
            width: 0%;
            height: 3px;
            background: var(--primary-color);
            z-index: 9999;
            transition: width 0.3s ease;
        ` + "`" + `;
        document.body.appendChild(progressBar);
        
        // 更新进度
        function updateProgress() {
            const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
            const scrollHeight = document.documentElement.scrollHeight - window.innerHeight;
            const progress = (scrollTop / scrollHeight) * 100;
            progressBar.style.width = Math.min(100, Math.max(0, progress)) + '%';
        }
        
        window.addEventListener('scroll', updateProgress);
        updateProgress();
    }
    
    // 工具函数：平滑滚动
    function smoothScrollTo(element) {
        element.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
        });
    }
    
    // 工具函数：防抖
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
})();`

	jsPath := filepath.Join(g.config.OutputDir, "static", "js", "main.js")
	return ioutil.WriteFile(jsPath, []byte(js), 0644)
}
