# Request
A pkg like ***python requests*** for Golang.



# How to install

`go get github.com/tanyaofei/requests`



# Examples

### Fast send a request

```go
resp1, err := requests.Request("GET", "http://www.example.com")
resp, err := requests.Post("http://www.example.com", 
                           // also support for http.header
                           requests.Headers{"Referer":"http://www.github.com"},
                           requests.Data{"username":"tanyaofei"}, // also support for requests.Json
                           // also support for []*http.Cookies
                           requests.Cookies{"authentication": "AuthenticationValue"},
                           requests.Params("query":"queryValue"),                     
                           reuqests.Verify(false), 		// do not verify https, default: true
                           requests.Redirect(false),	// do not auto redirect for 301, 302..., default: true
                           requests.Proxy("http://proxy.com"))
```



## Send Request With Session

```go
session, err := requests.NewSession(
    requests.Header{"Authentication", "ABC"},
    requests.Verify(false),
    requests.Redirect(false),
    requests.Params{"p1":"v1"},
    requests.Cookies{"cookie":"cookieValue"},
    requests.Proxy("http://proxy.com"),
    requests.Redirect(false),
)
```

session will store cookies and redirect hostory, use `requests.Response.History` to get request redirect history



# Supported Options

`requests.Data`: form data  `Content-Type : application/x-www-form-urlencoded`

`requests.Json`: json data `Content-Type: application/json`

`requests.Params`: 

`requets.Verify` : 

`requets.Redirect`: 

`requets.Headers`

`http.Header`

`requests.Cookies`

`[]*https.Cookie`

