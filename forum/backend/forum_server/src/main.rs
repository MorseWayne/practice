use actix_web::{App, HttpResponse, HttpServer, Responder, get, middleware::Logger, post, web};
use log::info;
use log4rs;
use redis::Client;
use serde::{Deserialize, Serialize};
use std::sync::Arc;

// 新增问题数据模型
#[derive(Debug, Serialize, Deserialize, Clone, Default)]
struct Question {
    title: String,
    tags: Vec<String>,
    details: String,
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

// 共享状态类型定义
// 修改AppState为Redis连接管理器
type AppState = Arc<Client>;

// 新增创建问题API端点
#[post("/api/questions")]
async fn add_question(data: web::Data<AppState>, question: web::Json<Question>) -> impl Responder {
    let conn = match data.get_multiplexed_async_connection().await {
        Ok(conn) => conn,
        Err(e) => return HttpResponse::InternalServerError().body(format!("{}", e)),
    };
    
    // let mut questions = data.lock().unwrap();
    // questions.push(question.into_inner());
    actix_web::HttpResponse::Ok().json("Question added successfully")
}

// 新增查询所有问题API端点
#[get("/api/questions")]
async fn get_questions(data: web::Data<AppState>) -> impl Responder {
    let questions = Question::default();
    actix_web::HttpResponse::Ok().json(&questions)
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
            .service(add_question)
            .service(get_questions)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
