#!/usr/bin/env python3
"""
Creeper 封面生成器
用于生成自定义的 SVG 格式小说封面
"""

import argparse
import os
from typing import Dict, Any

# 预设主题配置
THEMES = {
    "default": {
        "bg_gradient": ["#2c3e50", "#3498db"],
        "text_color": "#ffffff",
        "accent_color": "#f1c40f",
        "style": "modern"
    },
    "fantasy": {
        "bg_gradient": ["#8e44ad", "#2c3e50", "#1a1a2e"],
        "text_color": "#ffffff",
        "accent_color": "#e74c3c",
        "style": "fantasy"
    },
    "modern": {
        "bg_gradient": ["#667eea", "#764ba2"],
        "text_color": "#ffffff",
        "accent_color": "#ffffff",
        "style": "geometric"
    },
    "classical": {
        "bg_gradient": ["#8b4513", "#a0522d", "#654321"],
        "text_color": "#8b4513",
        "accent_color": "#ffd700",
        "style": "ornate"
    },
    "scifi": {
        "bg_gradient": ["#0a0a23", "#1a1a2e", "#000000"],
        "text_color": "#00ffff",
        "accent_color": "#0080ff",
        "style": "tech"
    }
}

def generate_svg_cover(title: str, subtitle: str = "", theme: str = "default", 
                      width: int = 300, height: int = 400) -> str:
    """生成 SVG 封面"""
    
    theme_config = THEMES.get(theme, THEMES["default"])
    bg_colors = theme_config["bg_gradient"]
    text_color = theme_config["text_color"]
    accent_color = theme_config["accent_color"]
    
    # 创建渐变定义
    if len(bg_colors) == 2:
        gradient = f'''
        <linearGradient id="bgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:{bg_colors[0]};stop-opacity:1" />
            <stop offset="100%" style="stop-color:{bg_colors[1]};stop-opacity:1" />
        </linearGradient>'''
    else:
        gradient = f'''
        <radialGradient id="bgGradient" cx="50%" cy="30%" r="80%">
            <stop offset="0%" style="stop-color:{bg_colors[0]};stop-opacity:1" />
            <stop offset="50%" style="stop-color:{bg_colors[1]};stop-opacity:1" />
            <stop offset="100%" style="stop-color:{bg_colors[2]};stop-opacity:1" />
        </radialGradient>'''
    
    # 根据主题添加装饰元素
    decorations = ""
    if theme == "fantasy":
        decorations = '''
        <!-- 星星装饰 -->
        <g fill="#ffffff" opacity="0.8">
            <circle cx="50" cy="60" r="1"/>
            <circle cx="250" cy="80" r="1.5"/>
            <circle cx="80" cy="320" r="1"/>
        </g>
        <!-- 城堡剪影 -->
        <g fill="#000000" opacity="0.3">
            <rect x="120" y="280" width="60" height="50" rx="5"/>
            <polygon points="140,280 150,260 160,280"/>
        </g>'''
    elif theme == "scifi":
        decorations = '''
        <!-- 星空 -->
        <g fill="#ffffff">
            <circle cx="50" cy="50" r="0.5" opacity="0.8"/>
            <circle cx="250" cy="80" r="1" opacity="0.6"/>
            <circle cx="80" cy="300" r="0.5" opacity="0.9"/>
        </g>
        <!-- 科技线条 -->
        <g stroke="#00ffff" stroke-width="1" fill="none" opacity="0.6">
            <path d="M50 250 L100 270 L150 250 L200 270 L250 250"/>
        </g>'''
    elif theme == "modern":
        decorations = '''
        <!-- 几何装饰 -->
        <g fill="#ffffff" opacity="0.2">
            <circle cx="250" cy="100" r="60"/>
            <rect x="50" y="280" width="40" height="40" transform="rotate(45 70 300)"/>
        </g>'''
    elif theme == "classical":
        decorations = '''
        <!-- 装饰边框 -->
        <rect x="30" y="30" width="240" height="340" fill="none" stroke="#ffd700" stroke-width="2" rx="10"/>
        <!-- 装饰花纹 -->
        <g fill="#ffd700" opacity="0.6">
            <circle cx="150" cy="80" r="20" fill="none" stroke="#ffd700" stroke-width="2"/>
        </g>'''
    
    # 生成 SVG
    svg = f'''<?xml version="1.0" encoding="UTF-8"?>
<svg width="{width}" height="{height}" viewBox="0 0 {width} {height}" xmlns="http://www.w3.org/2000/svg">
    <defs>
        {gradient}
    </defs>
    
    <!-- 背景 -->
    <rect width="{width}" height="{height}" fill="url(#bgGradient)"/>
    
    {decorations}
    
    <!-- 装饰元素 -->
    <g transform="translate({width//2}, {height-50})">
        <circle cx="0" cy="0" r="8" fill="{accent_color}" opacity="0.4"/>
        <circle cx="0" cy="0" r="4" fill="{accent_color}" opacity="0.7"/>
    </g>
</svg>'''
    
    return svg

def main():
    parser = argparse.ArgumentParser(description="生成小说封面")
    parser.add_argument("title", help="小说标题")
    parser.add_argument("-s", "--subtitle", default="", help="副标题")
    parser.add_argument("-t", "--theme", choices=list(THEMES.keys()), 
                       default="default", help="主题风格")
    parser.add_argument("-o", "--output", help="输出文件名")
    parser.add_argument("--width", type=int, default=300, help="宽度")
    parser.add_argument("--height", type=int, default=400, help="高度")
    parser.add_argument("--list-themes", action="store_true", help="列出所有主题")
    
    args = parser.parse_args()
    
    if args.list_themes:
        print("可用主题:")
        for theme_name, theme_config in THEMES.items():
            print(f"  {theme_name}: {theme_config['style']} 风格")
        return
    
    # 生成封面
    svg_content = generate_svg_cover(
        args.title, 
        args.subtitle, 
        args.theme,
        args.width,
        args.height
    )
    
    # 确定输出文件名
    if args.output:
        output_file = args.output
    else:
        safe_title = "".join(c for c in args.title if c.isalnum() or c in (' ', '-', '_')).strip()
        safe_title = safe_title.replace(' ', '-')
        output_file = f"static/images/{safe_title}-cover.svg"
    
    # 确保输出目录存在
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    
    # 写入文件
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(svg_content)
    
    print(f"✅ 封面已生成: {output_file}")
    print(f"📖 标题: {args.title}")
    if args.subtitle:
        print(f"📝 副标题: {args.subtitle}")
    print(f"🎨 主题: {args.theme}")
    print(f"📐 尺寸: {args.width}x{args.height}")

if __name__ == "__main__":
    main()
