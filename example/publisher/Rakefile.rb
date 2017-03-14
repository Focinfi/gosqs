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

def message(count)
  host = 'http://localhost:5546'
  base = 268435456
  index = Time.now.to_i
  RestClient.post(
    "#{host}/message",
    {
      "userID" => 1,
      "queueName" => "greeting",
      "content" => "#{index}_#{count}",
      "index" => index * base + count
    }.to_json,
    {content_type: :json, accept: :json}
  )

  count += 1
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

  puts res.body
end

task :m do
  10.times {|i| message(i + 1) }
end

task :r do
  register()
end

task :a do
  apply_message_id_range(10)
end