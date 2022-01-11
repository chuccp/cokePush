SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o cokePush.exe github.com/chuccp/cokePush