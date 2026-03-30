package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes/mockresponses"
	"github.com/gorilla/mux"
)

const mockStateKey = "PROJECTS"

var (
	// in-memory fallback (used when MOCK_STATE_TABLE is unset)
	adminMu       sync.RWMutex
	adminProjects []mockresponses.ProjectConfig

	// DynamoDB client, initialized once
	dynamoOnce   sync.Once
	dynamoClient *dynamodb.Client
)

func initDynamo() {
	dynamoOnce.Do(func() {
		cfgOpts := []func(*config.LoadOptions) error{
			config.WithRegion("us-east-1"),
		}
		cfg, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
		if err != nil {
			panic(fmt.Sprintf("mockserver: failed to load AWS config: %v", err))
		}
		var opts []func(*dynamodb.Options)
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			opts = append(opts, func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(endpoint)
			})
		}
		dynamoClient = dynamodb.NewFromConfig(cfg, opts...)
	})
}

func GetAdminProjects() ([]mockresponses.ProjectConfig, error) {
	tableName := os.Getenv("MOCK_STATE_TABLE")
	if tableName == "" {
		adminMu.RLock()
		defer adminMu.RUnlock()
		return adminProjects, nil
	}

	initDynamo()

	result, err := dynamoClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName:      aws.String(tableName),
		ConsistentRead: aws.Bool(true),
		Key: map[string]types.AttributeValue{
			"MockKey": &types.AttributeValueMemberS{Value: mockStateKey},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("GetAdminProjects: DynamoDB GetItem failed: %w", err)
	}
	if result.Item == nil {
		return nil, nil
	}

	av, ok := result.Item["ProjectsJSON"]
	if !ok {
		return nil, nil
	}
	jsonAttr, ok := av.(*types.AttributeValueMemberS)
	if !ok {
		return nil, fmt.Errorf("GetAdminProjects: ProjectsJSON attribute is not a string")
	}

	var stored []projectInput
	if err := json.Unmarshal([]byte(jsonAttr.Value), &stored); err != nil {
		return nil, fmt.Errorf("GetAdminProjects: failed to unmarshal ProjectsJSON: %w", err)
	}

	projects := make([]mockresponses.ProjectConfig, 0, len(stored))
	for _, p := range stored {
		date, err := time.Parse("2006-01-02", p.Date)
		if err != nil {
			return nil, fmt.Errorf("GetAdminProjects: invalid date %q: %w", p.Date, err)
		}
		projects = append(projects, mockresponses.ProjectConfig{
			Name:       p.Name,
			Date:       date,
			Id:         p.Id,
			CampaignId: p.CampaignId,
		})
	}
	return projects, nil
}

func SetAdminProjects(projects []mockresponses.ProjectConfig) error {
	tableName := os.Getenv("MOCK_STATE_TABLE")
	if tableName == "" {
		adminMu.Lock()
		defer adminMu.Unlock()
		adminProjects = projects
		return nil
	}

	initDynamo()

	stored := make([]projectInput, 0, len(projects))
	for _, p := range projects {
		stored = append(stored, projectInput{
			Name:       p.Name,
			Date:       p.Date.Format("2006-01-02"),
			Id:         p.Id,
			CampaignId: p.CampaignId,
		})
	}
	data, err := json.Marshal(stored)
	if err != nil {
		return fmt.Errorf("SetAdminProjects: failed to marshal projects: %w", err)
	}

	_, err = dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"MockKey":      &types.AttributeValueMemberS{Value: mockStateKey},
			"ProjectsJSON": &types.AttributeValueMemberS{Value: string(data)},
		},
	})
	if err != nil {
		return fmt.Errorf("SetAdminProjects: DynamoDB PutItem failed: %w", err)
	}
	return nil
}

type setProjectsRequest struct {
	Projects []projectInput `json:"projects"`
}

type projectInput struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Id         string `json:"id"`
	CampaignId string `json:"campaignId"`
}

func RegisterAdminRoutes(r *mux.Router) {
	r.HandleFunc("/admin/set-projects", func(w http.ResponseWriter, r *http.Request) {
		var req setProjectsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var projects []mockresponses.ProjectConfig
		for _, p := range req.Projects {
			date, err := time.Parse("2006-01-02", p.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid date: " + p.Date})
				return
			}
			projects = append(projects, mockresponses.ProjectConfig{
				Name:       p.Name,
				Date:       date,
				Id:         p.Id,
				CampaignId: p.CampaignId,
			})
		}

		if err := SetAdminProjects(projects); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("POST")
}
