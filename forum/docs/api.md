# 论坛系统 REST API

## 创建文章

### HTTP请求
`POST /api/articles`

### 请求体
包含文章信息的JSON对象，具体字段如下：

| 字段名 | 类型 | 描述 | 是否必需 |
|--------|------|------|----------|
| title | string | 文章标题 | 是 |
| content | string | 文章内容 | 是 |
| user_id | integer | 作者ID | 是 |
| tags | array | 标签列表 | 否 |
| view_count | integer | 浏览数 | 否，默认为0 |
| vote_count | integer | 投票数 | 否，默认为0 |
| answer_count | integer | 回答数 | 否，默认为0 |

### 响应体
成功时返回包含文章ID的JSON字符串：
```json
"Article added successfully with id: {article_id}"
```

### 错误响应
| 状态码 | 描述 |
|--------|------|
| 500 | 服务器内部错误 |

## 获取文章列表

### HTTP请求
`GET /api/articles`

### 查询参数
| 参数名 | 类型 | 描述 | 可选值 | 默认值 |
|--------|------|------|--------|--------|
| sort | string | 排序字段 | time, score | time |
| order | string | 排序方向 | asc, desc | desc |
| limit | integer | 最大返回数量 | 1-100 | 10 |

### 响应体
成功时返回文章对象数组：
```json
[
  {
    "title": "文章标题",
    "content": "文章内容",
    "user_id": 123,
    "tags": ["标签1", "标签2"],
    "view_count": 0,
    "vote_count": 0,
    "answer_count": 0
  }
]
```

## 文章投票

### HTTP请求
`PUT /api/articles/{article_id}/vote`

### 路径参数
| 参数名 | 类型 | 描述 |
|--------|------|------|
| article_id | integer | 文章ID |

### 请求体
包含投票用户ID的JSON对象：
```json
{
  "user_id": 123
}
```

### 响应体
成功时返回操作结果：
```json
"operation successful."
```

### 错误响应
| 状态码 | 描述 |
|--------|------|
| 400 | 已投票或参数错误 |