require "redis"
require "json"
redis = Redis.new
data = {
  "alarm_type": "mtest2",
  "status": 99,
  "target": "endpoint11",
  "metric": "cpu.idle",
  "current_step": 1,
  "event_time": Time.now.to_i,
  "priority": 9,
  "trigger_id": 1,
  "trigger_description": "resp == 1",
  "trigger_condition": "0 != 1",
  "note": "check external tester",
  "extended_blob": {
    "idc": "idc1",
    "platform": "owl-test",
    "contact": "郝琦"
  }
}

redis.lpush("extnal_event:all", data.to_json)
