use actix_web::{App, HttpServer, middleware::Logger, web};
use log::info;
use log4rs;
use std::sync::Arc;

mod constants;
mod db;
mod handlers;
mod models;
mod routes;

use db::redis::create_redis_client;
use routes::config_routes;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    log4rs::init_file("log4rs.yml", Default::default()).expect("Failed to initialize log4rs");
    info!("Server starting");

    // 创建Redis客户端
    let client = create_redis_client("redis://127.0.0.1:6379/").map_err(|e| {
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
            .configure(config_routes)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
