var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/servers/f0ba5172-ce08-4ead-a615-d4bf36225d1c/attach",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"voluuid":"2049b3be-53d5-48b8-a038-d4a5359fa0ac"}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
