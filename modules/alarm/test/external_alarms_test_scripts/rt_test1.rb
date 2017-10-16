require "redis"
require "json"
redis = Redis.new
data = {
  "alarm_type": "mtest",
  "status": 0,
  "target": "cepave.com",
  "metric": "http.tester",
  "current_step": 2,
  "event_time": Time.now.to_i,
  "priority": 2,
  "trigger_id": 1,
  "trigger_description": "resp == 1",
  "trigger_condition": "0 != 1",
  "note": "check external tester",
  "extended_blob": {
    "idc": "idc1",
    "platform": "pt1",
    "contact": "masato"
  }
}

redis.lpush("extnal_event:all", data.to_json)
