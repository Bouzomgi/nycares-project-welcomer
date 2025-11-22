package dto

////// GetCampaign

type CampaignResponse []CampaignWrapper

type CampaignWrapper struct {
	Campaign Campaign `json:"campaign"`
}

type Campaign struct {
	Name                          string     `json:"Name"`
	RecordType                    RecordType `json:"RecordType"`
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
	RepOnSite                     Contact    `json:"Rep_on_site__r"`
	CommunityPartnerName          string     `json:"Community_Partner_Name__c"`
	Directions                    string     `json:"Directions__c"`
	FullCapacity                  string     `json:"Full_Capacity__c"`
	AWSChimeChannelArn            string     `json:"AWS_Chime_Channel_Arn__c"`
	TeamLeaderContact             string     `json:"Team_Leader_Contact__c"`
	TeamLeadersList               *string    `json:"Team_Leaders_List__c"`
	PinnedChatMessage             *string    `json:"Pinned_Chat_Message__c"`
	NumOfRegistration             int        `json:"Num_of_Registration__c"`
	CapacityRemaining             int        `json:"Capacity_Remaining__c"`
	Agency                        Agency     `json:"agency__r"`
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
	CampaignSiblings              []Campaign `json:"CampaignSiblings__tl"`
}

type RecordType struct {
	Attributes RecordTypeAttributes `json:"attributes"`
	Name       string               `json:"Name"`
}

type RecordTypeAttributes struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Contact struct {
	Attributes RecordTypeAttributes `json:"attributes"`
	Name       string               `json:"Name"`
}

type Agency struct {
	Attributes  RecordTypeAttributes `json:"attributes"`
	Name        string               `json:"Name"`
	Description string               `json:"Description"`
}

////// SendMessage

type SendMessageResponse struct {
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

/////// GetMessages

type ChannelMessagesResponse []struct {
	ChannelMessages []ChannelMessage `json:"ChannelMessages"`
}

type ChannelMessage struct {
	MessageId            string          `json:"MessageId"`
	Content              string          `json:"Content"`
	Metadata             MessageMetadata `json:"Metadata"`
	Type                 string          `json:"Type"`
	CreatedTimestamp     string          `json:"CreatedTimestamp"`
	LastUpdatedTimestamp string          `json:"LastUpdatedTimestamp"`
	Sender               MessageSender   `json:"Sender"`
	Redacted             bool            `json:"Redacted"`
	ChannelType          string          `json:"ChannelType"`
	ProjectIdTL          string          `json:"Project_Id__tl"`
	JSONRepresentationTL string          `json:"JSON_Representation__tl"`
}

type MessageMetadata struct {
	SenderSFIDTL string        `json:"Sender_SFID__tl"`
	Attachments  []interface{} `json:"attachments"`
}

type MessageSender struct {
	Arn                  string `json:"Arn"`
	Name                 string `json:"Name"`
	RoleC                string `json:"Role__c"`
	FirstNameLastInitial string `json:"FirstNameLastInitial__tl"`
	PhotoURL             string `json:"PhotoUrl__c"`
}

/////// PostPinMessage

type PostPinResponse struct {
	Scalar bool `json:"scalar"`
}
