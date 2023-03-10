package _

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ddosify/go-faker/faker"
	"gitlab.com/ptami_lib/util"
)

type FakeType uint

type FakerOption struct {
	Count  int
	Fields []FakerOptions
}

type FakerOptions struct {
	FieldName string
	Type      string
	Prefix    *string
	Suffix    *string
	IsUnique  *bool
	Child     *FakerOption
	Children  *FakerOption
}

func fakerSetValue(fakerType string, dataType reflect.Kind, option FakerOptions) (value interface{}, err error) {
	var c = faker.NewFaker()
	if dataType == reflect.String {

	}
	switch fakerType {
	case "fullname":
		value = c.RandomPersonFullName()
		return
	case "address":
		value = c.RandomAddressStreetAddress()
		return
	case "city":
		value = c.RandomAddressCity()
		return
	case "number":
		value = c.RandomInt()
		return
	case "uuid":
		value = *option.Prefix + util.GetUlid()
		return
	default:
		err = errors.New("Faker not found !")
		return
	}
	return
}

func ifChildrenDummy[T any](option FakerOption) (data []map[string]interface{}, err error) {
	var _modelLv1 T

	if reflect.ValueOf(_modelLv1).Kind() != reflect.Struct {
		err = errors.New("Type data must be struct")
		return
	}

	if option.Count < 1 {
		err = errors.New("count must be greater than one")
		return
	}

	var metaValLv1 = reflect.ValueOf(&_modelLv1)
	// Level 1
	for i := 0; i < option.Count; i++ {
		var mapDataLv1 = map[string]interface{}{}
		for _, field := range option.Fields {
			if metaValLv1.FieldByName(field.FieldName) == (reflect.Value{}) {
				err = errors.New(fmt.Sprintf("fied %s is not defined in %s", field.FieldName, reflect.TypeOf(metaValLv1).Name()))
				return
			} else {
				var value interface{}
				if field.Child != nil {
					value, err = ifChildDummy(*field.Child)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value

				} else if field.Children != nil {
					value, err = ifChildrenDummy(*field.Child)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value
				} else {
					value, err = fakerSetValue(field.Type, metaValLv1.FieldByName(field.FieldName).Kind(), field)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value
				}
			}
		}
		data = append(data, mapDataLv1)
	}

	return
}

func ifChildDummy[T any](option FakerOption) (data map[string]interface{}, err error) {
	var _modelLv1 T

	if reflect.ValueOf(_modelLv1).Kind() != reflect.Struct {
		err = errors.New("Type data must be struct")
		return
	}

	var metaValLv1 = reflect.ValueOf(&_modelLv1)
	// Level 1
	for _, field := range option.Fields {
		if metaValLv1.FieldByName(field.FieldName) == (reflect.Value{}) {
			err = errors.New(fmt.Sprintf("fied %s is not defined in %s", field.FieldName, reflect.TypeOf(metaValLv1).Name()))
			return
		} else {
			var value interface{}
			if field.Child != nil {
				value, err = ifChildDummy(*field.Child)
				if err != nil {
					return
				}
				data[field.FieldName] = value

			} else if field.Children != nil {
				value, err = ifChildrenDummy(*field.Child)
				if err != nil {
					return
				}
				data[field.FieldName] = value
			} else {
				value, err = fakerSetValue(field.Type, metaValLv1.FieldByName(field.FieldName).Kind(), field)
				if err != nil {
					return
				}
				data[field.FieldName] = value
			}
		}
	}

	return
}

func CreateDummyData[T any, V any](option FakerOption) (_models []T, err error) {
	var _modelLv1 T
	if reflect.ValueOf(_modelLv1).Kind() != reflect.Struct {
		err = errors.New("Type data must be struct")
		return
	}

	if option.Count < 1 {
		err = errors.New("count must be greater than one")
		return
	}

	var metaValLv1 = reflect.ValueOf(&_modelLv1)

	var mapsDataLv1 = []map[string]interface{}{}

	// Level 1
	for i := 0; i < option.Count; i++ {
		var mapDataLv1 = map[string]interface{}{}
		for _, field := range option.Fields {
			if metaValLv1.FieldByName(field.FieldName) == (reflect.Value{}) {
				err = errors.New(fmt.Sprintf("fied %s is not defined", field.FieldName))
				return
			} else {
				var value interface{}
				if field.Child != nil {
					value, err = ifChildDummy[V](*field.Child)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value

				} else if field.Children != nil {
					value, err = ifChildrenDummy[V](*field.Child)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value
				} else {
					value, err = fakerSetValue(field.Type, metaValLv1.FieldByName(field.FieldName).Kind(), field)
					if err != nil {
						return
					}
					mapDataLv1[field.FieldName] = value
				}

				mapsDataLv1 = append(mapsDataLv1, mapDataLv1)
			}
		}
	}

	util.MapToStruct(mapsDataLv1, &_models)

	return
}
