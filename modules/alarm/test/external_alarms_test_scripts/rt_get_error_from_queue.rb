require "redis"
require "json"
redis = Redis.new

p redis.lpop("error_event:all")
