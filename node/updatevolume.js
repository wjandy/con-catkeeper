var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/volumes/70087c8f-1be3-484f-823e-ca1ef3fa79c8/update",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"name":"newname","description":"new desc"}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
