var request = require("request")
request({
  uri: "http://127.0.0.1/catkeeper/v1/volumes/wandy/3e087548-0efa-40ec-8d85-814a10f44afd/resize",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"size":3}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
