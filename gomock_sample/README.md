gomock sample
===

gomock is mock library for golang interface.

## explain

### sample program

*Request* interface has Request method that is give url and return status code.

*StatusChecker* is check status code and return boolean true(http status 200) or false(not 200)
Library user can customise request with Request interface.

## generate mock

```
# make mock directory
mkdir -p mock_request

# generage mock
mockgen -source request.go -destination mock_request/mock_request.go
```

### mock test

I want to do test case that StatusChecker check `http://localhost/` and return status `500` then return `false`.

I create reuqest mock and that mock Request return 500 when call Request with `http://localhost/`.

```
go test -v
```

## refs:
- [Go Mockでインタフェースのモックを作ってテストする ](http://qiita.com/tenntenn/items/24fc34ec0c31f6474e6d)
