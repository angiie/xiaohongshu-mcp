# 小红书发布文章 HTTP API 接口文档

## 接口概述

本文档描述了小红书MCP服务提供的HTTP发布文章接口，支持通过HTTP请求直接发布图文内容到小红书平台。

## 基础信息

- **服务地址**: `http://localhost:18060` (默认端口)
- **接口路径**: `/api/v1/content/publish`
- **请求方法**: `POST`
- **内容类型**: `application/json`

## 接口详情

### 发布文章接口

#### 请求地址
```
POST /api/v1/content/publish
```

#### 请求头
```
Content-Type: application/json
```

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 | 示例 |
|--------|------|------|------|------|
| title | string | 是 | 文章标题，最多20个中文字符 | "我的第一篇小红书" |
| content | string | 是 | 文章正文内容 | "这是文章的正文内容，可以包含多行文字。" |
| images | array | 是 | 图片路径数组，至少1张图片 | ["https://example.com/image1.jpg", "/path/to/local/image2.png"] |
| tags | array | 否 | 话题标签数组 | ["美食", "生活", "分享"] |

#### 请求示例

```json
{
  "title": "我的第一篇小红书",
  "content": "这是文章的正文内容，可以包含多行文字。\n\n支持换行和各种文字描述。",
  "images": [
    "https://example.com/image1.jpg",
    "/Users/admin/Pictures/image2.png"
  ],
  "tags": ["美食", "生活", "分享"]
}
```

#### 响应格式

##### 成功响应 (200 OK)

```json
{
  "success": true,
  "data": {
    "success": true,
    "post_id": "post_12345",
    "title": "我的第一篇小红书",
    "status": "published",
    "message": "文章发布成功"
  },
  "message": "文章发布成功"
}
```

##### 错误响应 (400/500)

```json
{
  "error": "发布失败",
  "code": "PUBLISH_FAILED",
  "details": "具体错误信息"
}
```

#### 响应参数说明

| 参数名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 请求是否成功 |
| data.success | boolean | 发布是否成功 |
| data.post_id | string | 发布后的文章ID |
| data.title | string | 文章标题 |
| data.status | string | 发布状态 |
| data.message | string | 发布结果消息 |

## 图片支持

接口支持两种图片格式：

1. **HTTP/HTTPS 图片链接**: 系统会自动下载网络图片
   - 示例: `"https://example.com/image.jpg"`

2. **本地图片绝对路径**: 直接使用本地图片文件
   - 示例: `"/Users/admin/Pictures/image.png"`

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| INVALID_REQUEST | 请求参数错误 |
| PUBLISH_FAILED | 发布失败 |
| STATUS_CHECK_FAILED | 登录状态检查失败 |

## 使用前提

1. 确保小红书MCP服务已启动
2. 确保已完成小红书账号登录（可通过 `/api/v1/login/status` 检查登录状态）
3. 如果未登录，需要先通过 `/api/v1/login/qrcode` 获取二维码进行登录

## 注意事项

1. 图片数量：至少需要1张图片，建议不超过9张
2. 标题长度：小红书限制标题最多20个中文字符或英文单词
3. 内容格式：支持换行符，建议内容丰富且有吸引力
4. 标签使用：合理使用标签有助于内容被更多用户发现
5. 图片质量：建议使用高质量图片，支持常见格式（JPG、PNG等）

## 健康检查

可以通过以下接口检查服务状态：

```
GET /health
```

响应：
```json
{
  "status": "ok",
  "message": "服务运行正常",
  "version": "1.0.0"
}
```