// Creeper 小说站点脚本
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
                
                return `<div class="search-result-item" onclick="location.href='${item.url}'">
                    <div class="search-result-title">[${typeText}] ${item.title}</div>
                    ${metaText ? `<div class="search-result-meta">${metaText}</div>` : ''}
                </div>`;
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
        progressBar.style.cssText = `
            position: fixed;
            top: 0;
            left: 0;
            width: 0%;
            height: 3px;
            background: var(--primary-color);
            z-index: 9999;
            transition: width 0.3s ease;
        `;
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
})();