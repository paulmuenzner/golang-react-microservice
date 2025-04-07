package config

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
)

// RegexPatterns struct holds compiled regex expressions
type RegexPatterns struct {
	Ipv4     *regexp.Regexp
	Ipv6     *regexp.Regexp
	UUID     *regexp.Regexp
	URL      *regexp.Regexp
	Email    *regexp.Regexp
	Password *regexp.Regexp
}

// Global regex variable
var Regex RegexPatterns

// Raw regex patterns from JSON
type rawRegexPatterns struct {
	Ipv4     string `json:"Ipv4"`
	Ipv6     string `json:"Ipv6"`
	UUID     string `json:"UUID"`
	URL      string `json:"URL"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// LoadRegexConfig loads regex patterns from JSON file
func LoadRegexConfig(filepath string) {
	// Read file
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading regex config file: %v", err)
	}

	// Parse JSON
	var raw rawRegexPatterns
	err = json.Unmarshal(file, &raw)
	if err != nil {
		log.Fatalf("Error parsing regex config JSON: %v", err)
	}

	// Compile regex patterns
	Regex.Ipv4 = compileRegex(raw.Ipv4, "Ipv4")
	Regex.Ipv6 = compileRegex(raw.Ipv6, "Ipv6")
	Regex.UUID = compileRegex(raw.UUID, "UUID")
	Regex.URL = compileRegex(raw.URL, "URL")
	Regex.Email = compileRegex(raw.Email, "Email")
	Regex.Password = compileRegex(raw.Password, "Password")

	log.Println("✅ Regex patterns loaded successfully from", filepath)
}

// compileRegex is a helper function to compile regex patterns
func compileRegex(pattern string, name string) *regexp.Regexp {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Error compiling %s regex: %v", name, err)
	}
	return compiled
}

// USAGE

// // Load regex config from shared data
// 	config.LoadRegexConfig("../shared/data/regexConfig.json")

// 	// Example usage
// 	if config.Regex.Email.MatchString("test@example.com") {
// 		println("✅ Valid Email!")
// 	} else {
// 		println("❌ Invalid Email!")
// 	}
