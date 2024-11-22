package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bytes"

	"github.com/beevik/etree"
	"github.com/go-playground/validator/v10"
)

// User struct represents the user data
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func validateJSON(data []byte, user *User) error {
	validate := validator.New()
	err := json.Unmarshal(data, user)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	err = validate.Struct(user)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func validateXML(data []byte, user *User) error {
	// Define the XML schema in a string (in a real scenario, you'd load this from a file)
	xmlSchema := `
    <?xml version="1.0" encoding="UTF-8" ?>
    <xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
        <xs:element name="user">
            <xs:complexType>
                <xs:sequence>
                    <xs:element name="name" type="xs:string"/>
                    <xs:element name="age" type="xs:positiveInteger"/>
                    <xs:element name="email" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
        </xs:element>
    </xs:schema>
    `

	// Parse the XML schema
	schemaDoc, err := etree.ParseFromString(xmlSchema)
	if err != nil {
		return fmt.Errorf("error parsing XML schema: %w", err)
	}

	// Parse the incoming XML data
	dataDoc, err := etree.ParseFromReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error parsing XML data: %w", err)
	}

	// Validate the data against the schema
	err = validateXMLDocument(schemaDoc, dataDoc)
	if err != nil {
		return fmt.Errorf("XML validation failed: %w", err)
	}

	return nil
}

func validateXMLDocument(schema *etree.Document, data *etree.Document) error {
	// Helper function to perform element validation
	validateElement := func(schemaElem *etree.Element, dataElem *etree.Element) error {
		if schemaElem.Tag != dataElem.Tag {
			return fmt.Errorf("element tags mismatch: expected %s, got %s", schemaElem.Tag, dataElem.Tag)
		}

		for _, schemaChild := range schemaElem.ChildElements() {
			found := false
			for _, dataChild := range dataElem.ChildElements() {
				if err := validateElement(schemaChild, dataChild); err != nil {
					return err
				}
				found = true
				break
			}
			if !found {
				return fmt.Errorf("missing element %s", schemaChild.Tag)
			}
		}

		return nil
	}

	// Find the root elements in schema and data
	schemaRoot := schema.Root()
	dataRoot := data.Root()

	if schemaRoot == nil || dataRoot == nil {
		return fmt.Errorf("invalid schema or data document")
	}

	return validateElement(schemaRoot, dataRoot)
}

func main() {
	jsonData := []byte(`{"name": "John", "age": 30, "email": "john@example.com"}`)
	xmlData := []byte(`
    <?xml version="1.0" encoding="UTF-8" ?>
    <user>
        <name>Jane</name>
        <age>25</age>
        <email>jane@example.com</email>
    </user>
    `)

	var user User

	if err := validateJSON(jsonData, &user); err != nil {
		fmt.Println("JSON validation failed:", err)
		return
	}
	fmt.Println("JSON validation successful")

	if err := validateXML(xmlData, &user); err != nil {
		fmt.Println("XML validation failed:", err)
		return
	}
	fmt.Println("XML validation successful")
}
