var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/volumes/e0270635-8e2d-4ea3-9a0a-c8c966c150a9",
  method: "GET",
  timeout: 10000,
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
