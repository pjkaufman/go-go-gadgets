package linter

import (
	"regexp"
	"strconv"
	"strings"
)

// regex and other validation information is based on https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch04s13.html
var (
	isbn10Regex               = regexp.MustCompile("^(ISBN(-10)?:? )?([0-9X]{10}$|[- 0-9X]{13}$)")
	isbn10DigitChecksRegex    = regexp.MustCompile("^[0-9]{1,5}[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9X]$")
	isbn13Regex               = regexp.MustCompile("^(ISBN(-13)?:? )?([0-9]{13}$|([- 0-9]{17}$))")
	isbn13SeparatorCheckRegex = regexp.MustCompile("([0-9]+[- ]){4}")
	isbn13DigitChecksRegex    = regexp.MustCompile("^97[89][- ]?[0-9]{1,5}[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9]$")
)

// IsValidISBN checks if the given string is a valid ISBN-10 or ISBN-13
func IsValidISBN(isbn string) bool {
	isbn10Submatches := isbn10Regex.FindStringSubmatch(isbn)
	if len(isbn10Submatches) != 0 {
		if isbn10DigitChecksRegex.MatchString(isbn10Submatches[3]) {
			var (
				sum     int
				cleaned = cleanupIsbn(isbn)
				chars   = strings.Split(cleaned, "")
				last    = chars[len(chars)-1]
				check   string
			)
			chars = chars[:len(chars)-1]

			if len(chars) == 9 {
				// checksum is done in reverse for isbn 10
				for l, r := 0, len(chars)-1; l < r; l, r = l+1, r-1 {
					chars[l], chars[r] = chars[r], chars[l]
				}

				for i := 0; i < len(chars); i++ {
					digit, _ := strconv.Atoi(chars[i])
					sum += (i + 2) * digit
				}
				checkDigit := 11 - (sum % 11)
				if checkDigit == 10 {
					check = "X"
				} else if checkDigit == 11 {
					check = "0"
				} else {
					check = strconv.Itoa(checkDigit)
				}

				if check == last {
					return true
				}
			}
		}
	}

	isbn13Submatches := isbn13Regex.FindStringSubmatch(isbn)
	if len(isbn13Submatches) != 0 {
		if isbn13DigitChecksRegex.MatchString(isbn13Submatches[3]) && (isbn13Submatches[4] == "" || isbn13SeparatorCheckRegex.MatchString(isbn13Submatches[4])) {
			var (
				sum     int
				cleaned = cleanupIsbn(isbn)
				chars   = strings.Split(cleaned, "")
				last    = chars[len(chars)-1]
				check   string
			)
			chars = chars[:len(chars)-1]

			for i := 0; i < len(chars); i++ {
				digit, _ := strconv.Atoi(chars[i])
				sum += (i%2*2 + 1) * digit
			}
			checkDigit := 10 - (sum % 10)
			if checkDigit == 10 {
				check = "0"
			} else {
				check = strconv.Itoa(checkDigit)
			}

			if check == last {
				return true
			}
		}
	}

	return false
}

func cleanupIsbn(isbn string) string {
	cleaned := strings.TrimPrefix(isbn, "ISBN-10: ")
	cleaned = strings.TrimPrefix(cleaned, "ISBN-13: ")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	return cleaned
}
