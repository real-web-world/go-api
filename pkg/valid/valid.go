package valid

import (
	"reflect"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/real-web-world/go-api/pkg/gin"
)

type (
	BoolStr          = string
	DefaultValidator struct {
		once     sync.Once
		validate *validator.Validate
	}
)

const (
	BoolStrTrue    BoolStr = "true"
	BoolStrFalse   BoolStr = "false"
	defaultTagName         = "binding"
	LangZH                 = "zh"
	LangEN                 = "en"
	LangZHtw               = "zh_tw"
)

var (
	phoneReg = regexp.MustCompile(`^1[3456789]\d{9}$`)
)

var _ binding.StructValidator = &DefaultValidator{}

func validPhone(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return phoneReg.MatchString(val)
}
func validPhoneOrEmpty(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return val == "" || phoneReg.MatchString(val)
}
func validBoolStr(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return val == BoolStrTrue || val == BoolStrFalse
}

// 仅允许编辑非root用户 root由系统自带
func validUserLevel(fl validator.FieldLevel) bool {
	level := ginApp.UserLevel(fl.Field().String())
	switch level {
	case ginApp.LevelAdmin, ginApp.LevelAuthor, ginApp.LevelGeneral:
		return true
	default:
		return false
	}
}
func validGender(fl validator.FieldLevel) bool {
	level := ginApp.Gender(fl.Field().String())
	switch level {
	case ginApp.GenderMan, ginApp.GenderWoman, ginApp.GenderUnknown:
		return true
	default:
		return false
	}
}
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}
func (v *DefaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}
func (v *DefaultValidator) InitTrans(transMap map[string]*ut.Translator) {
	v.lazyInit()
	if transMap != nil {
		for locale := range transMap {
			switch locale {
			case LangZH:
				regTrans := v.validate.RegisterTranslation
				t, ok := transMap[locale]
				if !ok {
					panic("require zh Translator")
				}
				_ = regTrans("validBoolStr", *t, func(ut ut.Translator) error {
					return ut.Add("validBoolStr", "布尔值{0}格式不正确,只能为 true 或 false",
						true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validBoolStr", fe.Value().(string))
					return t
				})
				_ = regTrans("phone", *t, func(ut ut.Translator) error {
					return ut.Add("phone", "手机号{0}格式不正确",
						true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("phone", fe.Value().(string))
					return t
				})
				_ = regTrans("phoneOrEmpty", *t, func(ut ut.Translator) error {
					return ut.Add("phoneOrEmpty", "手机号{0}格式不正确",
						true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("phoneOrEmpty", fe.Value().(string))
					return t
				})
				_ = regTrans("validUserLevel", *t, func(ut ut.Translator) error {
					return ut.Add("validUserLevel", "用户等级 {0} 不是合法的值",
						true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validUserLevel", string(fe.Value().(ginApp.UserLevel)))
					return t
				})
				_ = regTrans("validGender", *t, func(ut ut.Translator) error {
					return ut.Add("validGender", "性别 {0} 不是合法的值",
						true)
				}, func(ut ut.Translator, fe validator.FieldError) string {
					t, _ := ut.T("validGender", string(fe.Value().(ginApp.Gender)))
					return t
				})
			case LangZHtw:
			case LangEN:
			default:
			}
		}
	}
}
func (v *DefaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName(defaultTagName)
		reg := v.validate.RegisterValidation
		_ = reg("validBoolStr", validBoolStr)
		_ = reg("phone", validPhone)
		_ = reg("phoneOrEmpty", validPhoneOrEmpty)
		_ = reg("validUserLevel", validUserLevel)
		_ = reg("validGender", validGender)
	})
}

func RegRule(rules map[string]validator.Func) {
	v := binding.Validator.Engine().(*validator.Validate)
	for rule, fn := range rules {
		_ = v.RegisterValidation(rule, fn)
	}
}
func RegTrans(tarns map[string][]interface{}, t *ut.Translator) {
	v := binding.Validator.Engine().(*validator.Validate)
	for tag, fn := range tarns {
		_ = v.RegisterTranslation(tag, *t, fn[0].(func(ut ut.Translator) error),
			fn[1].(func(ut ut.Translator, fe validator.FieldError) string))
	}
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
