falcon-fe
===

鉴于很多用户反馈UIC太难安装了（虽然我觉得其实很容易……），用Go语言重新实现了一个，也就是这个falcon-fe了。

另外，监控系统组件比较多，有不少web组件，比如uic、portal、alarm、dashboard，没有一个统一的地方汇总查看，
falcon-fe也做了一些快捷配置，类似监控系统的hao123导航了

**问题**
目前用Go实现的这个UIC，没有集成LDAP登录功能，我对这块不熟，社区哪位大侠可以贡献个patch就完美了，没有patch，
有个使用LDAP登录校验的例子也行啊

# 安装Go语言环境

```
cd ~
wget http://dinp.qiniudn.com/go1.4.1.linux-amd64.tar.gz
tar zxvf go1.4.1.linux-amd64.tar.gz
mkdir -p workspace/src
echo "" >> .bashrc
echo 'export GOROOT=$HOME/go' >> .bashrc
echo 'export GOPATH=$HOME/workspace' >> .bashrc
echo 'export PATH=$GOROOT/bin:$GOPATH/bin:$PATH' >> .bashrc
echo "" >> .bashrc
source .bashrc
```

# 编译安装fe模块

```
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/fe.git
cd fe
go get ./...
./control build
./control start
```

# 二进制安装

```
wget http://7xiumq.com1.z0.glb.clouddn.com/open-falcon-fe-0.0.1.tar.gz
tar zxvf open-falcon-fe-0.0.1.tar.gz
# modify configuration
./control start
```

# 配置介绍

```
{
    "log": "debug",
    "http": {
        "enabled": true,
        "listen": "0.0.0.0:1234" 自己随便搞个端口，别跟现有的重复了，可以使用8080，与老版本保持一致
    },
    "cache": {
        "enabled": true,
        "redis": "127.0.0.1:6379", 这个redis跟judge、alarm用的redis不同，这个只是作为缓存来用
        "idle": 10,
        "max": 1000,
        "timeout": {
            "conn": 10000,
            "read": 5000,
            "write": 5000
        }
    },
    "salt": "0i923fejfd3", 搞一个随机字符串
    "canRegister": true,
    "uic": {
        "addr": "root:password@tcp(127.0.0.1:3306)/fe?charset=utf8&loc=Asia%2FChongqing",
        "idle": 10,
        "max": 100
    },
    "shortcut": {
        "falconPortal": "http://11.11.11.11:5050/", 浏览器可访问的portal地址
        "falconDashboard": "http://11.11.11.11:7070/", 浏览器可访问的dashboard地址
        "falconAlarm": "http://11.11.11.11:6060/" 浏览器可访问的alarm的http地址
    }
}
```

# 设置root账号的密码

该项目中的注册用户是有不同角色的，目前分三种角色：普通用户、管理员、root账号。系统启动之后第一件事情应该是设置root的密码，浏览器访问：http://fe.example.com/root?password=abc （此处假设你的项目访问地址是fe.example.com，也可以使用ip）,这样就设置了root账号的密码为abc。普通用户可以支持注册。
