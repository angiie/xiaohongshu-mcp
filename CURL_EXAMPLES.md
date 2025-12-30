# 小红书发布文章 cURL 调用示例

## 基础调用示例

### 1. 简单发布示例

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "我的第一篇小红书",
    "content": "这是我在小红书上的第一篇文章，分享一些生活中的美好时刻。",
    "images": [
      "https://picsum.photos/800/600?random=1"
    ],
    "tags": ["生活", "分享", "第一篇"]
  }'
```

### 2. 多图片发布示例

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "美食探店记录",
    "content": "今天去了一家超棒的餐厅！\n\n环境很好，菜品也很精致，强烈推荐给大家。\n\n📍 地址：市中心美食街\n💰 人均：100-150元",
    "images": [
      "https://picsum.photos/800/600?random=2",
      "https://picsum.photos/800/600?random=3",
      "https://picsum.photos/800/600?random=4"
    ],
    "tags": ["美食", "探店", "推荐", "餐厅"]
  }'
```

### 3. 使用本地图片示例

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "今日穿搭分享",
    "content": "今天的穿搭主题是简约风格，选择了基础款的搭配。\n\n上衣：白色T恤\n下装：牛仔裤\n鞋子：小白鞋\n\n简单却不失时尚感！",
    "images": [
      "/Users/admin/Pictures/outfit1.jpg",
      "/Users/admin/Pictures/outfit2.jpg"
    ],
    "tags": ["穿搭", "简约", "时尚", "OOTD"]
  }'
```

### 4. 无标签发布示例

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "随手拍的风景",
    "content": "路过公园时随手拍的，觉得很美就想分享给大家。",
    "images": [
      "https://picsum.photos/800/600?random=5"
    ]
  }'
```

## 高级调用示例

### 5. 带详细内容的发布

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "护肤心得分享",
    "content": "最近发现了一些很好用的护肤品，想和大家分享一下使用心得。\n\n🌟 产品清单：\n1. 洁面乳 - 温和不刺激\n2. 爽肤水 - 补水效果很好\n3. 精华液 - 美白淡斑\n4. 面霜 - 保湿锁水\n\n💡 使用技巧：\n- 洁面要彻底但不要过度\n- 爽肤水要用化妆棉轻拍\n- 精华液要按摩至吸收\n- 面霜要从内向外涂抹\n\n坚持使用一个月，皮肤状态明显改善！",
    "images": [
      "https://picsum.photos/800/600?random=6",
      "https://picsum.photos/800/600?random=7",
      "https://picsum.photos/800/600?random=8",
      "https://picsum.photos/800/600?random=9"
    ],
    "tags": ["护肤", "美妆", "心得", "分享", "护肤品推荐"]
  }'
```

### 6. 旅行游记发布

```bash
curl -X POST http://localhost:18060/api/v1/content/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "三亚之旅完美收官",
    "content": "为期5天的三亚之旅圆满结束！\n\n🏖️ 行程亮点：\n• Day1: 到达酒店，海边漫步\n• Day2: 天涯海角，南山寺\n• Day3: 蜈支洲岛一日游\n• Day4: 亚龙湾，椰梦长廊\n• Day5: 免税店购物，返程\n\n🌊 最难忘的时刻：\n在蜈支洲岛看日出，海水清澈见底，仿佛置身仙境。\n\n📸 拍照技巧：\n- 黄昏时分光线最美\n- 利用海浪作前景\n- 穿亮色衣服更出片\n\n下次还想再来！",
    "images": [
      "https://picsum.photos/800/600?random=10",
      "https://picsum.photos/800/600?random=11",
      "https://picsum.photos/800/600?random=12"
    ],
    "tags": ["旅行", "三亚", "海边", "度假", "游记"]
  }'
```

## 错误处理示例

### 7. 检查登录状态

```bash
# 发布前先检查登录状态
curl -X GET http://localhost:18060/api/v1/login/status
```

### 8. 获取登录二维码（如果未登录）

```bash
# 如果未登录，获取二维码
curl -X GET http://localhost:18060/api/v1/login/qrcode
```

### 9. 健康检查

```bash
# 检查服务是否正常运行
curl -X GET http://localhost:18060/health
```

## 批量发布脚本示例

### 10. Bash 脚本批量发布

```bash
#!/bin/bash

# 批量发布脚本
BASE_URL="http://localhost:18060/api/v1/content/publish"

# 发布内容数组
declare -a posts=(
  '{"title":"早安分享","content":"新的一天开始了，给大家分享一张美丽的日出图。","images":["https://picsum.photos/800/600?random=20"],"tags":["早安","日出","正能量"]}'
  '{"title":"午餐时光","content":"今天的午餐很丰盛，营养搭配也很均衡。","images":["https://picsum.photos/800/600?random=21"],"tags":["午餐","美食","健康"]}'
  '{"title":"晚安时刻","content":"一天结束了，分享一张温馨的夜景图。","images":["https://picsum.photos/800/600?random=22"],"tags":["晚安","夜景","温馨"]}'
)

# 循环发布
for post in "${posts[@]}"; do
  echo "正在发布: $post"
  
  response=$(curl -s -X POST "$BASE_URL" \
    -H "Content-Type: application/json" \
    -d "$post")
  
  echo "响应: $response"
  echo "---"
  
  # 间隔5秒避免频繁请求
  sleep 5
done

echo "批量发布完成！"
```

## Python 调用示例

### 11. Python requests 示例

```python
import requests
import json

def publish_to_xiaohongshu(title, content, images, tags=None):
    """发布内容到小红书"""
    url = "http://localhost:18060/api/v1/content/publish"
    
    payload = {
        "title": title,
        "content": content,
        "images": images
    }
    
    if tags:
        payload["tags"] = tags
    
    headers = {
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.post(url, json=payload, headers=headers)
        response.raise_for_status()
        
        result = response.json()
        print(f"发布成功: {result}")
        return result
        
    except requests.exceptions.RequestException as e:
        print(f"发布失败: {e}")
        return None

# 使用示例
if __name__ == "__main__":
    publish_to_xiaohongshu(
        title="Python自动发布测试",
        content="这是通过Python脚本自动发布的内容。",
        images=["https://picsum.photos/800/600?random=30"],
        tags=["Python", "自动化", "测试"]
    )
```

## 注意事项

1. **请求频率**: 建议发布间隔至少5秒，避免被平台限制
2. **图片大小**: 建议图片大小不超过5MB
3. **内容审核**: 发布的内容需要符合小红书社区规范
4. **网络超时**: 建议设置合理的超时时间（如30秒）
5. **错误重试**: 建议添加重试机制处理网络异常

## 常见错误及解决方案

| 错误信息 | 可能原因 | 解决方案 |
|----------|----------|----------|
| "请求参数错误" | JSON格式错误或缺少必填字段 | 检查JSON格式和必填参数 |
| "发布失败" | 未登录或网络问题 | 检查登录状态，重试请求 |
| "连接被拒绝" | 服务未启动 | 启动小红书MCP服务 |
| "图片下载失败" | 图片链接无效 | 检查图片URL是否可访问 |