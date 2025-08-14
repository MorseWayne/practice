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
| 参数名 | 类型   | 描述                | 可选值          | 默认值 |
|--------|--------|---------------------|-----------------|--------|
| sort   | string | 排序字段            | time, score     | time   |
| order  | string | 排序方向            | asc, desc       | desc   |
| limit  | number | 最大返回文章数量    | 1-100           | 10     |

### 请求示例
#### 1. 获取最新的10篇文章（默认排序）
```http
GET /api/articles
```

#### 2. 获取评分最高的5篇文章
```http
GET /api/articles?sort=score&order=desc&limit=5
```

#### 3. 获取最早发布的20篇文章
```http
GET /api/articles?sort=time&order=asc&limit=20
```

### 响应体
成功时返回文章对象数组：
```json
[
  {
    "title": "Rust入门指南",
    "content": "本文介绍了Rust编程语言的基础知识...",
    "user_id": 101,
    "tags": ["Rust", "编程", "入门"],
    "view_count": 156,
    "vote_count": 28,
    "answer_count": 5
  },
  {
    "title": "Actix-web框架实战",
    "content": "探讨如何使用Actix-web构建高性能Web应用...",
    "user_id": 102,
    "tags": ["Rust", "Web开发", "Actix"],
    "view_count": 98,
    "vote_count": 15,
    "answer_count": 3
  }
]
```

### 响应状态码
| 状态码 | 描述                |
|--------|---------------------|
| 200    | 请求成功            |
| 400    | 请求参数格式错误    |
| 500    | 服务器内部错误      |

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