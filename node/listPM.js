var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://10.72.84.145/catkeeper/v1/pms",
  method: "GET",
  timeout: 10000,
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
