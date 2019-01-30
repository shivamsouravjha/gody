package gody

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrInvalidBody struct {
	Kind reflect.Kind
}

func (e *ErrInvalidBody) Error() string {
	return fmt.Sprintln("invalid body:", e.Kind)
}

type ErrInvalidTag struct {
	Format string
}

func (e *ErrInvalidTag) Error() string {
	return fmt.Sprintln("invalid tag:", e.Format)
}

type (
	Field struct {
		Name  string
		Value string
		Tags  map[string]string
	}
)

func Serialize(b interface{}) ([]Field, error) {
	valueOf := reflect.ValueOf(b)
	if valueOf.Kind().String() != "struct" {
		return nil, &ErrInvalidBody{Kind: valueOf.Kind()}
	}

	typeOf := reflect.TypeOf(b)
	fields := make([]Field, 0)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		tagString := field.Tag.Get("validate")
		if tagString == "" {
			continue
		}

		tagFormats := strings.Fields(tagString)
		tags := make(map[string]string, 0)
		for i := 0; i < len(tagFormats); i++ {
			tagFormatSplitted := strings.Split(tagFormats[i], "=")
			if len(tagFormatSplitted) != 2 {
				return nil, &ErrInvalidTag{Format: tagFormats[i]}
			}
			tags[tagFormatSplitted[0]] = tagFormatSplitted[1]
		}

		fields = append(fields, Field{
			Name:  strings.ToLower(field.Name),
			Value: fmt.Sprintf("%s", valueOf.FieldByName(field.Name)),
			Tags:  tags})
	}

	return fields, nil
}
