# aliyun-oss-action

用于 docusaurus 博客上传静态资源到阿里云 OSS。

sdk 下载:

https://github.com/aliyun/aliyun-oss-go-sdk

https://pkg.go.dev/github.com/aliyun/aliyun-oss-go-sdk/oss

sdk 文档:

https://help.aliyun.com/document_detail/88601.html

Docker Hub:

https://hub.docker.com/repository/docker/garfield998/aliyun-oss-action/general

## Get started

登录 Docker：

```bash
$ docker login -u <用户名> -p <密码>
```

构建镜像（注意命令最后有一个 `.`）：

```bash
$ docker build -t garfield998/aliyun-oss-action:1.0 .
```

查看构建出来的镜像：

```bash
$ docker image ls
```

发布镜像：

```bash
$ docker push garfield998/aliyun-oss-action:1.0
```
