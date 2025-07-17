use actix_web::{App, HttpServer, middleware::Logger};
use actix_files::Files;
use log::info;
use log4rs;  // 导入log4rs

#[actix_web::main]
async fn main() -> std::io::Result<()>
{
    // 从配置文件初始化log4rs
    log4rs::init_file("log4rs.yml", Default::default())
        .expect("Failed to initialize log4rs");

    info!("Server starting");

    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())  // HTTP请求日志中间件
            .service(Files::new("/", "/home/wayne/source/practice/forum/frontend")
                .index_file("index.html")
                .prefer_utf8(true))
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
