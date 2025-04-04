# Example Project Go lang model MVC

Project tutorial, learning for golang model MVC

## Recheck have start services
- Go-Version: 1.24
- Database: MYSQL
- Cache: Redis

## Install Project
- Init project
```
go mod download
```

- Create or copy from .env.example (*Command for bash or zsh)
```
cp .env.example .env
```

- Run project
```
go run main.go
```

## (Optional) Learning step tutorial
- Create folder
```
mkdir <projectname>
```

- move to folder
```
cd <projectname>
```

- init project
```
go mod init
```

- Download package
```
go get github.com/gofiber/fiber/v2
```
```
go get github.com/go-playground/validator/v10
```
```
go get github.com/dgrijalva/jwt-go
```
```
go get github.com/joho/godotenv
```
```
go get github.com/redis/go-redis/v9
```
```
go get -u gorm.io/gorm
```
```
go get gorm.io/driver/mysql
```
```
go get golang.org/x/crypto/bcrypt
```
```
go get github.com/spf13/viper
```

## (Optional) Docker command

- MYSQL
```
docker pull mysql
```
```
docker run --name <currentname> -e MYSQL_ROOT_PASSWORD=<currentpassword> -p 3306:3306 -d mysql
```
OR
```
docker run --name <currentname> -e MYSQL_ALLOW_EMPTY_PASSWORD=true -p 3306:3306 -d mysql
```

- Redis
```
docker pull redis
```
```
docker run --name <currentname> -p 6379:6379 -d redis
```
