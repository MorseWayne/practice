use super::handlers::articles::*;
use actix_web::web;

pub fn config_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(add_article)
        .service(get_articles)
        .service(vote_article);
}
