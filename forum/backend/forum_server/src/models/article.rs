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
pub struct VoteRequest {
    pub user_id: u32,
}