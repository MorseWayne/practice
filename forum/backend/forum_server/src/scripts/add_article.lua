-- 添加文章的Lua脚本
-- 参数: KEYS[1]=ARTICLE_KEY, KEYS[2]=CREATE_TIME_KEY
-- ARGV[1]=title, ARGV[2]=content
local cjson = require 'cjson'
local article_id = redis.call('INCR', KEYS[1])
local time = redis.call('TIME')
local article_key = KEYS[1] .. ':' .. article_id
redis.call('ZADD', KEYS[2], time[1], article_key)
local tags = cjson.decode(ARGV[4])
redis.call('HSET', article_key, 'title', ARGV[1], 'content', ARGV[2], 'user_id',
           ARGV[3], 'tags', cjson.encode(tags))
return article_id
