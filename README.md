### Metrics

Не понимаю, в чем проблема



когда был go 1.21.0, падала ошибка, что верия старая
обновил до ```go version go1.23.0 linux/amd64```
падает эта ошибка

```
Run go install golang.org/x/tools/cmd/goimports@latest
go: downloading golang.org/x/tools v0.31.0
go: golang.org/x/tools/cmd/goimports@latest: golang.org/x/tools@v0.31.0 requires go >= 1.23.0 (running go 1.22.12; GOTOOLCHAIN=local)
Error: Process completed with exit code 1.
```
