package models

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type Auth struct {
	Cookies []Cookie `json:"cookies"`
}

type Project struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
	LastUpdated      string `json:"lastUpdated"`
}

type SendableMessage struct {
	Type        string `json:"type"`
	TemplateRef string `json:"templateRef"`
}

type MessageType int

const (
	Welcome MessageType = iota
	Reminder
)

func (m MessageType) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case Reminder:
		return "reminder"
	default:
		return "unknown"
	}
}

// ======

type CampaignResponse []struct {
	Campaign Campaign `json:"campaign"`
}

type Campaign struct {
	Name                          string            `json:"Name"`
	RecordType                    RecordType        `json:"RecordType"`
	ProgramType                   string            `json:"Program_Type__c"`
	Status                        string            `json:"Status"`
	WebTitle                      string            `json:"Web_Title_FF__c"`
	StartDate                     string            `json:"StartDate"`
	ParentId                      string            `json:"ParentId"`
	ActivityStartTime             string            `json:"Activity_Start_Time__c"`
	EndDate                       string            `json:"EndDate"`
	ActivityEndTime               string            `json:"Activity_End_Time__c"`
	WebPublicationStartDate       string            `json:"Web_Publication_Start_Date__c"`
	WebPublicationEndDate         string            `json:"Web_Publication_End_Date__c"`
	SpecialProject                string            `json:"Special_Project__c"`
	CommittedProjectDateRange     *string           `json:"Committed_Project_Date_Range__c"`
	Borough                       string            `json:"Borough__c"`
	WebsiteAddress                string            `json:"Website_Address__c"`
	TeamLeaderToolsId             *string           `json:"Team_Leader_Tools_Id__c"`
	DrupalID                      *string           `json:"Drupal_ID__c"`
	AttendanceTaken               bool              `json:"Attendance_Taken__c"`
	ProjectDescription            string            `json:"Project_Description__c"`
	TeamLeaderNotes               *string           `json:"Team_Leader_Notes__c"`
	ProjectLogistics              *string           `json:"Project_Logistics__c"`
	RepOnSite                     RepOnSite         `json:"Rep_on_site__r"`
	CommunityPartnerName          string            `json:"Community_Partner_Name__c"`
	Directions                    string            `json:"Directions__c"`
	FullCapacity                  string            `json:"Full_Capacity__c"`
	AWSChimeChannelArn            string            `json:"AWS_Chime_Channel_Arn__c"`
	TeamLeaderContact             string            `json:"Team_Leader_Contact__c"`
	TeamLeadersList               *string           `json:"Team_Leaders_List__c"`
	PinnedChatMessage             *string           `json:"Pinned_Chat_Message__c"`
	NumOfRegistration             int               `json:"Num_of_Registration__c"`
	CapacityRemaining             int               `json:"Capacity_Remaining__c"`
	Agency                        Agency            `json:"Agency__r"`
	AgencyDescription             *string           `json:"Agency_Description__c"`
	GeneralInterestCampaign       bool              `json:"General_Interest_Campaign__c"`
	AttendanceSignedUpCount       int               `json:"Attendance_Signed_Up_Count__c"`
	DatetimeState                 string            `json:"Datetime_State__c"`
	OrientationVIFNotRequired     bool              `json:"Orientation_VIF_Not_Required__c"`
	HumanReadableDate             string            `json:"Human_Readable_Date__c"`
	Id                            string            `json:"Id"`
	RecordTypeTL                  string            `json:"RecordType__tl"`
	IsTeamLeader                  bool              `json:"IsTeamLeader__tl"`
	SpecialProjectTL              []string          `json:"SpecialProject__tl"`
	IsMultiSession                bool              `json:"IsMultiSession__tl"`
	IsFirstSession                bool              `json:"IsFirstSession__tl"`
	IsCommittedProject            bool              `json:"IsCommittedProject__tl"`
	CommittedProjectDateRangeTL   []string          `json:"CommittedProjectDateRange__tl"`
	IsTeenFriendly                bool              `json:"IsTeenFriendly__tl"`
	IsFamilyFriendly              bool              `json:"IsFamilyFriendly__tl"`
	StartDateTimeTL               string            `json:"StartDateTime__tl"`
	EndDateTimeTL                 string            `json:"EndDateTime__tl"`
	ProjectOccurrenceStateTL      string            `json:"ProjectOccurrenceState__tl"`
	ProjectIsUpcoming             bool              `json:"ProjectIsUpcoming__tl"`
	CommunityPartnerTL            string            `json:"CommunityPartner__tl"`
	SiteLocationTL                string            `json:"SiteLocation__tl"`
	RepOnSiteNameTL               string            `json:"RepOnSiteName__tl"`
	SiteAddressTL                 string            `json:"SiteAddress__tl"`
	AttendanceTakenTL             bool              `json:"Attendance_Taken__tl"`
	SiteDescriptionTL             string            `json:"SiteDescription__tl"`
	CommunityPartnerNameTL        string            `json:"CommunityPartnerName__tl"`
	CommunityPartnerDescriptionTL string            `json:"CommunityPartnerDescription__tl"`
	AWSChimeChannelId             string            `json:"AWS_Chime_Channel_Id__c"`
	CurrentUserIdTL               string            `json:"CurrentUserId__tl"`
	ProjectManagerIdTL            string            `json:"ProjectManagerId__tl"`
	ProjectManagerNameTL          string            `json:"ProjectManagerName__tl"`
	ProjectManagerChannelLinkTL   string            `json:"ProjectManager_ChannelLink__tl"`
	TeamLeaderIdTL                string            `json:"TeamLeaderId__tl"`
	TeamLeaderNameTL              string            `json:"TeamLeaderName__tl"`
	TeamLeaderFirstNameTL         string            `json:"TeamLeaderFirstName__tl"`
	TeamLeaderChannelLinkTL       string            `json:"TeamLeader_ChannelLink__tl"`
	Bookmarked                    bool              `json:"Bookmarked__tl"`
	IsRecent                      bool              `json:"IsRecent__tl"`
	CampaignSiblings              []CampaignSibling `json:"CampaignSiblings__tl"`
}

type RecordType struct {
	Attributes struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"attributes"`
	Name string `json:"Name"`
}

type RepOnSite struct {
	Attributes struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"attributes"`
	Name string `json:"Name"`
}

type Agency struct {
	Attributes struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"attributes"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type CampaignSibling struct {
	Name              string     `json:"Name"`
	RecordType        RecordType `json:"RecordType"`
	StartDate         string     `json:"StartDate"`
	EndDate           string     `json:"EndDate"`
	HumanReadableDate string     `json:"Human_Readable_Date__c"`
	Id                string     `json:"Id"`
	AWSChimeChannelId string     `json:"AWS_Chime_Channel_Id__c"`
}

type MessageResponse struct {
	ChannelArn string `json:"ChannelArn"`
	MessageId  string `json:"MessageId"`
	Metadata   struct {
		StatusCode   int    `json:"statusCode"`
		EffectiveUri string `json:"effectiveUri"`
		Headers      struct {
			Date           string `json:"date"`
			ContentType    string `json:"content-type"`
			ContentLength  string `json:"content-length"`
			Connection     string `json:"connection"`
			XAmznRequestID string `json:"x-amzn-requestid"`
		} `json:"headers"`
		TransferStats struct {
			HTTP []any `json:"http"`
		} `json:"transferStats"`
	} `json:"@metadata"`
}
