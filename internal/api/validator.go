package api

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// 注册自定义验证器
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("username", validateUsername)
	
	// 注册标签名称函数
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值让required标签处理
	}
	
	// 中国手机号码正则
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 {
		return false
	}
	
	// 密码必须包含数字和字母
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	
	return hasNumber && hasLetter
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	
	// 用户名只能包含字母、数字和下划线
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

func ValidateStruct(c *gin.Context, obj interface{}) error {
	if err := validate.Struct(obj); err != nil {
		HandleError(c, err)
		return err
	}
	return nil
}

func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		HandleError(c, err)
		return err
	}
	
	return ValidateStruct(c, obj)
}

func BindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		HandleError(c, err)
		return err
	}
	
	return ValidateStruct(c, obj)
}

func BindURIAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		HandleError(c, err)
		return err
	}
	
	return ValidateStruct(c, obj)
}