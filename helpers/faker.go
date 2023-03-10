package helpers

import (
	"errors"

	"github.com/ddosify/go-faker/faker"
	"gitlab.com/ptami_lib/util"
)

type FakerOption struct {
	Count  int
	Fields []FakerOptions
}

type FakerOptions struct {
	FieldName string
	Type      string
	Prefix    *string
	Suffix    *string
	Child     *FakerOption
	Children  *FakerOption
}

func fakerSetValue(fakerType string, option FakerOptions) (value interface{}, err error) {
	var c = faker.NewFaker()
	// fullname, address, city, number, uuid, lorem, date, job, image-url, email, product, price
	switch fakerType {
	case "fullname":
		value = c.RandomPersonFullName()
	case "address":
		value = c.RandomAddressStreetAddress()
	case "city":
		value = c.RandomAddressCity()
	case "number":
		value = c.RandomInt()
	case "uuid":
		value = util.GetUlid()
	case "lorem":
		value = c.RandomLoremText()
	case "date":
		value = c.RandomDateRecent()
	case "job":
		value = c.RandomJobTitle()
	case "phone":
		value = c.RandomPhoneNumber()
	case "image-url":
		value = c.RandomAvatarImage()
	case "email":
		value = c.RandomEmail()
	case "product":
		value = c.RandomProduct()
	case "price":
		value = c.RandomPrice()
	default:
		err = errors.New("Faker not found !")
		return
	}

	if option.Prefix != nil {
		value = *option.Prefix + value.(string)
	}

	if option.Suffix != nil {
		value = value.(string) + *option.Suffix
	}

	return
}

func ifChildrenDummy(option FakerOption) (data []map[string]interface{}, err error) {

	if option.Count < 1 {
		err = errors.New("count must be greater than one")
		return
	}

	for i := 0; i < option.Count; i++ {
		var mapDataLv1 = map[string]interface{}{}
		for _, field := range option.Fields {
			var value interface{}
			if field.Child != nil {
				value, err = ifChildDummy(*field.Child)
				if err != nil {
					return
				}
			} else if field.Children != nil {
				value, err = ifChildrenDummy(*field.Children)
				if err != nil {
					return
				}
			} else {
				value, err = fakerSetValue(field.Type, field)
				if err != nil {
					return
				}
			}
			mapDataLv1[field.FieldName] = value
		}
		data = append(data, mapDataLv1)
	}

	return
}

func ifChildDummy(option FakerOption) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	for _, field := range option.Fields {
		var value interface{}
		if field.Child != nil {
			value, err = ifChildDummy(*field.Child)
			if err != nil {
				return
			}
		} else if field.Children != nil {
			value, err = ifChildrenDummy(*field.Children)
			if err != nil {
				return
			}
		} else {
			value, err = fakerSetValue(field.Type, field)
			if err != nil {
				return
			}
		}
		data[field.FieldName] = value

	}

	return
}

func CreateDummyData(option FakerOption, _entities interface{}) (err error) {

	if option.Count < 1 {
		err = errors.New("count must be greater than one")
		return
	}

	var mapsDataLv1 = []map[string]interface{}{}

	// Level 1
	for i := 0; i < option.Count; i++ {
		var mapDataLv1 = map[string]interface{}{}
		for _, field := range option.Fields {
			var value interface{}
			if field.Child != nil {
				value, err = ifChildDummy(*field.Child)
				if err != nil {
					return
				}
			} else if field.Children != nil {
				value, err = ifChildrenDummy(*field.Children)
				if err != nil {
					return
				}
			} else {
				value, err = fakerSetValue(field.Type, field)
				if err != nil {
					return
				}
			}
			mapDataLv1[field.FieldName] = value
		}
		mapsDataLv1 = append(mapsDataLv1, mapDataLv1)
	}

	util.MapToStruct(mapsDataLv1, _entities)

	return
}
