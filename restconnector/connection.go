package restconnector

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type ConnectionFile map[string]Connection

func (cm ConnectionFile) Names() []string {
	return maputil.Keys(cm)
}

func (cm ConnectionFile) WriteExternalAttributeJSONFiles(dir string, templateNameAsDir, templateNameAsFilePrefix, writeHTTPParamsFile bool) error {
	for _, c := range cm {
		err := c.WriteExternalAttributeJSONFiles(dir, templateNameAsDir, templateNameAsFilePrefix, writeHTTPParamsFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadConnectionFile(filename string) (*ConnectionFile, error) {
	connMap := &ConnectionFile{}
	_, err := jsonutil.UnmarshalFile(filename, connMap)
	return connMap, err
}

type Connection struct {
	ConnectionKey             string        `json:"_connectionKey"`
	ConnectionDescription     string        `json:"connectiondescription"`
	ConnectionName            string        `json:"connectionname"`
	Connectorms               bool          `json:"connectorms"`
	ConnectorType             *string       `json:"connectorType"`
	CredentialChangeConfig    *string       `json:"credentialChangeConfig"`
	ExternalConnectionType    string        `json:"externalconnectiontype"`
	MSConnectorVersion        *string       `json:"msconnectorversion"`
	SSLCertificate            *string       `json:"sslCertificate"`
	SSLCertificateID          *string       `json:"sslCertificateId"`
	Status                    int           `json:"status"`
	StatusForEnableDisable    int           `json:"statusForEnableDisable"`
	TemplateMandatoryData     string        `json:"templateMandatoryData"`
	TemplateName              string        `json:"templateName"`
	VaultConfig               *string       `json:"vaultConfig"`
	VaultCredentialConnection *string       `json:"vaultCredentialConnection"`
	ExternalAttrs             ExternalAttrs `json:"EXTERNAL_ATTR"`
}

type ExternalAttrs []ExternalAttr

func (eas ExternalAttrs) Inflate() ExternalAttrs {
	inflated := ExternalAttrs{}
	namesMap := map[string]int{}
	for _, ea := range eas {
		inflated = append(inflated, ea)
		namesMap[ea.AttributeName]++
	}
	allNames := ExternalAttributeNames()
	for _, extraName := range allNames {
		if _, ok := namesMap[extraName]; !ok {
			inflated = append(inflated, ExternalAttr{
				AttributeName: extraName,
			})
		}
	}
	return inflated
}

func (eas ExternalAttrs) CallBodies() []string {
	var bodies []string
	for _, ea := range eas {
		eav := strings.TrimSpace(ea.EncryptedAttributeValue)
		if strings.Index(eav, "{") != 0 {
			continue
		}
		ci := &CallInfo{}
		err := json.Unmarshal([]byte(eav), ci)
		if err != nil {
			continue
		}
		cb := ci.CallBodies()
		if len(cb) > 0 {
			bodies = append(bodies, cb...)
		}
	}
	return bodies
}

func (eas ExternalAttrs) Markdown() []string {
	lines := []string{}

	for _, ea := range eas {
		eav := strings.TrimSpace(ea.EncryptedAttributeValue)
		if strings.Index(eav, "{") == 0 {
			fmt.Println(eav)
		}
	}

	return lines
}

type ExternalAttrNamesOpts struct {
	ToUpper      bool
	Dedupe       bool
	Sort         bool
	RequireValue bool
}

func (e ExternalAttrs) Names(opts ExternalAttrNamesOpts) []string {
	var names []string
	for _, ea := range e {
		name := ea.AttributeName
		name = strings.TrimSpace(name)
		if opts.ToUpper {
			name = strings.ToUpper(name)
		}
		if opts.RequireValue {
			if v := strings.TrimSpace(ea.EncryptedAttributeValue); v == "" {
				continue
			}
		}
		if name != "" {
			names = append(names, name)
		}
	}
	if opts.Dedupe {
		names = slicesutil.Dedupe(names)
	}
	if opts.Sort {
		sort.Strings(names)
	}
	return names
}

func (c Connection) WriteExternalAttributeJSONFiles(dir string, templateNameAsDir, templateNameAsFilePrefix, writeHTTPParamsFile bool) error {
	for _, ea := range c.ExternalAttrs {
		if !ea.HasEncryptedAttributeValueJSONMap() {
			continue
		}
		outfileBase := strings.TrimSpace(ea.CanonicalAttrributeName())
		if outfileBase == "" {
			return errors.New("no attribute name")
		}
		tmplName := strings.TrimSpace(c.TemplateName)
		outfileExtAttrJSON := outfileBase
		if templateNameAsFilePrefix {
			outfileBase = strings.Join([]string{tmplName, outfileBase}, ".")
		}
		err := jsonutil.WriteFileIndentBytes(outfileExtAttrJSON+".json", []byte(ea.EncryptedAttributeValue), "", "  ", 0600)
		//err = os.WriteFile(outfileExtAttrJSON, out, 0600)
		if err != nil {
			return err
		}
		if writeHTTPParamsFile {

		}
	}
	return nil
}

type ExternalAttr struct {
	AttributeName           string  `json:"attributename"`
	CredentialSaveStatus    bool    `json:"credentialSaveStatus"`
	EncryptedAttributeValue string  `json:"encryptedattributevalue"`
	Formdata                *string `json:"formdata"`
}

func (ea ExternalAttr) CanonicalAttrributeName() string {
	return CanonicalizeExternalAttributeName(ea.AttributeName)
}

func (ea ExternalAttr) HasEncryptedAttributeValue() bool {
	return strings.TrimSpace(ea.EncryptedAttributeValue) != ""
}

func (ea ExternalAttr) HasEncryptedAttributeValueJSONMap() bool {
	return strings.Index(strings.TrimSpace(ea.EncryptedAttributeValue), "{") == 0 &&
		stringsutil.ReverseIndex(strings.TrimSpace(ea.EncryptedAttributeValue), "}") == 0
}

func (ea ExternalAttr) ParseEncryptedAttributeValueJSONMap() (ExternalAttributeValue, error) {
	if !ea.HasEncryptedAttributeValueJSONMap() {
		return ExternalAttributeValue{}, errors.New("encryptedAttributeValue JSON map not present")
	}
	eav := ExternalAttributeValue{}
	return eav, json.Unmarshal([]byte(ea.EncryptedAttributeValue), &eav)
}

func CanonicalizeExternalAttributeName(s string) string {
	m := ExternalAttributeNamesCanonicalMap()
	if v, ok := m[s]; ok {
		return v
	}
	return s
}

func ExternalAttributeNamesCanonicalMap() map[string]string {
	m := map[string]string{}
	names := ExternalAttributeNames()
	for _, name := range names {
		nameLc := strings.TrimSpace(strings.ToLower(name))
		m[nameLc] = name
	}
	return m
}

func ExternalAttributeNames() []string {
	return []string{
		"AddAccessJSON",
		"AddFFIDAccessJSON",
		"ChangePassJSON",
		"ConfigJSON",
		"ConnectionJSON",
		"CreateAccountJSON",
		"CreateTicketJSON",
		"DisableAccountJSON",
		"EnableAccountJSON",
		"ImportAccountEntJSON",
		"ImportUserJSON",
		"PasswdPolicyJSON",
		"RemoveAccessJSON",
		"RemoveAccountJSON",
		"RemoveFFIDAccessJSON",
		"SendOtpJSON",
		"TicketStatusJSON",
		"UpdateAccountJSON",
		"UpdateUserJSON",
		"ValidateOtpJSON",
		"ENDPOINTS_FILTER",
		// "MODIFYUSERDATAJSON",
		"ModifyUserDataJSON",
		"PAM_CONFIG",
		"STATUS_THRESHOLD_CONFIG",
	}
}
