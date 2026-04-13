package bedrockservice

import (
	"context"
	"fmt"
)

type MockBedrockService struct{}

func NewMockBedrockService() *MockBedrockService {
	return &MockBedrockService{}
}

func (s *MockBedrockService) GenerateThankYouMessage(_ context.Context, _, projectName, _ string) (string, error) {
	return fmt.Sprintf("Thank you so much for leading \"%s\" today — your dedication makes a real difference!", projectName), nil
}
