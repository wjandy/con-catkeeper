var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/logs/42701ee5-4329-498c-bb2d-2af82b48e858",
  method: "GET",
  timeout: 10000,
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
