package main

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

const (
	USER_SCHEMA = `{
        "type": "object",
        "properties": {
            "Name": {
                "type": "string",
				"minLength": 0,
				"maxLength": 10
            },
            "Age": {
                "type": "integer",
				"minimum": 0,
				"maximum": 130
            },
            "Articles": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "Title": {
                            "type": "string"
                        },
                        "Content": {
                            "type": "string"
                        }
                    },
                    "additionalProperties": true,
                    "required": ["Title", "Content"]
                }
            }
        },
        "additionalProperties": false,
        "required": ["Name", "Age", "Articles"]
    }`
)

type User struct {
	Name     string
	Age      int
	Articles []Article
}

type Article struct {
	Title   string
	Content string
}

func NewUserValidater() *UserValidater {
	return &UserValidater{sl: gojsonschema.NewStringLoader(USER_SCHEMA)}
}

type UserValidater struct {
	sl gojsonschema.JSONLoader
}

func (v *UserValidater) Validate(u *User) (error, []string) {
	dl := gojsonschema.NewGoLoader(u)
	return validate(v.sl, dl)
}

func (v *UserValidater) ValidateJSON(s string) (error, []string) {
	dl := gojsonschema.NewStringLoader(s)
	return validate(v.sl, dl)
}

func validate(sl gojsonschema.JSONLoader, dl gojsonschema.JSONLoader) (error, []string) {
	result, err := gojsonschema.Validate(sl, dl)
	if err != nil {
		return err, nil
	}
	if result.Valid() {
		return nil, nil
	} else {
		ss := make([]string, 0, len(result.Errors()))
		for _, desc := range result.Errors() {
			ss = append(ss, desc.Description())
		}

		return nil, ss
	}
}

func main() {
	u := &User{
		Name: "Ken",
		Age:  32,
		Articles: []Article{
			Article{
				Title:   "Title A",
				Content: "content A",
			},
			Article{
				Title:   "Title B",
				Content: "content B",
			},
		},
	}

	uv := NewUserValidater()
	validateAndShow(uv, u)
	validateAndShow(uv, &User{})

	u2 := &User{
		Name: "Ken",
		Age:  -1,
		Articles: []Article{
			Article{
				Title:   "Title A",
				Content: "content A",
			},
			Article{
				Title:   "Title B",
				Content: "content B",
			},
		},
	}
	validateAndShow(uv, u2)

	ju := `{"Name":"Ken", "Age": 32, "Articles":[{"Title": "Title A", "Content":"content A"}, {"Title":"Title B", "Content":"content B"}]}`
	validateJSONAndShow(uv, ju)
	ju2 := `{"Name":"Ken", "Age": -1, "Articles":[{"Title": "Title A", "Content":"content A", "Other": "other A"}, {"Title":"Title B", "Content":"content B"}]}`
	validateJSONAndShow(uv, ju2)
}

func validateAndShow(v *UserValidater, u *User) {
	err, errStrings := v.Validate(u)
	if err != nil {
		panic(err.Error())
	}
	if errStrings != nil {
		for _, desc := range errStrings {
			fmt.Println(desc)
		}
	} else {
		fmt.Println(u)
	}
}

func validateJSONAndShow(v *UserValidater, s string) {
	err, errStrings := v.ValidateJSON(s)
	if err != nil {
		panic(err.Error())
	}
	if errStrings != nil {
		for _, desc := range errStrings {
			fmt.Println(desc)
		}
	} else {
		fmt.Println(s)
	}
}
