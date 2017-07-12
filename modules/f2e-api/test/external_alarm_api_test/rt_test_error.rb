require "redis"
require "json"
require 'net/http'

data = {
  "alarm_type": "mtest21",
  "status": 0,
  "target": "endpoint11",
  "metric": "cpu.idle",
  "current_step": 1,
  "event_time": Time.now.to_i,
  "priority": 3,
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

uri = URI('http://localhost:8088/api/v1/alarm/external_feeds')
# please refer cfg.json -> "services" of f2e-api 
tokens = {name: "servicename",sig: "servicesig"}
Net::HTTP.start(uri.host, uri.port) do |http|
  req = Net::HTTP::Post.new(uri)
  req['Content-Type'] = 'application/json'
  req['Apitoken'] = tokens.to_json

  req.body = data.to_json
  res = http.request(req)
  puts "response #{res.body}"
end
