
build:
	go build sendgmail.go


osx:
	GOOS=darwin GOARCH=amd64 go build -o sendgmail.osx sendgmail.go

linux-arm:
	GOOS=linux GOARCH=arm go build -o sendgmail.linux-arm sendgmail.go

linux-x86:
	GOOS=linux GOARCH=386 go build -o sendgmail.linux-x86 sendgmail.go

linux-x64:
	GOOS=linux GOARCH=amd64 go build -o sendgmail.linux-x86 sendgmail.go

freebsd-x86:
	GOOS=freebsd GOARCH=386 go build -o sendgmail.freebsd-x86 sendgmail.go

freebsd-x64:
	GOOS=freebsd GOARCH=amd64 go build -o sendgmail.freebsd-x64 sendgmail.go

windows-x86:
	GOOS=windows GOARCH=386 go build -o sendgmail-x86.exe sendgmail.go

windows-x64:
	GOOS=windows GOARCH=amd64 go build -o sendgmail-x64.exe sendgmail.go
