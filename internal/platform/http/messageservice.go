package httpservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
)

type MessageService interface {
	GetProjectChannelId(ctx context.Context, projectId string) (string, error)
	SendMessage(ctx context.Context, channelId, messageContent string) (string, error)
	PinMessage(ctx context.Context, channelId, messageId string) error
	SetCookies(cookies []*http.Cookie) error
}

type campaignResponse []campaignWrapper

type campaignWrapper struct {
	Campaign campaign `json:"campaign"`
}

type campaign struct {
	Name                          string     `json:"Name"`
	RecordType                    recordType `json:"RecordType"`
	ProgramType                   string     `json:"Program_Type__c"`
	Status                        string     `json:"Status"`
	WebTitleFF                    string     `json:"Web_Title_FF__c"`
	StartDate                     string     `json:"StartDate"`
	ParentId                      string     `json:"ParentId"`
	ActivityStartTime             string     `json:"Activity_Start_Time__c"`
	EndDate                       string     `json:"EndDate"`
	ActivityEndTime               string     `json:"Activity_End_Time__c"`
	WebPublicationStartDate       string     `json:"Web_Publication_Start_Date__c"`
	WebPublicationEndDate         string     `json:"Web_Publication_End_Date__c"`
	SpecialProject                string     `json:"Special_Project__c"`
	CommittedProjectDateRange     *string    `json:"Committed_Project_Date_Range__c"`
	Borough                       string     `json:"Borough__c"`
	WebsiteAddress                string     `json:"Website_Address__c"`
	TeamLeaderToolsId             *string    `json:"Team_Leader_Tools_Id__c"`
	DrupalID                      *string    `json:"Drupal_ID__c"`
	AttendanceTaken               bool       `json:"Attendance_Taken__c"`
	ProjectDescription            string     `json:"Project_Description__c"`
	TeamLeaderNotes               *string    `json:"Team_Leader_Notes__c"`
	ProjectLogistics              *string    `json:"Project_Logistics__c"`
	RepOnSite                     contact    `json:"Rep_on_site__r"`
	CommunityPartnerName          string     `json:"Community_Partner_Name__c"`
	Directions                    string     `json:"Directions__c"`
	FullCapacity                  string     `json:"Full_Capacity__c"`
	AWSChimeChannelArn            string     `json:"AWS_Chime_Channel_Arn__c"`
	TeamLeaderContact             string     `json:"Team_Leader_Contact__c"`
	TeamLeadersList               *string    `json:"Team_Leaders_List__c"`
	PinnedChatMessage             *string    `json:"Pinned_Chat_Message__c"`
	NumOfRegistration             int        `json:"Num_of_Registration__c"`
	CapacityRemaining             int        `json:"Capacity_Remaining__c"`
	Agency                        agency     `json:"agency__r"`
	AgencyDescription             *string    `json:"agency_Description__c"`
	GeneralInterestCampaign       bool       `json:"General_Interest_Campaign__c"`
	AttendanceSignedUpCount       int        `json:"Attendance_Signed_Up_Count__c"`
	DatetimeState                 string     `json:"Datetime_State__c"`
	OrientationVIFNotRequired     bool       `json:"Orientation_VIF_Not_Required__c"`
	HumanReadableDate             string     `json:"Human_Readable_Date__c"`
	Id                            string     `json:"Id"`
	RegistrationId                string     `json:"Registration_Id__tl"`
	UserStatus                    string     `json:"UserStatus__tl"`
	UserRole                      string     `json:"UserRole__tl"`
	RecordTypeTL                  string     `json:"RecordType__tl"`
	IsTeamLeader                  bool       `json:"IsTeamLeader__tl"`
	SpecialProjectTL              []string   `json:"SpecialProject__tl"`
	IsMultiSession                bool       `json:"IsMultiSession__tl"`
	IsFirstSession                bool       `json:"IsFirstSession__tl"`
	IsCommittedProject            bool       `json:"IsCommittedProject__tl"`
	CommittedProjectDateRangeTL   []string   `json:"CommittedProjectDateRange__tl"`
	IsTeenFriendly                bool       `json:"IsTeenFriendly__tl"`
	IsFamilyFriendly              bool       `json:"IsFamilyFriendly__tl"`
	StartDateTimeTL               string     `json:"StartDateTime__tl"`
	EndDateTimeTL                 string     `json:"EndDateTime__tl"`
	ProjectOccurrenceState        string     `json:"ProjectOccurrenceState__tl"`
	ProjectIsUpcoming             bool       `json:"ProjectIsUpcoming__tl"`
	CommunityPartnerTL            string     `json:"CommunityPartner__tl"`
	SiteLocation                  string     `json:"SiteLocation__tl"`
	RepOnSiteName                 string     `json:"RepOnSiteName__tl"`
	SiteAddressTL                 string     `json:"SiteAddress__tl"`
	AttendanceTakenTL             bool       `json:"Attendance_Taken__tl"`
	SiteDescription               string     `json:"SiteDescription__tl"`
	CommunityPartnerNameTL        string     `json:"CommunityPartnerName__tl"`
	CommunityPartnerDescriptionTL string     `json:"CommunityPartnerDescription__tl"`
	AWSChimeChannelID             string     `json:"AWS_Chime_Channel_Id__c"`
	CurrentUserID                 string     `json:"CurrentUserId__tl"`
	ProjectManagerID              string     `json:"ProjectManagerId__tl"`
	ProjectManagerName            string     `json:"ProjectManagerName__tl"`
	ProjectManagerChannelLink     string     `json:"ProjectManager_ChannelLink__tl"`
	TeamLeaderID                  string     `json:"TeamLeaderId__tl"`
	TeamLeaderName                string     `json:"TeamLeaderName__tl"`
	TeamLeaderFirstName           string     `json:"TeamLeaderFirstName__tl"`
	TeamLeaderChannelLink         string     `json:"TeamLeader_ChannelLink__tl"`
	Bookmarked                    bool       `json:"Bookmarked__tl"`
	IsRecent                      bool       `json:"IsRecent__tl"`
	CampaignSiblings              []campaign `json:"CampaignSiblings__tl"`
}

type recordType struct {
	Attributes recordTypeAttributes `json:"attributes"`
	Name       string               `json:"Name"`
}

type recordTypeAttributes struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type contact struct {
	Attributes recordTypeAttributes `json:"attributes"`
	Name       string               `json:"Name"`
}

type agency struct {
	Attributes  recordTypeAttributes `json:"attributes"`
	Name        string               `json:"Name"`
	Description string               `json:"Description"`
}

func (s *HttpService) GetProjectChannelId(ctx context.Context, projectId string) (string, error) {

	req, err := s.buildCampaignRequest(projectId)
	if err != nil {
		return "", fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var campaignResp campaignResponse
	if err := json.Unmarshal(body, &campaignResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign response: %w", err)
	}

	if len(campaignResp) == 0 {
		return "", fmt.Errorf("campaign response was empty: %w", err)
	}

	channelId := campaignResp[0].Campaign.AWSChimeChannelID

	return channelId, nil
}

func (s *HttpService) buildCampaignRequest(projectId string) (*http.Request, error) {
	getCampaignBaseUrl := endpoints.JoinPaths(endpoints.BaseUrl, endpoints.GetCampaignPath)
	urlStr := fmt.Sprintf("%s/%s", getCampaignBaseUrl, projectId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return req, nil
}

/////////

type sendMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	ChannelArn string   `json:"ChannelArn"`
	MessageId  string   `json:"MessageId"`
	Metadata   Metadata `json:"@metadata"`
}

type Metadata struct {
	StatusCode    int           `json:"statusCode"`
	EffectiveURI  string        `json:"effectiveUri"`
	Headers       Headers       `json:"headers"`
	TransferStats TransferStats `json:"transferStats"`
}

type Headers struct {
	Date           string `json:"date"`
	ContentType    string `json:"content-type"`
	ContentLength  string `json:"content-length"`
	Connection     string `json:"connection"`
	XAmznRequestID string `json:"x-amzn-requestid"`
}

type TransferStats struct {
	HTTP [][]interface{} `json:"http"`
}

func (s *HttpService) SendMessage(ctx context.Context, channelId, messageContent string) (string, error) {
	req, err := s.buildSendMessageRequest(channelId, messageContent)
	if err != nil {
		return "", fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return "", fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var sendMessageResp sendMessageResponse
	if err := json.Unmarshal(body, &sendMessageResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal campaign response: %w", err)
	}

	messageId := sendMessageResp.Data.MessageId
	return messageId, nil
}

func (s *HttpService) buildSendMessageRequest(channelId, messageContent string) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormField("message")
	if err != nil {
		return nil, err
	}

	io.WriteString(part, messageContent)
	writer.Close()

	urlStr := fmt.Sprintf("%s/api/messenger/channel/%s/message/post", endpoints.BaseUrl, channelId)

	req, err := http.NewRequest("POST", urlStr, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	return req, nil
}

////////

func (s *HttpService) PinMessage(ctx context.Context, channelId, messageId string) error {
	req, err := s.buildPinMessageRequest(channelId, messageId)
	if err != nil {
		return fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return fmt.Errorf("schedule request failed: %w", err)
	}

	return nil
}

func (s *HttpService) buildPinMessageRequest(channelId, messageId string) (*http.Request, error) {

	body := map[string]string{
		"MessageId": messageId,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	urlStr := fmt.Sprintf("%s/api/messenger/create-pin-message/%s", endpoints.BaseUrl, channelId)

	req, err := http.NewRequest(
		"POST",
		urlStr,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

	return req, nil
}
