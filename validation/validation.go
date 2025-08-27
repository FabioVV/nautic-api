package validation

import "github.com/go-playground/validator"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func MsgForErrorTag(tag string) string {
	switch tag {
		case "required":
			return "This field is required"
		case "email":
			return "Invalid email"
	}
	return ""
}

func FmtErrReturn(err error)[]map[string]string {
	out := make([]map[string]string, len(err.(validator.ValidationErrors)))

	for i, err := range err.(validator.ValidationErrors) {
		out[i] = map[string]string{err.Field(): MsgForErrorTag(err.Tag())}
	}

	return out
}
