BEGIN TRANSACTION;CREATE TABLE "users" (
	"id"	TEXT NOT NULL UNIQUE,
	"DisplayName"	TEXT NOT NULL,
	"Mail"	TEXT NOT NULL UNIQUE,
	"JobTitle"	TEXT,
	"UserPrincipalName"	TEXT,
	"AccessToken"	TEXT NOT NULL UNIQUE,
	"AccessTokenActive"	INTEGER,
	"RefreshToken"	TEXT,
	PRIMARY KEY("id")
)