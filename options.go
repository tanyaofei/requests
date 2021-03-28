package requests

// Http 链接参数, e.g: www.example.com?param1=value1&param2=value2
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Params map[string]string

// Http Body
// 如果使用了 Json 则 Data 和 Files 不再生效
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Json map[string]string

// Http Body
// 如果使用了 Json, Data 不再生效
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Data map[string]string

// Http Body
// 如果使用了 Json, Files 不再生效
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Files map[string]*[]byte

// Http Headers
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Headers map[string]string

// Http Cookies
// 可以在任何情况下使用:
//	1. 直接使用 requests.Request() 等方法中使用;
//	2. 在使用 requests.NewSession() 创建 session 时使用;
//	3. 使用已经创建的 session 中发送请求时使用
type Cookies map[string]string // also support []*http.Cookie

// 是否自动跳转
// 指定 3xxx 是否自动跳转
// 可以作用与:
// 	1. 直接使用 requests.Request(), requests.Get(), requests.Post().... 等方法
//  2. 在使用 requests.NewSession() 方法创建 request.session 对象时指定
// 注意, 在已经创建的 session 结构体中发送请求不再生效, 如 session.Get() 方法, 这种操作不被支持
type Redirect bool

// 是否校验 SSL 证书
// 指定 3xxx 是否自动跳转
// 可以作用与:
// 	1. 直接使用 requests.Request(), requests.Get(), requests.Post().... 等方法
//  2. 在使用 requests.NewSession() 方法创建 request.session 对象时指定
// 注意, 在已经创建的 session 结构体中发送请求不再生效, 如 session.Get() 方法, 这种操作不被支持
type Verify bool

// 指定代理
// 指定 3xxx 是否自动跳转
// 可以作用与:
// 	1. 直接使用 requests.Request(), requests.Get(), requests.Post().... 等方法
//  2. 在使用 requests.NewSession() 方法创建 request.session 对象时指定
// 注意, 在已经创建的 session 结构体中发送请求不再生效, 如 session.Get() 方法, 这种操作不被支持
type Proxy string
