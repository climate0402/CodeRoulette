package services

import (
	"encoding/json"
	"math/rand"
	"time"

	"coderoulette/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemService struct {
	db *gorm.DB
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ProblemData struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Language    string     `json:"language"`
	TestCases   []TestCase `json:"test_cases"`
	Solution    string     `json:"solution"`
}

func NewProblemService(db *gorm.DB) *ProblemService {
	return &ProblemService{db: db}
}

// GetRandomProblem returns a random problem based on difficulty and language
func (s *ProblemService) GetRandomProblem(difficulty, language string) (*ProblemData, error) {
	var problem database.Problem

	query := s.db.Where("difficulty = ? AND language = ?", difficulty, language)

	// Get count for random selection
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Get random offset
	rand.Seed(time.Now().UnixNano())
	offset := rand.Intn(int(count))

	// Get random problem
	if err := query.Offset(offset).First(&problem).Error; err != nil {
		return nil, err
	}

	// Parse test cases
	var testCases []TestCase
	if err := json.Unmarshal([]byte(problem.TestCases), &testCases); err != nil {
		return nil, err
	}

	return &ProblemData{
		ID:          problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
		Difficulty:  problem.Difficulty,
		Language:    problem.Language,
		TestCases:   testCases,
		Solution:    problem.Solution,
	}, nil
}

// GetProblemByID returns a problem by its ID
func (s *ProblemService) GetProblemByID(id uuid.UUID) (*ProblemData, error) {
	var problem database.Problem
	if err := s.db.First(&problem, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Parse test cases
	var testCases []TestCase
	if err := json.Unmarshal([]byte(problem.TestCases), &testCases); err != nil {
		return nil, err
	}

	return &ProblemData{
		ID:          problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
		Difficulty:  problem.Difficulty,
		Language:    problem.Language,
		TestCases:   testCases,
		Solution:    problem.Solution,
	}, nil
}

// CreateProblem creates a new problem
func (s *ProblemService) CreateProblem(data *ProblemData) error {
	// Serialize test cases
	testCasesJSON, err := json.Marshal(data.TestCases)
	if err != nil {
		return err
	}

	problem := &database.Problem{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Difficulty:  data.Difficulty,
		Language:    data.Language,
		TestCases:   string(testCasesJSON),
		Solution:    data.Solution,
	}

	return s.db.Create(problem).Error
}

// GetProblems returns a list of problems with pagination
func (s *ProblemService) GetProblems(page, limit int, difficulty, language string) ([]*ProblemData, int64, error) {
	var problems []database.Problem
	var count int64

	query := s.db.Model(&database.Problem{})

	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if language != "" {
		query = query.Where("language = ?", language)
	}

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get problems with pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	// Convert to ProblemData
	result := make([]*ProblemData, len(problems))
	for i, problem := range problems {
		var testCases []TestCase
		json.Unmarshal([]byte(problem.TestCases), &testCases)

		result[i] = &ProblemData{
			ID:          problem.ID,
			Title:       problem.Title,
			Description: problem.Description,
			Difficulty:  problem.Difficulty,
			Language:    problem.Language,
			TestCases:   testCases,
			Solution:    problem.Solution,
		}
	}

	return result, count, nil
}

// SeedProblems creates some sample problems for testing
func (s *ProblemService) SeedProblems() error {
	sampleProblems := []*ProblemData{
		{
			ID:          uuid.New(),
			Title:       "Two Sum",
			Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
			Difficulty:  "easy",
			Language:    "go",
			TestCases: []TestCase{
				{Input: "[2,7,11,15], 9", Output: "[0,1]"},
				{Input: "[3,2,4], 6", Output: "[1,2]"},
				{Input: "[3,3], 6", Output: "[0,1]"},
			},
			Solution: "func twoSum(nums []int, target int) []int {\n    m := make(map[int]int)\n    for i, num := range nums {\n        if j, ok := m[target-num]; ok {\n            return []int{j, i}\n        }\n        m[num] = i\n    }\n    return nil\n}",
		},
		{
			ID:          uuid.New(),
			Title:       "Reverse String",
			Description: "Write a function that reverses a string. The input string is given as an array of characters s.",
			Difficulty:  "easy",
			Language:    "go",
			TestCases: []TestCase{
				{Input: "[\"h\",\"e\",\"l\",\"l\",\"o\"]", Output: "[\"o\",\"l\",\"l\",\"e\",\"h\"]"},
				{Input: "[\"H\",\"a\",\"n\",\"n\",\"a\",\"h\"]", Output: "[\"h\",\"a\",\"n\",\"n\",\"a\",\"H\"]"},
			},
			Solution: "func reverseString(s []byte) {\n    left, right := 0, len(s)-1\n    for left < right {\n        s[left], s[right] = s[right], s[left]\n        left++\n        right--\n    }\n}",
		},
		{
			ID:          uuid.New(),
			Title:       "Valid Parentheses",
			Description: "Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.",
			Difficulty:  "medium",
			Language:    "go",
			TestCases: []TestCase{
				{Input: "\"()\"", Output: "true"},
				{Input: "\"()[]{}\"", Output: "true"},
				{Input: "\"(]\"", Output: "false"},
				{Input: "\"([)]\"", Output: "false"},
			},
			Solution: "func isValid(s string) bool {\n    stack := []rune{}\n    pairs := map[rune]rune{\n        ')': '(',\n        '}': '{',\n        ']': '[',\n    }\n    \n    for _, char := range s {\n        if char == '(' || char == '{' || char == '[' {\n            stack = append(stack, char)\n        } else if len(stack) > 0 && stack[len(stack)-1] == pairs[char] {\n            stack = stack[:len(stack)-1]\n        } else {\n            return false\n        }\n    }\n    \n    return len(stack) == 0\n}",
		},
	}

	for _, problem := range sampleProblems {
		if err := s.CreateProblem(problem); err != nil {
			// Continue if problem already exists
			continue
		}
	}

	return nil
}
