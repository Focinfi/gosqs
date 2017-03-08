require 'rest-client'
require 'json'

host = 'http://localhost:5546'
count = 1
base = 268435456

# def resgister()
#   RestClient.post(
#     "#{host}/resgister",
#     {
#       "userID": 1,
#       "queueName": "greeting",
#       "clientID": 1,
#       "addresses": [":55446/greeting/1"]
#     }.to_json,
#     {content_type: :json, accept: :json}
#   )
# end

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


10.times {|i| message(i + 1) }
