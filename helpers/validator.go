package helpers

import (
	"regexp"
	"strings"
)

// Fungsi ini memeriksa apakah string adalah URL yang valid dan berakhir dengan salah satu ekstensi yang diinginkan dari slice.
func IsValidURLWithDesiredExtension(url string, desiredExtensions []string) bool {
	// Pattern regex untuk URL yang sederhana (silakan sesuaikan sesuai kebutuhan).
	// Ini hanya mendukung URL HTTP/HTTPS yang paling umum.
	urlPattern := `^(https?://)?([a-zA-Z0-9.-]+(\.[a-zA-Z]{2,})+(/.*)?)?$`

	// Compile regex pattern
	regex, err := regexp.Compile(urlPattern)
	if err != nil {
		return false
	}

	// Pecahkan URL untuk mendapatkan nama file terakhir
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return false // URL tidak valid
	}

	filename := parts[len(parts)-1]

	// Periksa apakah URL adalah URL yang valid
	if !regex.MatchString(url) {
		return false // URL tidak valid
	}

	// Periksa apakah ekstensi file cocok dengan salah satu yang diinginkan
	for _, ext := range desiredExtensions {
		if HasDesiredExtension(filename, ext) {
			return true // URL valid dengan ekstensi yang diinginkan
		}
	}

	return false // URL tidak valid atau tidak memiliki ekstensi yang diinginkan
}

// Fungsi ini memeriksa apakah string berakhir dengan salah satu ekstensi yang diinginkan dari slice.
func HasDesiredExtension(filename string, desiredExtension string) bool {
	// Pecahkan string ke dalam potongan berdasarkan tanda titik (.)
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return false // Tidak ada ekstensi dalam nama file
	}

	// Ambil ekstensi dari potongan terakhir
	extension := parts[len(parts)-1]

	// Periksa apakah ekstensi ada dalam daftar ekstensi yang diinginkan
	return strings.EqualFold(extension, desiredExtension)
}