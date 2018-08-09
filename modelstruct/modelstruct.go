package modelstruct

import (
	"reflect"
	"strings"
)

type StructField struct {
	reflect.StructField
	ModelIndex int
}

func subFieldOf(parent, child reflect.StructField) *StructField {
	child.Index = append(parent.Index, child.Index...)
	return &StructField{child, -1}
}

type ModelStruct struct {
	Value        interface{}
	Type         reflect.Type
	Fields       []*StructField
	FieldsByName map[string]*StructField
	NumFields    int
}

func Get(value interface{}) (model *ModelStruct) {
	var reflectType reflect.Type
	var rt bool
	switch vt := value.(type) {
	case reflect.Type:
		reflectType = vt
		rt = true
	default:
		reflectType = reflect.ValueOf(value).Type()
	}

	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	m, ok := modelsStructMap.Load(reflectType)
	if ok {
		return m.(*ModelStruct)
	}

	var (
		reflectField reflect.StructField
		subField     *StructField
	)

	if rt {
		value = reflect.New(reflectType).Interface()
	}

	model = &ModelStruct{Type: reflectType, Value: value}

	addField := func(field *StructField) {
		for i, f := range model.Fields {
			if f.Name == field.Name {
				model.Fields = append(model.Fields[:i], model.Fields[i+1:]...)
				break
			}
		}
		model.Fields = append(model.Fields, field)
	}

	for numFields, i := reflectType.NumField(), 0; i < numFields; i++ {
		reflectField = reflectType.Field(i)
		if reflectField.Anonymous {
			for _, subField = range Get(reflect.New(reflectField.Type).Interface()).Fields {
				subField = subFieldOf(reflectField, subField.StructField)
				addField(subField)
			}
		} else {
			addField(&StructField{reflectField, -1})
		}
	}

	model.FieldsByName = map[string]*StructField{}

	for i, f := range model.Fields {
		f.ModelIndex = i
		model.FieldsByName[strings.ToLower(f.Name)] = f
	}

	model.NumFields = len(model.Fields)

	modelsStructMap.Store(reflectType, model)
	return model
}
