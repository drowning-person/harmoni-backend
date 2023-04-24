package validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

// 全局翻译器
var (
	trans ut.Translator
	v     = validator.New()
)

// InitTrans locale 指定你想要的翻译 环境
func InitTrans(locale string) (err error) {
	//中文
	zhT := zh.New()
	//英文
	enT := en.New()
	// 第一个参数 是备用的语言
	uni := ut.New(enT, zhT, enT)

	//local 一般会在前端的请求头中 定义Accept-Language
	var ok bool
	trans, ok = uni.GetTranslator(locale)
	if !ok {
		return fmt.Errorf("uni.GetTranslator failed:%s ", locale)
	}

	switch locale {
	case "en":
		err = enTrans.RegisterDefaultTranslations(v, trans)
	case "zh":
		err = zhTrans.RegisterDefaultTranslations(v, trans)
	default:
		// 默认是英文
		err = enTrans.RegisterDefaultTranslations(v, trans)
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
	return err
}

func Validate(dataStruct interface{}) error {
	//验证器注册翻译器
	err := v.Struct(dataStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(trans))
		}
	}
	return nil
}
