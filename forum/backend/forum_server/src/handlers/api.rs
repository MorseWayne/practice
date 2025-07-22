use actix_web::{web, HttpResponse, Responder};
use redis::{AsyncCommands, RedisError};
use crate::{constants::*, db::redis::AppState, models::article::Article};

// 新增创建问题API端点
#[post("/api/articles")]
pub async fn add_article(data: web::Data<AppState>, article: web::Json<Article>) -> impl Responder {
    let mut conn = match data.get_multiplexed_async_connection().await {
        Ok(conn) => conn,
        Err(e) => return HttpResponse::InternalServerError().body(format!("{}", e)),
    };

    let article_id: u32 = conn.incr(ARTICLE_KEY, 1).await.unwrap();
    let article_key = format!("article:{}", article_id);
    let time: Vec<String> = redis::cmd("TIME").query_async(&mut conn).await.unwrap();
    redis::cmd("ZADD")
        .arg(CREATE_TIME_KEY)
        .arg(&time[0])
        .arg(&article_key)
        .exec_async(&mut conn)
        .await
        .unwrap();

    let article = article.into_inner();
    let res: Result<(), RedisError> = redis::cmd("HSET")
        .arg(article_key)
        .arg("title")
        .arg(article.title)
        .arg("content")
        .arg(article.content)
        .exec_async(&mut conn)
        .await;

    match res {
        Ok(_) => HttpResponse::Ok().json("Article added successfully"),
        Err(e) => HttpResponse::InternalServerError().body(format!("{}", e)),
    }
}

// 新增查询所有问题API端点
#[get("/api/articles")]
pub async fn get_articles(_data: web::Data<AppState>) -> impl Responder {
    let article = Article::default();
    HttpResponse::Ok().json(&article)
}

#[put("/api/articles/{article_id}/vote")]
pub async fn vote_article(
    data: web::Data<AppState>,
    path: web::Path<u32>,
    vote_req: web::Json<crate::models::article::VoteRequest>,
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