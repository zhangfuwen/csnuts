csnuts
======
A blog program based on golang and google app engine.

上传当前程序到github: 
==
        git push -u origin master

本地运行当前实例的办法:
==
        dev_server.py .
		在浏览器中打开localhost:8080

从GAE上下载存储着的实例的方法:
==
        appcfg.py download_app -A usrccapp -V 1-17 .
		其中usrccapp是当前csnuts.com的应用名.

将当前文件夹下的程序更新到GAE方法:
==
        app_cfg update .

安装sessions:
==
cd csnuts
git clone https://github.com/gorilla/sessions.git

Log May 14, 2014
==
当前代码是从gae上下载的最新实例1-17版本，也是当前正在使用的版本。这个版本利用github.com/gorilla/sessions实现了用户的注册和登陆的功能。目前尚不清楚这一功能好用到什么程序，可知的是，当前的版本注册后不会自动进入已登陆状态，而且登录和注册的页面实在是很难看。
下载代码时使用的命令是：
 cd demos
 appcfg.py download_app -A usrccapp -V 1-17 www.csnuts.com
 如果下载失败，rm -Rf www.csnuts.com，之后得新下载。

Log May 14, 2014
==
美化login页面的模板，从网上找一个login form的网页，包含一个html文件和两个css文件。现在登陆页面可用。
TODO:登陆后的重定向总是弹出一个确定对话框，重定向的页面不自动进入登陆状态。

Log May 14, 2014
==
完成注册页的美化，使用了与登陆页同样的模板。
TODO:重定向时页面不自动进入登陆状态。

Log May 15, 2014
==
完成注册页的美化，使用了与登陆页同样的模板。
修复重定向时页面不自动进入登陆状态。
TODO:美化注册及登陆出错返回页面

Log May 18, 2014
==
美化注册及登陆出错返回页面
TODO:把发文表单从首页移到单独的页面。

Log May 18,2014
==
把发文表单从首页移到post页面
