// package main
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// User represents the user data with schema annotations
type User struct {
	Name  string `json:"name" xml:"name" validate:"required"`
	Age   int    `json:"age" xml:"age" validate:"gte=0,lte=150"`
	Email string `json:"email" xml:"email" validate:"required,email"`
}

func validateJSON(jsonData []byte, user *User) error {
	err := json.Unmarshal(jsonData, user)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		return fmt.Errorf("JSON schema validation failed: %w", err)
	}

	return nil
}
func validateXML(xmlData []byte, user *User) error {
	err := xml.Unmarshal(xmlData, user)
	if err != nil {
		return fmt.Errorf("error unmarshalling XML: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		return fmt.Errorf("XML schema validation failed: %w", err)
	}

	return nil
}
func main() {
	// JSON Data Example
	jsonData := []byte(`
	{
		"name": "John",
		"age": 30,
		"email": "john@example.com"
	}
	`)

	// XML Data Example
	xmlData := []byte(`
	<user>
		<name>Jane</name>
		<age>25</age>
		<email>jane@example.com</email>
	</user>
	`)

	var user User

	// Validate JSON
	err := validateJSON(jsonData, &user)
	if err != nil {
		fmt.Println("JSON Validation Error:", err)
		return
	}
	fmt.Println("JSON Data Validated:", user)

	// Validate XML
	err = validateXML(xmlData, &user)
	if err != nil {
		fmt.Println("XML Validation Error:", err)
		return
	}
	fmt.Println("XML Data Validated:", user)
}
