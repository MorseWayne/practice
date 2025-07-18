```sql
    HASH1 {
       key article:id
       field title
       field content
       field user_id
    }

    ZSET1 {
        key create_time
        score created_at
        member article-id
    }

    ZSET2 {
        key score
        score article_score
        member score:article-id
    }

    SET1 {
        key tags:article-id
        member tag
    }

    SET1 {
        key vote:article-id
        member user_id
    }
```