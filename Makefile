docs:
	@~/go/bin/swag init --parseDependency -g app.go 

build:
	@GOOS=linux go build -o api-offers
	@upx api-offers


 

