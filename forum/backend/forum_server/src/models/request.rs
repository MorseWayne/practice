use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Article {
    pub title: String,
    pub tags: Vec<String>,
    pub content: String,
    pub user_id: u32,
    #[serde(default)]
    pub view_count: u32,
    #[serde(default)]
    pub vote_count: i32,
    #[serde(default)]
    pub answer_count: u32,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Vote {
    pub user_id: u32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ArticleQuery {
    #[serde(default = "default_sort")]
    sort: String,
    #[serde(default = "default_order")]
    order: String,
    #[serde(default = "default_limit")]
    limit: u32,
}

fn default_sort() -> String {
    "time".to_string()
}

fn default_order() -> String {
    "desc".to_string()
}

fn default_limit() -> u32 {
    10
}
