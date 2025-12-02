package main

import "github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"

// HTTP API 响应类型

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// MCP 相关类型（用于内部转换）

// MCPToolResult MCP 工具结果（内部使用）
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP 内容（内部使用）
type MCPContent struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// FeedDetailRequest Feed详情请求
type FeedDetailRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
}

type SearchFeedsRequest struct {
	Keyword string                   `json:"keyword" binding:"required"`
	Filters xiaohongshu.FilterOption `json:"filters,omitempty"`
}

// FeedDetailResponse Feed详情响应
type FeedDetailResponse struct {
	FeedID string `json:"feed_id"`
	Data   any    `json:"data"`
}

// PostCommentRequest 发表评论请求
type PostCommentRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// PostCommentResponse 发表评论响应
type PostCommentResponse struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserProfileRequest 用户主页请求
type UserProfileRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
}

// ActionResult 通用动作响应（点赞/收藏等）
type ActionResult struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PublishContentRequest 简化的发布文章请求
type PublishContentRequest struct {
	Title   string   `json:"title" binding:"required" example:"我的第一篇小红书"`
	Content string   `json:"content" binding:"required" example:"这是文章的正文内容，可以包含多行文字。"`
	Images  []string `json:"images" binding:"required,min=1" example:"[\"https://example.com/image1.jpg\", \"/path/to/local/image2.png\"]"`
	Tags    []string `json:"tags,omitempty" example:"[\"美食\", \"生活\", \"分享\"]"`
}

// PublishContentResponse 简化的发布文章响应
type PublishContentResponse struct {
	Success bool   `json:"success" example:"true"`
	PostID  string `json:"post_id,omitempty" example:"post_12345"`
	Title   string `json:"title" example:"我的第一篇小红书"`
	Status  string `json:"status" example:"published"`
	Message string `json:"message" example:"文章发布成功"`
}
