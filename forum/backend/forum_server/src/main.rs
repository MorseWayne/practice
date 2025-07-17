use actix_files::Files;
use actix_web::{App, HttpServer, Responder, get, middleware::Logger, post, web};
use log::info;
use log4rs;
use serde::{Deserialize, Serialize}; // 新增导入
use std::sync::{Arc, Mutex}; // 新增导入

// 新增问题数据模型
#[derive(Debug, Serialize, Deserialize, Clone)]
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
type AppState = Arc<Mutex<Vec<Question>>>;

// 新增创建问题API端点
#[post("/api/questions")]
async fn add_question(data: web::Data<AppState>, question: web::Json<Question>) -> impl Responder {
    let mut questions = data.lock().unwrap();
    questions.push(question.into_inner());
    actix_web::HttpResponse::Ok().json("Question added successfully")
}

// 新增查询所有问题API端点
#[get("/api/questions")]
async fn get_questions(data: web::Data<AppState>) -> impl Responder {
    let questions = data.lock().unwrap();
    actix_web::HttpResponse::Ok().json(&*questions)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    log4rs::init_file("log4rs.yml", Default::default()).expect("Failed to initialize log4rs");

    info!("Server starting");

    // 创建共享状态
    let app_state = Arc::new(Mutex::new(Vec::<Question>::new()));

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::new(app_state.clone())) // 注册共享状态
            .service(add_question) // 注册API服务
            .service(get_questions) // 注册API服务
            .service(
                Files::new("/", "/home/wayne/source/practice/forum/frontend")
                    .index_file("index.html")
                    .prefer_utf8(true),
            )
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
