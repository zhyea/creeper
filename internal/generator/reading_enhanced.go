package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateEnhancedJS 生成增强的阅读体验 JavaScript
func (g *Generator) generateEnhancedJS() error {
	js := `// Creeper 增强阅读体验脚本
(function() {
    'use strict';
    
    let searchData = [];
    let searchTimeout;
    let readingSettings = {
        theme: 'light',
        fontSize: 16,
        lineHeight: 1.6,
        pageWidth: 800,
        autoScroll: false,
        fullScreen: false
    };
    
    // 初始化
    document.addEventListener('DOMContentLoaded', function() {
        initSearch();
        initKeyboardNavigation();
        initReadingProgress();
        initReadingSettings();
        initThemeSwitcher();
        initAutoScroll();
        initFullScreen();
        loadUserSettings();
    });
    
    // 加载用户设置
    function loadUserSettings() {
        const saved = localStorage.getItem('creeper-reading-settings');
        if (saved) {
            try {
                readingSettings = {...readingSettings, ...JSON.parse(saved)};
                applySettings();
            } catch (e) {
                console.warn('读取用户设置失败:', e);
            }
        }
    }
    
    // 保存用户设置
    function saveUserSettings() {
        localStorage.setItem('creeper-reading-settings', JSON.stringify(readingSettings));
    }
    
    // 应用设置
    function applySettings() {
        const root = document.documentElement;
        
        // 应用主题
        root.setAttribute('data-theme', readingSettings.theme);
        
        // 应用字体设置
        root.style.setProperty('--reading-font-size', readingSettings.fontSize + 'px');
        root.style.setProperty('--reading-line-height', readingSettings.lineHeight);
        root.style.setProperty('--reading-page-width', readingSettings.pageWidth + 'px');
        
        // 更新设置面板
        updateSettingsPanel();
    }
    
    // 初始化主题切换器
    function initThemeSwitcher() {
        // 创建主题切换按钮
        const themeBtn = createToolButton('🌙', '切换主题', toggleTheme);
        addToToolbar(themeBtn);
    }
    
    // 切换主题
    function toggleTheme() {
        const themes = ['light', 'dark', 'sepia', 'green'];
        const currentIndex = themes.indexOf(readingSettings.theme);
        const nextIndex = (currentIndex + 1) % themes.length;
        
        readingSettings.theme = themes[nextIndex];
        applySettings();
        saveUserSettings();
        
        // 更新按钮图标
        const btn = document.querySelector('[data-action="theme"]');
        if (btn) {
            const icons = {'light': '🌙', 'dark': '☀️', 'sepia': '📜', 'green': '🌿'};
            btn.textContent = icons[readingSettings.theme] || '🌙';
        }
    }
    
    // 初始化阅读设置
    function initReadingSettings() {
        const settingsBtn = createToolButton('⚙️', '阅读设置', toggleSettingsPanel);
        addToToolbar(settingsBtn);
        
        createSettingsPanel();
    }
    
    // 创建设置面板
    function createSettingsPanel() {
        const panel = document.createElement('div');
        panel.id = 'settings-panel';
        panel.className = 'settings-panel';
        panel.innerHTML = ` + "`" + `
            <div class="settings-header">
                <h3>阅读设置</h3>
                <button class="close-btn" onclick="toggleSettingsPanel()">×</button>
            </div>
            <div class="settings-content">
                <div class="setting-group">
                    <label>字体大小</label>
                    <div class="font-size-controls">
                        <button onclick="adjustFontSize(-1)">A-</button>
                        <span id="font-size-display">16px</span>
                        <button onclick="adjustFontSize(1)">A+</button>
                    </div>
                </div>
                
                <div class="setting-group">
                    <label>行间距</label>
                    <div class="line-height-controls">
                        <button onclick="adjustLineHeight(-0.1)">-</button>
                        <span id="line-height-display">1.6</span>
                        <button onclick="adjustLineHeight(0.1)">+</button>
                    </div>
                </div>
                
                <div class="setting-group">
                    <label>页面宽度</label>
                    <div class="page-width-controls">
                        <button onclick="adjustPageWidth(-50)">窄</button>
                        <span id="page-width-display">800px</span>
                        <button onclick="adjustPageWidth(50)">宽</button>
                    </div>
                </div>
                
                <div class="setting-group">
                    <label>阅读主题</label>
                    <div class="theme-controls">
                        <button onclick="setTheme('light')" class="theme-btn light">明亮</button>
                        <button onclick="setTheme('dark')" class="theme-btn dark">夜间</button>
                        <button onclick="setTheme('sepia')" class="theme-btn sepia">护眼</button>
                        <button onclick="setTheme('green')" class="theme-btn green">绿色</button>
                    </div>
                </div>
                
                <div class="setting-group">
                    <label>
                        <input type="checkbox" id="auto-scroll" onchange="toggleAutoScroll()">
                        自动滚动
                    </label>
                </div>
                
                <div class="setting-group">
                    <button onclick="resetSettings()" class="reset-btn">恢复默认</button>
                </div>
            </div>
        ` + "`" + `;
        
        document.body.appendChild(panel);
    }
    
    // 切换设置面板
    function toggleSettingsPanel() {
        const panel = document.getElementById('settings-panel');
        if (panel) {
            panel.classList.toggle('active');
        }
    }
    
    // 调整字体大小
    function adjustFontSize(delta) {
        readingSettings.fontSize = Math.max(12, Math.min(24, readingSettings.fontSize + delta));
        applySettings();
        saveUserSettings();
    }
    
    // 调整行间距
    function adjustLineHeight(delta) {
        readingSettings.lineHeight = Math.max(1.2, Math.min(2.5, readingSettings.lineHeight + delta));
        applySettings();
        saveUserSettings();
    }
    
    // 调整页面宽度
    function adjustPageWidth(delta) {
        readingSettings.pageWidth = Math.max(600, Math.min(1200, readingSettings.pageWidth + delta));
        applySettings();
        saveUserSettings();
    }
    
    // 设置主题
    function setTheme(theme) {
        readingSettings.theme = theme;
        applySettings();
        saveUserSettings();
    }
    
    // 更新设置面板显示
    function updateSettingsPanel() {
        const fontSizeDisplay = document.getElementById('font-size-display');
        const lineHeightDisplay = document.getElementById('line-height-display');
        const pageWidthDisplay = document.getElementById('page-width-display');
        const autoScrollCheck = document.getElementById('auto-scroll');
        
        if (fontSizeDisplay) fontSizeDisplay.textContent = readingSettings.fontSize + 'px';
        if (lineHeightDisplay) lineHeightDisplay.textContent = readingSettings.lineHeight.toFixed(1);
        if (pageWidthDisplay) pageWidthDisplay.textContent = readingSettings.pageWidth + 'px';
        if (autoScrollCheck) autoScrollCheck.checked = readingSettings.autoScroll;
        
        // 更新主题按钮状态
        document.querySelectorAll('.theme-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        const activeThemeBtn = document.querySelector(` + "`" + `.theme-btn.${readingSettings.theme}` + "`" + `);
        if (activeThemeBtn) {
            activeThemeBtn.classList.add('active');
        }
    }
    
    // 重置设置
    function resetSettings() {
        readingSettings = {
            theme: 'light',
            fontSize: 16,
            lineHeight: 1.6,
            pageWidth: 800,
            autoScroll: false,
            fullScreen: false
        };
        applySettings();
        saveUserSettings();
    }
    
    // 初始化自动滚动
    function initAutoScroll() {
        const autoScrollBtn = createToolButton('📜', '自动滚动', toggleAutoScroll);
        addToToolbar(autoScrollBtn);
    }
    
    // 切换自动滚动
    function toggleAutoScroll() {
        readingSettings.autoScroll = !readingSettings.autoScroll;
        saveUserSettings();
        
        if (readingSettings.autoScroll) {
            startAutoScroll();
        } else {
            stopAutoScroll();
        }
    }
    
    let autoScrollInterval;
    
    // 开始自动滚动
    function startAutoScroll() {
        stopAutoScroll(); // 确保没有重复的定时器
        autoScrollInterval = setInterval(() => {
            window.scrollBy(0, 1);
            
            // 如果到达页面底部，停止滚动
            if (window.innerHeight + window.pageYOffset >= document.body.offsetHeight - 100) {
                stopAutoScroll();
            }
        }, 50);
    }
    
    // 停止自动滚动
    function stopAutoScroll() {
        if (autoScrollInterval) {
            clearInterval(autoScrollInterval);
            autoScrollInterval = null;
        }
    }
    
    // 初始化全屏模式
    function initFullScreen() {
        const fullScreenBtn = createToolButton('⛶', '全屏阅读', toggleFullScreen);
        addToToolbar(fullScreenBtn);
    }
    
    // 切换全屏模式
    function toggleFullScreen() {
        if (!document.fullscreenElement) {
            document.documentElement.requestFullscreen().then(() => {
                readingSettings.fullScreen = true;
                document.body.classList.add('fullscreen-reading');
            });
        } else {
            document.exitFullscreen().then(() => {
                readingSettings.fullScreen = false;
                document.body.classList.remove('fullscreen-reading');
            });
        }
    }
    
    // 创建工具栏按钮
    function createToolButton(icon, title, onClick) {
        const btn = document.createElement('button');
        btn.className = 'tool-btn';
        btn.textContent = icon;
        btn.title = title;
        btn.onclick = onClick;
        return btn;
    }
    
    // 添加到工具栏
    function addToToolbar(element) {
        let toolbar = document.getElementById('reading-toolbar');
        if (!toolbar) {
            toolbar = createToolbar();
        }
        toolbar.appendChild(element);
    }
    
    // 创建阅读工具栏
    function createToolbar() {
        const toolbar = document.createElement('div');
        toolbar.id = 'reading-toolbar';
        toolbar.className = 'reading-toolbar';
        
        // 只在章节页面显示
        if (document.querySelector('.chapter-content')) {
            document.body.appendChild(toolbar);
        }
        
        return toolbar;
    }
    
    // 初始化搜索功能（保持原有功能）
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
    
    // 初始化键盘导航（增强版）
    function initKeyboardNavigation() {
        document.addEventListener('keydown', function(e) {
            // 如果设置面板打开，不处理导航快捷键
            if (document.getElementById('settings-panel')?.classList.contains('active')) {
                return;
            }
            
            // 章节页面的键盘导航
            if (document.querySelector('.chapter-content')) {
                switch(e.key) {
                    case 'ArrowLeft':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            goToPrevChapter();
                        }
                        break;
                    case 'ArrowRight':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            goToNextChapter();
                        }
                        break;
                    case 'ArrowUp':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            goToToc();
                        }
                        break;
                    case 'f':
                    case 'F':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            toggleFullScreen();
                        }
                        break;
                    case 't':
                    case 'T':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            toggleTheme();
                        }
                        break;
                    case 's':
                    case 'S':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            toggleSettingsPanel();
                        }
                        break;
                    case 'a':
                    case 'A':
                        if (e.ctrlKey) {
                            e.preventDefault();
                            toggleAutoScroll();
                        }
                        break;
                    case 'Escape':
                        hideSearchResults();
                        closeSettingsPanel();
                        break;
                }
            }
        });
    }
    
    // 导航函数
    function goToPrevChapter() {
        const prevLink = document.querySelector('a[href*="chapter-"]:nth-of-type(1)');
        if (prevLink && prevLink.textContent.includes('上一章')) {
            location.href = prevLink.href;
        }
    }
    
    function goToNextChapter() {
        const nextLink = document.querySelector('a[href*="chapter-"]:last-of-type');
        if (nextLink && nextLink.textContent.includes('下一章')) {
            location.href = nextLink.href;
        }
    }
    
    function goToToc() {
        const tocLink = document.querySelector('a[href="./index.html"]');
        if (tocLink) {
            location.href = tocLink.href;
        }
    }
    
    // 关闭设置面板
    function closeSettingsPanel() {
        const panel = document.getElementById('settings-panel');
        if (panel) {
            panel.classList.remove('active');
        }
    }
    
    // 初始化阅读进度（增强版）
    function initReadingProgress() {
        const chapterContent = document.querySelector('.chapter-content');
        if (!chapterContent) return;
        
        // 创建进度条
        const progressContainer = document.createElement('div');
        progressContainer.className = 'progress-container';
        progressContainer.innerHTML = ` + "`" + `
            <div class="progress-bar">
                <div class="progress-fill"></div>
            </div>
            <div class="progress-text">0%</div>
        ` + "`" + `;
        
        document.body.appendChild(progressContainer);
        
        // 更新进度
        function updateProgress() {
            const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
            const scrollHeight = document.documentElement.scrollHeight - window.innerHeight;
            const progress = Math.min(100, Math.max(0, (scrollTop / scrollHeight) * 100));
            
            const progressFill = document.querySelector('.progress-fill');
            const progressText = document.querySelector('.progress-text');
            
            if (progressFill) progressFill.style.width = progress + '%';
            if (progressText) progressText.textContent = Math.round(progress) + '%';
        }
        
        window.addEventListener('scroll', updateProgress);
        updateProgress();
        
        // 添加章节信息
        addChapterInfo();
    }
    
    // 添加章节信息
    function addChapterInfo() {
        const chapterHeader = document.querySelector('.chapter-header');
        if (chapterHeader) {
            const infoDiv = document.createElement('div');
            infoDiv.className = 'chapter-reading-info';
            
            const chapterContent = document.querySelector('.chapter-content');
            const wordCount = chapterContent ? chapterContent.textContent.length : 0;
            const readingTime = Math.ceil(wordCount / 300); // 假设每分钟300字
            
            infoDiv.innerHTML = ` + "`" + `
                <span class="reading-time">预计阅读时间: ${readingTime} 分钟</span>
                <span class="word-count">字数: ${wordCount}</span>
            ` + "`" + `;
            
            chapterHeader.appendChild(infoDiv);
        }
    }
    
    // 暴露全局函数
    window.adjustFontSize = adjustFontSize;
    window.adjustLineHeight = adjustLineHeight;
    window.adjustPageWidth = adjustPageWidth;
    window.setTheme = setTheme;
    window.toggleAutoScroll = toggleAutoScroll;
    window.toggleSettingsPanel = toggleSettingsPanel;
    window.resetSettings = resetSettings;
    
})();`

	jsPath := filepath.Join(g.config.OutputDir, "static", "js", "reading-enhanced.js")
	return os.WriteFile(jsPath, []byte(js), 0644)
}

// generateEnhancedCSS 生成增强的阅读体验 CSS
func (g *Generator) generateEnhancedCSS() error {
	css := fmt.Sprintf(`/* Creeper 增强阅读体验样式 */

/* CSS 变量定义 */
:root {
    /* 原有变量 */
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
    
    /* 阅读体验变量 */
    --reading-font-size: 16px;
    --reading-line-height: 1.6;
    --reading-page-width: 800px;
}

/* 主题样式 */
[data-theme="light"] {
    --theme-bg: #ffffff;
    --theme-text: #333333;
    --theme-secondary: #666666;
    --theme-border: #e1e5e9;
    --theme-card-bg: #ffffff;
}

[data-theme="dark"] {
    --theme-bg: #1a1a1a;
    --theme-text: #e0e0e0;
    --theme-secondary: #b0b0b0;
    --theme-border: #404040;
    --theme-card-bg: #2d2d2d;
}

[data-theme="sepia"] {
    --theme-bg: #f7f3e9;
    --theme-text: #5c4b37;
    --theme-secondary: #8b7355;
    --theme-border: #d4c4a8;
    --theme-card-bg: #faf6ed;
}

[data-theme="green"] {
    --theme-bg: #e8f5e8;
    --theme-text: #2d5016;
    --theme-secondary: #5a7c47;
    --theme-border: #c1d5c1;
    --theme-card-bg: #f0f8f0;
}

/* 应用主题 */
body {
    background-color: var(--theme-bg);
    color: var(--theme-text);
    transition: background-color 0.3s, color 0.3s;
}

/* 增强的章节内容样式 */
.chapter-content {
    max-width: var(--reading-page-width);
    margin: 0 auto;
    background: var(--theme-card-bg);
    border: 1px solid var(--theme-border);
    border-radius: 8px;
    padding: 3rem;
    box-shadow: var(--shadow);
    font-size: var(--reading-font-size);
    line-height: var(--reading-line-height);
    transition: all 0.3s ease;
}

.chapter-content p {
    margin-bottom: 1.5rem;
    text-indent: 2em;
    color: var(--theme-text);
}

/* 阅读工具栏 */
.reading-toolbar {
    position: fixed;
    right: 20px;
    top: 50%;
    transform: translateY(-50%);
    background: var(--theme-card-bg);
    border: 1px solid var(--theme-border);
    border-radius: 25px;
    padding: 10px;
    box-shadow: var(--shadow-hover);
    display: flex;
    flex-direction: column;
    gap: 8px;
    z-index: 1000;
    transition: all 0.3s ease;
}

.reading-toolbar:hover {
    box-shadow: 0 8px 24px rgba(0,0,0,0.2);
}

.tool-btn {
    width: 40px;
    height: 40px;
    border: none;
    border-radius: 50%;
    background: var(--primary-color);
    color: white;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
}

.tool-btn:hover {
    background: var(--secondary-color);
    transform: scale(1.1);
}

/* 设置面板 */
.settings-panel {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 400px;
    max-width: 90vw;
    background: var(--theme-card-bg);
    border: 1px solid var(--theme-border);
    border-radius: 12px;
    box-shadow: 0 20px 40px rgba(0,0,0,0.3);
    z-index: 2000;
    display: none;
    transition: all 0.3s ease;
}

.settings-panel.active {
    display: block;
    animation: fadeInScale 0.3s ease;
}

@keyframes fadeInScale {
    from {
        opacity: 0;
        transform: translate(-50%, -50%) scale(0.9);
    }
    to {
        opacity: 1;
        transform: translate(-50%, -50%) scale(1);
    }
}

.settings-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid var(--theme-border);
}

.settings-header h3 {
    margin: 0;
    color: var(--theme-text);
}

.close-btn {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: var(--theme-secondary);
    width: 30px;
    height: 30px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.close-btn:hover {
    background: var(--theme-border);
}

.settings-content {
    padding: 20px;
}

.setting-group {
    margin-bottom: 20px;
}

.setting-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: var(--theme-text);
}

.font-size-controls,
.line-height-controls,
.page-width-controls {
    display: flex;
    align-items: center;
    gap: 10px;
}

.font-size-controls button,
.line-height-controls button,
.page-width-controls button {
    padding: 8px 12px;
    border: 1px solid var(--theme-border);
    background: var(--theme-card-bg);
    color: var(--theme-text);
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.font-size-controls button:hover,
.line-height-controls button:hover,
.page-width-controls button:hover {
    background: var(--primary-color);
    color: white;
}

.theme-controls {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
}

.theme-btn {
    padding: 12px;
    border: 2px solid var(--theme-border);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 14px;
    font-weight: 500;
}

.theme-btn.light {
    background: #ffffff;
    color: #333333;
}

.theme-btn.dark {
    background: #1a1a1a;
    color: #e0e0e0;
}

.theme-btn.sepia {
    background: #f7f3e9;
    color: #5c4b37;
}

.theme-btn.green {
    background: #e8f5e8;
    color: #2d5016;
}

.theme-btn.active {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.reset-btn {
    width: 100%;
    padding: 12px;
    background: var(--secondary-color);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    transition: background 0.3s ease;
}

.reset-btn:hover {
    background: var(--primary-color);
}

/* 进度条增强 */
.progress-container {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: var(--theme-border);
    z-index: 999;
}

.progress-bar {
    height: 100%;
    background: var(--theme-border);
    position: relative;
}

.progress-fill {
    height: 100%;
    background: var(--primary-color);
    transition: width 0.3s ease;
    position: relative;
}

.progress-text {
    position: fixed;
    top: 10px;
    right: 20px;
    background: var(--theme-card-bg);
    color: var(--theme-text);
    padding: 4px 8px;
    border-radius: 12px;
    font-size: 12px;
    border: 1px solid var(--theme-border);
    z-index: 1001;
}

/* 章节阅读信息 */
.chapter-reading-info {
    margin-top: 15px;
    padding: 10px;
    background: var(--theme-bg);
    border: 1px solid var(--theme-border);
    border-radius: 6px;
    display: flex;
    justify-content: space-between;
    font-size: 14px;
    color: var(--theme-secondary);
}

/* 全屏阅读模式 */
.fullscreen-reading .header,
.fullscreen-reading .footer {
    display: none;
}

.fullscreen-reading .main {
    padding: 20px 0;
}

.fullscreen-reading .chapter-content {
    max-width: 90vw;
    margin: 0 auto;
}

/* 移动端优化 */
@media (max-width: 768px) {
    .reading-toolbar {
        right: 10px;
        padding: 8px;
    }
    
    .tool-btn {
        width: 35px;
        height: 35px;
        font-size: 14px;
    }
    
    .settings-panel {
        width: 95vw;
        max-height: 80vh;
        overflow-y: auto;
    }
    
    .chapter-content {
        padding: 2rem 1.5rem;
        font-size: var(--reading-font-size);
    }
    
    .theme-controls {
        grid-template-columns: 1fr;
    }
    
    .chapter-reading-info {
        flex-direction: column;
        gap: 5px;
        text-align: center;
    }
}

@media (max-width: 480px) {
    .chapter-content {
        padding: 1.5rem 1rem;
    }
    
    .reading-toolbar {
        right: 5px;
        bottom: 20px;
        top: auto;
        transform: none;
        flex-direction: row;
        border-radius: 20px;
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

	cssPath := filepath.Join(g.config.OutputDir, "static", "css", "reading-enhanced.css")
	return os.WriteFile(cssPath, []byte(css), 0644)
}
