/// 对于每篇文章，每200个赞可以展示在热榜上一天
/// VOTE_SCORE也是每篇文章每增加一个赞可以在热榜上存活的时间
pub const VOTE_SCORE: i32 = 24 * 60 * 60 / 200;

pub const ARTICLE_KEY: &str = "article";
pub const CREATE_TIME_KEY: &str = "create_time";
pub const VOTE_KEY: &str = "vote";
pub const SCORE_KEY: &str = "score";