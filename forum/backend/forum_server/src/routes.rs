use crate::handlers::rest_api;
use actix_web::web;
pub fn config_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(rest_api::add_article)
        .service(rest_api::get_articles)
        .service(rest_api::vote_article);
}
