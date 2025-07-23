use redis::Client;
use std::sync::Arc;

/// 创建Redis客户端
pub fn create_redis_client(addr: &str) -> Result<Client, redis::RedisError> {
    let client = Client::open(addr)?;
    Ok(client)
}

/// 共享状态类型定义
pub type AppState = Arc<Client>;