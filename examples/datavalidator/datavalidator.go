package main

import (
	"fmt"

	"github.com/rk-the-dev/golib-core/pkg/datavalidator"
)

// Define a nested struct for address validation
type Address struct {
	City string
	Pin  int
}

// Define the main struct for validation
type User struct {
	Name       string
	Email      string
	Phone      string
	Age        int
	PAN        string
	Aadhaar    string
	Passport   string
	CreditCard string
	Address    Address
}

func main() {
	// Create a sample user with invalid data
	user := User{
		Name:       "",                // Invalid (empty)
		Email:      "invalid-email",   // Invalid (not a valid email)
		Phone:      "12345",           // Invalid (not a valid phone number)
		Age:        16,                // Invalid (should be at least 18)
		PAN:        "ABCDE12345",      // Invalid (should end with a letter)
		Aadhaar:    "123456789012",    // Invalid (should start with 2-9)
		Passport:   "A123456",         // Invalid (should be 8 chars)
		CreditCard: "123456781234567", // Invalid (should be 16 digits)
		Address: Address{
			City: "",  // Invalid (empty)
			Pin:  999, // Invalid (should be at least 1000)
		},
	}

	// Define validation rules
	rules := []datavalidator.ValidationRule{
		{Field: "Name", Validator: datavalidator.NotEmpty, ErrorMsg: "Name is required"},
		{Field: "Email", Validator: datavalidator.IsEmail, ErrorMsg: "Invalid email address"},
		{Field: "Phone", Validator: datavalidator.IsPhoneNumber, ErrorMsg: "Invalid phone number"},
		{Field: "Age", Validator: datavalidator.MinValue(18), ErrorMsg: "Age must be at least 18"},
		{Field: "PAN", Validator: datavalidator.IsPAN, ErrorMsg: "Invalid PAN format"},
		{Field: "Aadhaar", Validator: datavalidator.IsAadhaar, ErrorMsg: "Invalid Aadhaar number"},
		{Field: "Passport", Validator: datavalidator.IsPassport, ErrorMsg: "Invalid Passport format"},
		{Field: "CreditCard", Validator: datavalidator.IsCreditCard, ErrorMsg: "Invalid Credit Card number"},
		{Field: "Address", Validator: datavalidator.ValidateNestedStruct, ErrorMsg: "Address validation failed"},
	}

	// Validate struct
	validationErrors := datavalidator.ValidateStruct(user, rules)

	// Print validation errors
	if len(validationErrors) > 0 {
		fmt.Println("❌ Validation Errors:")
		for field, errMsg := range validationErrors {
			fmt.Printf("  - %s: %s\n", field, errMsg)
		}
	} else {
		fmt.Println("✅ All validations passed!")
	}
}
