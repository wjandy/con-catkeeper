var request = require("request")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://127.0.0.1/catkeeper/v1/volumes",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"name":"testvolume","size":2,"description":"test rbd volume","volumetype":"rbd"}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
