module intelligent-guidance-system/service/auth

go 1.21

require (
    github.com/go-kratos/kratos/v2 v2.7.2
    github.com/google/wire v0.5.0
    github.com/google/uuid v1.6.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/mysql v1.5.2
    github.com/gin-gonic/gin v1.9.1
    github.com/go-playground/validator/v10 v10.16.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/redis/go-redis/v9 v9.3.0
)

replace intelligent-guidance-system => ../../