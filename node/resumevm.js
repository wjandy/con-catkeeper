var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/servers/dac3cc08-3ff0-4eee-83fc-bc28f4c8eb39/resume",
  method: "POST",
  timeout: 10000,
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
