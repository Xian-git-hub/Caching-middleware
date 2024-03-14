# Caching-middleware
这是一个简单的缓存中间件，可以在客户端与服务端之间添加一个redis数据库，将热点数据缓存下来。

目前实现的功能有热点数据缓存，带缓存的日志。

日志会在每天零点新创建，并且日志的名字为创建的日期。值得注意的是，在使用过程中日志文件会略大。

关于服务器，日志文件路径等配置，可以在setting文件夹中的json文件中进行配置。配置之后需重启程序才能生效。

有任何使用问题请联系我，邮箱:2213630742@qq.com。

项目中有遇到的问题和一些思考我记录在`开发日志`中，可以在项目中看到，希望会有帮助。

