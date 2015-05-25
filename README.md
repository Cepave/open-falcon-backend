falcon-fe
===

鉴于很多用户反馈UIC太难安装了（虽然我觉得其实很容易……），用Go语言重新实现了一个，也就是这个falcon-fe了。
另外，监控系统组件比较多，有不少web组件，比如uic、portal、alarm、dashboard，没有一个统一的地方汇总查看，
falcon-fe也做了一些快捷配置，类似监控系统的hao123导航了

问题：
目前用Go实现的这个UIC，没有集成LDAP登录功能，我对这块不熟，社区哪位大侠可以贡献个patch就完美了，没有patch，
有个使用LDAP登录校验的例子也行啊
