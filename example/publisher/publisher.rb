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

def message()
  host = 'http://localhost:5546'
  count = 1
  base = 268435456
  index = Time.now.to_i
  RestClient.post(
    "#{host}/message",
    {
      "userID" => 1,
      "queueName" => "greeting",
      "content" => index.to_s,
      "index" => index * base
    }.to_json,
    {content_type: :json, accept: :json}
  )

  count += 1
end


message()
