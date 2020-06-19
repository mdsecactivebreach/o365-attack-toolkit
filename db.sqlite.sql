BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (
	"id"	TEXT NOT NULL UNIQUE,
	"DisplayName"	TEXT NOT NULL,
	"Mail"	TEXT NOT NULL UNIQUE,
	"JobTitle"	TEXT,
	"UserPrincipalName"	TEXT,
	"AccessToken"	TEXT NOT NULL UNIQUE,
	"AccessTokenActive"	INTEGER,
	"RefreshToken"	TEXT,
	PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "mails" (
	"Id"	TEXT NOT NULL UNIQUE,
	"User"	TEXT NOT NULL,
	"Subject"	TEXT,
	"SenderEmail"	TEXT NOT NULL,
	"SenderName"	TEXT NOT NULL,
	"HasAttachments"	INTEGER,
	"BodyPreview"	TEXT NOT NULL,
	"BodyType"	TEXT NOT NULL,
	"BodyContent"	TEXT NOT NULL,
	PRIMARY KEY("Id")
);
COMMIT;
