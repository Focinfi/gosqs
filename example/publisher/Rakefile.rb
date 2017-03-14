require 'rest-client'
require 'json'

def register(user_id = 1, client_id = 1)
  host = 'http://localhost:5546'
  RestClient.post(
    "#{host}/register",
    {
      "userID": user_id,
      "queueName": "greeting",
      "clientID": client_id,
      "addresses": ["http://localhost:55446/greeting/#{client_id}"]
    }.to_json,
    {content_type: :json, accept: :json}
  )
end

def message(index)
  host = 'http://localhost:5546'
  base = 268435456
  res = RestClient.post(
    "#{host}/message",
    {
      "userID" => 1,
      "queueName" => "greeting",
      "content" => "#{index}",
      "index" => index
    }.to_json,
    {content_type: :json, accept: :json}
  )
  puts index, JSON.parse(res.body)
end

def apply_message_id_range(size = 10)
  host = 'http://localhost:5546'
  res = RestClient.put(
    "#{host}/messageID",
    {
      "userID" => 1,
      "queueName" => "greeting",
      "size" => size
    }.to_json,
    {content_type: :json, accept: :json}
  )

  
  data = JSON.parse(res.body)["data"]
  first = data["messageIDBegin"]
  last = data["messageIDEnd"]
  (first..last)
end

task :m do
  apply_message_id_range.each do |id|
    message(id)
  end
end

task :r do
  register()
end

task :a do
  apply_message_id_range()
end