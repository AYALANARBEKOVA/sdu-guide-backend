package utils

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"

)

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RemoveNumFromArrayOfNum(nums []int64, num int64) []int64 {
	index := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] == num {
			index = i
		}
	}
	return append(nums[:index], nums[index+1:]...)
}

func UniqueStrings(input []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, value := range input {
		if _, exists := uniqueMap[value]; !exists {
			uniqueMap[value] = true
			result = append(result, value)
		}
	}

	return result
}

const (
	chars    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenChars = len(chars)
)

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(lenChars)]
	}
	return string(result)
}

func UniqueInt64s(input []int64) []int64 {
	uniqueMap := make(map[int64]bool)
	var result []int64

	for _, value := range input {
		if _, exists := uniqueMap[value]; !exists {
			uniqueMap[value] = true
			result = append(result, value)
		}
	}

	return result
}

func JaccardSimilarity(setA, setB []int64) float64 {
	setAMap := make(map[int64]struct{}, len(setA))
	for _, id := range setA {
		setAMap[id] = struct{}{}
	}

	intersectionCount := 0
	unionCount := len(setA)
	for _, id := range setB {
		if _, found := setAMap[id]; found {
			intersectionCount++
		} else {
			unionCount++
		}
	}

	// Calculate Jaccard similarity
	if unionCount == 0 {
		return 0
	}
	return float64(intersectionCount) / float64(unionCount)
}
func GetTopNNums(recommendations map[int64]float64, n int) []int64 {
	type pair struct {
		ID    int64
		Score float64
	}

	var pairs []pair
	for id, score := range recommendations {
		pairs = append(pairs, pair{ID: id, Score: score})
	}

	// Sort by score in descending order
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Score > pairs[j].Score
	})

	var topIDs []int64
	for i := 0; i < len(pairs) && i < n; i++ {
		topIDs = append(topIDs, pairs[i].ID)
	}

	return topIDs
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func LoadTranslations(lang string) (map[string]string, error) {
	filePath := fmt.Sprintf("./internal/locales/%s.yml", lang)
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %s: %v", filePath, err)
	}

	var translations map[string]string
	if err := yaml.Unmarshal(file, &translations); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %v", err)
	}

	return translations, nil
}
