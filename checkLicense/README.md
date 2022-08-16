# UOS激活码批量查询

## 安装依赖

```shell
go mod tidy
```



## 编译

```shell
go build checkLicense.go
```



## 其他

* 目前仅支持查询 `.License.txt`， 这个是很早之前的格式，如果需要读取其他文件，续修改`ReadKey()`

  