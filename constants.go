package saviynt

const (
	RelURLLogin         = "/ECM/api/login"
	RelOAuthAccessToken = "/ECM/oauth/access_token" // #nosec G101
	RelURLECM           = "/ECM"
	RelURLAPI           = "/api/v5"

	RelURLLoginRuntimeControlsData = "/fetchRuntimeControlsDataV2" // API at https://documenter.getpostman.com/view/23973797/2s9XxwutWR#b821cc21-ee7c-49e3-9433-989ed87b2b03
	RelURLUserGetAccessDetails     = "/getAccessDetailsForUser"
	RelURLPasswordChange           = "/changePassword"
	RelURLUserGet                  = "/getUser"
	RelURLUserUpdate               = "/updateUser"

	EnvSaviyntServerURL = "SAVIYNT_SERVER_URL"
	EnvSaviyntUsername  = "SAVIYNT_USERNAME"
	EnvSaviyntPassword  = "SAVIYNT_PASSWORD"
)
