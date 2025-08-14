use crate::{constants::*, db::redis::AppState, models::request::{Article, ArticleQuery}};
use actix_web::{HttpResponse, Responder, get, post, put, web};
use redis::RedisError;

// 新增创建问题API端点
#[post("/api/articles")]
pub async fn add_article(data: web::Data<AppState>, article: web::Json<Article>) -> impl Responder {
    let mut conn = match data.get_multiplexed_async_connection().await {
        Ok(conn) => conn,
        Err(e) => return HttpResponse::InternalServerError().body(format!("{}", e)),
    };

    let article = article.into_inner();
    // 加载外部Lua脚本并进行语法检查
    let script_content = include_str!("../scripts/add_article.lua");
    let script = redis::Script::new(script_content);
    // 执行脚本并处理详细错误
    let tags_json = match serde_json::to_string(&article.tags) {
        Ok(json) => json,
        Err(e) => {
            return HttpResponse::InternalServerError().body(format!("serialize tags error: {}", e));
        }
    };

    let article_id: u32 = match script
        .key(ARTICLE_KEY)
        .key(CREATE_TIME_KEY)
        .arg(article.title)
        .arg(article.content)
        .arg(article.user_id)
        .arg(tags_json)
        .invoke_async(&mut conn)
        .await
    {
        Ok(id) => id,
        Err(e) => {
            log::error!("Failed to execute Lua script: {:?}", e);
            return HttpResponse::InternalServerError()
                .body(format!("Failed to add article: {}", e));
        }
    };
    HttpResponse::Ok().json(format!(
        "Article added successfully with id: {}",
        article_id
    ))
}

// 新增查询所有问题API端点
#[get("/api/articles")]
pub async fn get_articles(_data: web::Data<AppState>, query: web::Query<ArticleQuery>) -> impl Responder {
    let article = Article::default();
    HttpResponse::Ok().json(&article)
}

#[put("/api/articles/{article_id}/vote")]
pub async fn vote_article(
    data: web::Data<AppState>,
    path: web::Path<u32>,
    vote_req: web::Json<crate::models::request::Vote>,
) -> impl Responder {
    let article_id = path.into_inner();
    let vote_req = vote_req.into_inner();
    let mut conn = match data.get_multiplexed_async_connection().await {
        Ok(conn) => conn,
        Err(e) => return HttpResponse::InternalServerError().body(format!("{}", e)),
    };
    let vote_key = format!("{}:{}", VOTE_KEY, article_id);
    let res: Result<u64, RedisError> = redis::cmd("SADD")
        .arg(vote_key)
        .arg(vote_req.user_id)
        .query_async(&mut conn)
        .await;

    match res {
        Ok(0) => HttpResponse::Ok().json("already voted."),
        Ok(_) => {
            redis::cmd("ZINCRBY")
                .arg(SCORE_KEY)
                .arg(VOTE_SCORE)
                .arg(article_id)
                .exec_async(&mut conn)
                .await
                .unwrap();
            HttpResponse::Ok().json("operation successful.")
        }
        Err(e) => HttpResponse::InternalServerError().body(format!("{}", e)),
    }
}
