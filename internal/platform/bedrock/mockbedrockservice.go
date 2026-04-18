package bedrockservice

import "context"

type MockBedrockService struct{}

func NewMockBedrockService() *MockBedrockService {
	return &MockBedrockService{}
}

func (s *MockBedrockService) GenerateThankYouMessage(_ context.Context, _, _ string) (string, error) {
	return "Thank you so much for your help today — your dedication makes a real difference!", nil
}
