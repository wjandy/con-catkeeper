var request = require("request")
request({
  uri: "http://127.0.0.1/catkeeper/v1/servers/f0ba5172-ce08-4ead-a615-d4bf36225d1c/detach",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"voluuid":"3ecb15d5-db49-4c22-86d7-0cb0d790ac09","AttachPoint":"vdb"}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
