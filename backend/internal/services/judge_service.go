package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"coderoulette/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JudgeService struct {
	db *gorm.DB
}

type JudgeResult struct {
	Status    string           `json:"status"`    // passed, failed, error
	Score     int              `json:"score"`     // 0-100
	Runtime   int              `json:"runtime"`   // in milliseconds
	ErrorMsg  string           `json:"error_msg"` // error message if any
	TestCases []TestCaseResult `json:"test_cases"`
}

type TestCaseResult struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Passed   bool   `json:"passed"`
	Runtime  int    `json:"runtime"`
}

type SubmissionData struct {
	ID        uuid.UUID `json:"id"`
	MatchID   uuid.UUID `json:"match_id"`
	PlayerID  uuid.UUID `json:"player_id"`
	Code      string    `json:"code"`
	Language  string    `json:"language"`
	Status    string    `json:"status"`
	Score     int       `json:"score"`
	Runtime   int       `json:"runtime"`
	ErrorMsg  string    `json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`
}

func NewJudgeService() *JudgeService {
	return &JudgeService{}
}

func (s *JudgeService) SetDB(db *gorm.DB) {
	s.db = db
}

// SubmitCode submits code for judging
func (s *JudgeService) SubmitCode(ctx context.Context, matchID, playerID uuid.UUID, code, language string, testCases []TestCase) (*JudgeResult, error) {
	// Create submission record
	submission := &database.Submission{
		ID:       uuid.New(),
		MatchID:  matchID,
		PlayerID: playerID,
		Code:     code,
		Language: language,
		Status:   "running",
	}

	if err := s.db.Create(submission).Error; err != nil {
		return nil, err
	}

	// Judge the code
	result, err := s.judgeCode(code, language, testCases)
	if err != nil {
		submission.Status = "error"
		submission.ErrorMsg = err.Error()
		s.db.Save(submission)
		return nil, err
	}

	// Update submission with results
	submission.Status = result.Status
	submission.Score = result.Score
	submission.Runtime = result.Runtime
	submission.ErrorMsg = result.ErrorMsg
	s.db.Save(submission)

	return result, nil
}

// judgeCode executes the code and validates against test cases
func (s *JudgeService) judgeCode(code, language string, testCases []TestCase) (*JudgeResult, error) {
	switch language {
	case "go":
		return s.judgeGoCode(code, testCases)
	case "python":
		return s.judgePythonCode(code, testCases)
	case "javascript":
		return s.judgeJavaScriptCode(code, testCases)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// judgeGoCode judges Go code
func (s *JudgeService) judgeGoCode(code string, testCases []TestCase) (*JudgeResult, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "judge_go_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codeFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, err
	}

	// Compile Go code
	compileCmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "main"), codeFile)
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		return &JudgeResult{
			Status:   "error",
			Score:    0,
			ErrorMsg: fmt.Sprintf("Compilation error: %s", string(compileOutput)),
		}, nil
	}

	// Run test cases
	var testCaseResults []TestCaseResult
	passedCount := 0
	totalRuntime := 0

	for _, testCase := range testCases {
		start := time.Now()

		// Run the program with test input
		runCmd := exec.Command(filepath.Join(tempDir, "main"))
		runCmd.Stdin = strings.NewReader(testCase.Input)

		output, _ := runCmd.CombinedOutput()
		runtime := int(time.Since(start).Milliseconds())
		totalRuntime += runtime

		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(testCase.Output)

		passed := actual == expected
		if passed {
			passedCount++
		}

		testCaseResults = append(testCaseResults, TestCaseResult{
			Input:    testCase.Input,
			Expected: expected,
			Actual:   actual,
			Passed:   passed,
			Runtime:  runtime,
		})
	}

	// Calculate score
	score := (passedCount * 100) / len(testCases)
	status := "passed"
	if score < 100 {
		status = "failed"
	}

	return &JudgeResult{
		Status:    status,
		Score:     score,
		Runtime:   totalRuntime,
		TestCases: testCaseResults,
	}, nil
}

// judgePythonCode judges Python code
func (s *JudgeService) judgePythonCode(code string, testCases []TestCase) (*JudgeResult, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "judge_python_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codeFile := filepath.Join(tempDir, "main.py")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, err
	}

	// Run test cases
	var testCaseResults []TestCaseResult
	passedCount := 0
	totalRuntime := 0

	for _, testCase := range testCases {
		start := time.Now()

		// Run Python code with test input
		runCmd := exec.Command("python3", codeFile)
		runCmd.Stdin = strings.NewReader(testCase.Input)

		output, err := runCmd.CombinedOutput()
		if err != nil {
			return &JudgeResult{
				Status:   "error",
				Score:    0,
				ErrorMsg: fmt.Sprintf("Runtime error: %s", string(output)),
			}, nil
		}
		runtime := int(time.Since(start).Milliseconds())
		totalRuntime += runtime

		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(testCase.Output)

		passed := actual == expected
		if passed {
			passedCount++
		}

		testCaseResults = append(testCaseResults, TestCaseResult{
			Input:    testCase.Input,
			Expected: expected,
			Actual:   actual,
			Passed:   passed,
			Runtime:  runtime,
		})
	}

	// Calculate score
	score := (passedCount * 100) / len(testCases)
	status := "passed"
	if score < 100 {
		status = "failed"
	}

	return &JudgeResult{
		Status:    status,
		Score:     score,
		Runtime:   totalRuntime,
		TestCases: testCaseResults,
	}, nil
}

// judgeJavaScriptCode judges JavaScript code
func (s *JudgeService) judgeJavaScriptCode(code string, testCases []TestCase) (*JudgeResult, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "judge_js_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	codeFile := filepath.Join(tempDir, "main.js")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, err
	}

	// Run test cases
	var testCaseResults []TestCaseResult
	passedCount := 0
	totalRuntime := 0

	for _, testCase := range testCases {
		start := time.Now()

		// Run Node.js code with test input
		runCmd := exec.Command("node", codeFile)
		runCmd.Stdin = strings.NewReader(testCase.Input)

		output, err := runCmd.CombinedOutput()
		if err != nil {
			return &JudgeResult{
				Status:   "error",
				Score:    0,
				ErrorMsg: fmt.Sprintf("Runtime error: %s", string(output)),
			}, nil
		}
		runtime := int(time.Since(start).Milliseconds())
		totalRuntime += runtime

		actual := strings.TrimSpace(string(output))
		expected := strings.TrimSpace(testCase.Output)

		passed := actual == expected
		if passed {
			passedCount++
		}

		testCaseResults = append(testCaseResults, TestCaseResult{
			Input:    testCase.Input,
			Expected: expected,
			Actual:   actual,
			Passed:   passed,
			Runtime:  runtime,
		})
	}

	// Calculate score
	score := (passedCount * 100) / len(testCases)
	status := "passed"
	if score < 100 {
		status = "failed"
	}

	return &JudgeResult{
		Status:    status,
		Score:     score,
		Runtime:   totalRuntime,
		TestCases: testCaseResults,
	}, nil
}

// GetSubmission returns a submission by ID
func (s *JudgeService) GetSubmission(id uuid.UUID) (*SubmissionData, error) {
	var submission database.Submission
	if err := s.db.Preload("Player").Preload("Match").First(&submission, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &SubmissionData{
		ID:        submission.ID,
		MatchID:   submission.MatchID,
		PlayerID:  submission.PlayerID,
		Code:      submission.Code,
		Language:  submission.Language,
		Status:    submission.Status,
		Score:     submission.Score,
		Runtime:   submission.Runtime,
		ErrorMsg:  submission.ErrorMsg,
		CreatedAt: submission.CreatedAt,
	}, nil
}

// GetMatchSubmissions returns all submissions for a match
func (s *JudgeService) GetMatchSubmissions(matchID uuid.UUID) ([]*SubmissionData, error) {
	var submissions []database.Submission
	if err := s.db.Preload("Player").Where("match_id = ?", matchID).Find(&submissions).Error; err != nil {
		return nil, err
	}

	result := make([]*SubmissionData, len(submissions))
	for i, submission := range submissions {
		result[i] = &SubmissionData{
			ID:        submission.ID,
			MatchID:   submission.MatchID,
			PlayerID:  submission.PlayerID,
			Code:      submission.Code,
			Language:  submission.Language,
			Status:    submission.Status,
			Score:     submission.Score,
			Runtime:   submission.Runtime,
			ErrorMsg:  submission.ErrorMsg,
			CreatedAt: submission.CreatedAt,
		}
	}

	return result, nil
}
