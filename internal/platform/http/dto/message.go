package dto

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
