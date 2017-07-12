require "redis"
require "json"
require 'net/http'

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
