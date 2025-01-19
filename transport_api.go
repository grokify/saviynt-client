package saviynt

import (
	"errors"
	"net/http"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/strconv/strconvutil"
)

type TransportService struct {
	client *Client
}

func NewTransportService(client *Client) *TransportService {
	return &TransportService{client: client}
}

type ExportPackageOpts struct {
	// ExportOnline (MANDATORY)
	ExportOnline bool `json:"exportonline"`

	// ExportPath (MANDATORY if ExportOnline is `false`)
	ExportPath string `json:"exportpath,omitempty"`

	// EnvironmentName (MANDATORY if ExportOnline is `true`)
	EnvironmentName string `json:"environmentname,omitempty"`

	// UpdateUser specifies the username of the user exporting the package
	UpdateUser            string                   `json:"updateuser,omitempty"`
	BusinessJustification string                   `json:"businessjustification,omitempty"`
	ObjectsToExport       TransportObjectsToExport `json:"objectstoexport,omitempty"`

	// TransportOwner is an option option to transport owners for selected object
	TransportOwner   bool `json:"transportowner"`
	TransportMembers bool `json:"transportmembers"`
}

func (opts ExportPackageOpts) internal() (exportPackageOptInternal, error) {
	if opts.ExportOnline && strings.TrimSpace(opts.EnvironmentName) == "" {
		return exportPackageOptInternal{}, errors.New("environmentname is mandatory if exportonline is \"true\"")
	} else if !opts.ExportOnline && strings.TrimSpace(opts.ExportPath) == "" {
		return exportPackageOptInternal{}, errors.New("exportpoath is mandatory if exportonline is \"false\"")
	}

	return exportPackageOptInternal{
		ExportOnline:          strconvutil.Btoa(opts.ExportOnline),
		ExportPath:            opts.ExportPath,
		EnvironmentName:       opts.EnvironmentName,
		UpdateUser:            opts.UpdateUser,
		BusinessJustification: opts.BusinessJustification,
		ObjectsToExport:       opts.ObjectsToExport,
		TransportOwner:        strconvutil.Btoa(opts.TransportOwner),
		TransportMembers:      strconvutil.Btoa(opts.TransportMembers),
	}, nil
}

type exportPackageOptInternal struct {
	ExportOnline          string                   `json:"exportonline"`
	ExportPath            string                   `json:"exportpath,omitempty"`
	EnvironmentName       string                   `json:"environmentname,omitempty"`
	UpdateUser            string                   `json:"updateuser,omitempty"`
	BusinessJustification string                   `json:"businessjustification,omitempty"`
	ObjectsToExport       TransportObjectsToExport `json:"objectstoexport,omitempty"`
	TransportOwner        string                   `json:"transportowner"`
	TransportMembers      string                   `json:"transportmembers"`
}

type TransportObjectsToExport struct {
	SAVRoles        []string `json:"savRoles,omitempty"`
	EmailTemplate   []string `json:"emailTemplate,omitempty"`
	Roles           []string `json:"roles,omitempty"`
	AnalyticsV1     []string `json:"analyticsV1,omitempty"`
	AnalyticsV2     []string `json:"analyticsV2,omitempty"`
	GlobalConfig    []string `json:"globalConfig,omitempty"`
	Workflows       []string `json:"workflows,omitempty"`
	Connection      []string `json:"connection,omitempty"`
	AppOnboarding   []string `json:"appOnboarding,omitempty"`
	UserGroups      []string `json:"userGroups,omitempty"`
	ScanRules       []string `json:"scanRules,omitempty"`
	Organizations   []string `json:"organizations,omitempty"`
	SecuritySystems []string `json:"securitySystems,omitempty"`
}

func (svc *TransportService) ExportPackage(opts ExportPackageOpts) (*httpsimple.Request, *http.Response, error) {
	if svc.client == nil {
		return nil, nil, ErrClientNotSet
	} else if svc.client.SimpleClient == nil {
		return nil, nil, ErrSimpleClientNotSet
	}
	opts2, err := opts.internal()
	if err != nil {
		return nil, nil, err
	}

	sreq := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      svc.client.BuildURL(RelURLTransportExport),
		BodyType: httpsimple.BodyTypeJSON,
		Body:     opts2,
	}

	resp, err := svc.client.SimpleClient.Do(sreq)
	return &sreq, resp, err
}

type ImportPackageOpts struct {
	// PackageToImport (MANDATORY) specifies the local filepath of the package to import.
	PackageToImport string `json:"packagetoimport"`

	// UpdateUser (OPTIOAL) specifies the username of the user importing the package,
	UpdateUser            string `json:"updateuser,omitempty"`
	BusinessJustification string `json:"businessjustification,omitempty"`
}

func (svc *TransportService) ImportPackage(opts ImportPackageOpts) (*httpsimple.Request, *http.Response, error) {
	if svc.client == nil {
		return nil, nil, ErrClientNotSet
	} else if svc.client.SimpleClient == nil {
		return nil, nil, ErrSimpleClientNotSet
	}

	sreq := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      svc.client.BuildURL(RelURLTransportImport),
		BodyType: httpsimple.BodyTypeJSON,
		Body:     opts,
	}

	resp, err := svc.client.SimpleClient.Do(sreq)
	return &sreq, resp, err
}
