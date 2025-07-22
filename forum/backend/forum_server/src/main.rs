use actix_web::{
    App, HttpResponse, HttpServer, Responder, get, middleware::Logger, post, put, web,
};
use log::info;
use log4rs;
use redis::{AsyncCommands, Client};
use serde::{Deserialize, Serialize};
use std::{fmt::format, sync::Arc};

// 新增问题数据模型
#[derive(Debug, Serialize, Deserialize, Clone, Default)]
struct Article {
    title: String,
    tags: Vec<String>,
    content: String,
    // 用户的id
    poster: u32,
    // 添加新字段并设置默认值
    #[serde(default)]
    view_count: u32,
    #[serde(default)]
    vote_count: i32,
    #[serde(default)]
    answer_count: u32,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
struct VoteRequest {
    article_id: u32,
    user_id: u32,
}

const ARTICLE_KEY: &str = "article";
const CREATE_TIME_KEY: &str = "create_time";
const VOTE_KEY: &str = "vote";
// 共享状态类型定义
// 修改AppState为Redis连接管理器
type AppState = Arc<Client>;

// 新增创建问题API端点
#[post("/api/articles")]
async fn add_article(data: web::Data<AppState>, article: web::Json<Article>) -> impl Responder {
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
    let res: Result<(), redis::RedisError> = redis::cmd("HSET")
        .arg(article_key)
        .arg("title")
        .arg(article.title)
        .arg("content")
        .arg(article.content)
        .exec_async(&mut conn)
        .await;

    match res {
        Ok(_) => actix_web::HttpResponse::Ok().json("Article added successfully"),
        Err(e) => actix_web::HttpResponse::InternalServerError().body(format!("{}", e)),
    }
}

// 新增查询所有问题API端点
#[get("/api/articles")]
async fn get_articles(_data: web::Data<AppState>) -> impl Responder {
    let article = Article::default();
    actix_web::HttpResponse::Ok().json(&article)
}

#[put("/api/articles/{article_id}/vote")]
async fn vote_article(
    data: web::Data<AppState>,
    path: web::Path<u32>,
    vote_req: web::Json<VoteRequest>,
) -> impl Responder {
    let article_id = path.into_inner();
    let vote_req = vote_req.into_inner();
    let mut conn = match data.get_multiplexed_async_connection().await {
        Ok(conn) => conn,
        Err(e) => return HttpResponse::InternalServerError().body(format!("{}", e)),
    };
    let vote_key = format!("{}:{}", VOTE_KEY, article_id);
    let res: Result<(), redis::RedisError> = redis::cmd("SADD")
        .arg(vote_key)
        .arg(vote_req.user_id)
        .exec_async(&mut conn)
        .await;
    match res {
        Ok(_) => actix_web::HttpResponse::Ok().json("Vote added successfully"),
        Err(e) => actix_web::HttpResponse::InternalServerError().body(format!("{}", e)),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    log4rs::init_file("log4rs.yml", Default::default()).expect("Failed to initialize log4rs");

    info!("Server starting");

    // 创建Redis客户端
    let client = Client::open("redis://127.0.0.1:6379/").map_err(|e| {
        std::io::Error::new(
            std::io::ErrorKind::Other,
            format!("Redis client error: {}", e),
        )
    })?;

    // 创建共享状态
    let app_state = Arc::new(client);

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::new(app_state.clone()))
            .service(add_article)
            .service(get_articles)
            .service(vote_article)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
