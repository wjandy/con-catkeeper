var request = require("request")
//request.setHeader("Token", "d51d6ad0-ef4d-4714-9b5c-ce08533517b8")
request({
  headers: {
    "Token":"d51d6ad0-ef4d-4714-9b5c-ce08533517b8"
  },
  uri: "http://10.72.84.145/catkeeper/v1/servers",
  method: "POST",
  timeout: 10000,
  body: JSON.stringify({"Name":"wandytest17","imagename":"vm222","cpu":2,"mem":4,"disk":60,"hostipaddress":"10.72.84.145","description":"just for test"}),
  followRedirect: true,
  maxRedirects: 1,
}, function(error, response,body){
    console.log(body)
    console.log(error)
})
