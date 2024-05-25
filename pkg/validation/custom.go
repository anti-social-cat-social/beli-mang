package validation

import (
	"belimang/pkg/helper"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidNameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Improved regex pattern:
	// - Allows spaces, hyphens, and apostrophes within a name part
	// - Uses a character class for allowed characters (letters, spaces, hyphens, apostrophes)
	pattern := `^[[:alpha:]]+(?: [[:alpha:]]+|[-\'][:alpha:]]+)*$`

	// COmpile regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Println("Error compiling regex")
	}

	return re.MatchString(value)
}

// nipPrefix is considered by spesific role. Pass the correct / valid NIP prefix
func ValidNIP(nipPrefix string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		raw := fl.Field().Int()
		value := strconv.Itoa(int(raw))

		if len(value) < 13 {
			return false
		}

		inputRole := value[:3]
		if inputRole != nipPrefix {
			return false
		}

		inputSex := string(value[3:4])
		if inputSex != "1" && inputSex != "2" {
			return false
		}

		inputYearStr := value[4:8]
		if valid, _ := helper.NumberInRange(inputYearStr, 2000, time.Now().Year()); !valid {
			return false
		}

		inputMonthStr := value[8:10]
		if valid, _ := helper.NumberInRange(inputMonthStr, 1, 12); !valid {
			return false
		}

		lastNumber := value[10:]
		if _, err := strconv.Atoi(lastNumber); err != nil {
			return false
		}

		return true
	}
}
