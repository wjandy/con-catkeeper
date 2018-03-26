var request = require("request")
request({
  uri: "http://127.0.0.1/catkeeper/v1/volumes/6dabc9ca-49a1-4f87-bb76-a43c070a38c3/resize",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"size":5}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
