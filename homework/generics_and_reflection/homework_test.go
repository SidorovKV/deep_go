package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](object T) string {
	builder := &strings.Builder{}
	obj := reflect.TypeOf(object)

	for i := 0; i < obj.NumField(); i++ {
		fieldName := obj.Field(i).Tag.Get("properties")
		parts := strings.Split(fieldName, ",")
		value := reflect.ValueOf(object).Field(i)

		if fieldName == "" || value.IsZero() && slices.Contains(parts, "omitempty") {
			continue
		}

		_, _ = builder.WriteString(fmt.Sprintf("%s=%v\n", parts[0], value))
	}

	return strings.TrimSpace(builder.String())
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestSerializationAnother(t *testing.T) {
	type Building struct {
		Address   string `properties:"address"`
		Floors    int    `properties:"floors"`
		Parking   bool   `properties:"parking"`
		IsPrivate bool   `properties:"is_private,omitempty"`
	}
	tests := map[string]struct {
		person Building
		result string
	}{
		"test case with empty fields": {
			result: "address=\nfloors=0\nparking=false",
		},
		"test case with fields": {
			person: Building{
				Address: "Moscow, Kremlin, 1",
				Floors:  3,
				Parking: true,
			},
			result: "address=Moscow, Kremlin, 1\nfloors=3\nparking=true",
		},
		"test case with omitempty field": {
			person: Building{
				Address:   "Moscow, Kremlin, 1",
				Floors:    3,
				Parking:   true,
				IsPrivate: true,
			},
			result: "address=Moscow, Kremlin, 1\nfloors=3\nparking=true\nis_private=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
