## Go-kit golang基础库
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/comeonjy/go-kit)
![GitHub](https://img.shields.io/github/license/comeonjy/go-kit)
![GitHub issues](https://img.shields.io/github/issues/comeonjy/go-kit)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/comeonjy/go-kit)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/comeonjy/go-kit)
![GitHub pull requests](https://img.shields.io/github/issues-pr/comeonjy/go-kit)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/comeonjy/go-kit)
![GitHub last commit](https://img.shields.io/github/last-commit/comeonjy/go-kit)
![GitHub repo size](https://img.shields.io/github/repo-size/comeonjy/go-kit)
![GitHub language count](https://img.shields.io/github/languages/count/comeonjy/go-kit)
![Lines of code](https://img.shields.io/tokei/lines/github/comeonjy/go-kit)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/comeonjy/go-kit)
![GitHub contributors](https://img.shields.io/github/contributors-anon/comeonjy/go-kit)
![Sourcegraph for Repo Reference Count](https://img.shields.io/sourcegraph/rrc/github.com/comeonjy/go-kit)
![GitHub top language](https://img.shields.io/github/languages/top/comeonjy/go-kit)

### TODO

- [ ] kafka将消费者均匀分配在所有pod上
- [ ] kafka分区动态更新
- [ ] Elastic APM 链路追踪 https://www.modb.pro/db/211621
- [ ] grpc客户端拦截器、服务端拦截器、resolver
- [ ] app灵活注册server
- [ ] 增加CICD

### 约定
1. 环境变量存储保密信息、配置文件存储常规配置

### 使用
```shell
# 安装
go install github.com/comeonjy/go-kit/cmd/kit
# 创建新项目
kit new demo-project
```


### JetBrains OS licenses
Go-kit是根据JetBrains s.r.o 授予的免费JetBrains开源许可证与GoLand一起开发的，因此在此我要表示感谢。
<a href="https://www.jetbrains.com/?from=go-kit" target="_blank"><img src="https://tva1.sinaimg.cn/large/0081Kckwgy1gkl0xz7y4uj30zz0u042c.jpg" width="30%"  /></a>

### License
© JiangYang, 2020~time.Now

Released under the MIT [License](https://github.com/comeonjy/go-kit/blob/master/LICENSE)