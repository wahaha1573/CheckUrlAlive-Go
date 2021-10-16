# CheckUrlAlive-Go.exe

## 简介

CheckUrlAlive-Go.exe是为了解决在信息收集中如通过子域名收集或者其他方式得到大量的域名而难以快速验证域名存活的窘境而开发的。测试2500个域名验证存活只需4秒。

输出格式

- url

- 网站title
- 中间件信息
- status_code

## 用法

### url.txt文件内支持格式

- www.baidu.com
- http://www.baidu.com
- https://www.baidu.com

CheckUrlAlive-Go.exe 

![0](E:\goProject\src\github.com\Mdxjj\CheckUrlAlive-Go\images\0.png) 

CheckUrlAlive-Go.exe -f url.txt

![1](E:\goProject\src\github.com\Mdxjj\CheckUrlAlive-Go\images\1.png)

