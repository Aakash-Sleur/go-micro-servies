.PHONY: up down auth products orders kill-daemon

up:
	docker-compose up --build

down:
	docker-compose down

auth:
	cd auth && CompileDaemon -build="go build -o auth-service.exe ." -command="./auth-service.exe"

products:
	cd products && CompileDaemon -build="go build -o product-service.exe ." -command="./product-service.exe"

orders:
	cd orders && CompileDaemon -build="go build -o orders-service.exe ." -command="./orders-service.exe"

kill-daemon:
	powershell "Get-Process | Where-Object {$$_.ProcessName -like '*CompileDaemon*'} | Stop-Process -Force"

