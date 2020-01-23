package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	duoapi "github.com/duosecurity/duo_api_golang"
)

func buildAdminClient(url string, proxy func(*http.Request) (*url.URL, error)) *Client {
	ikey := "eyekey"
	skey := "esskey"
	host := strings.Split(url, "//")[1]
	userAgent := "GoTestClient"
	base := duoapi.NewDuoApi(ikey, skey, host, userAgent, duoapi.SetTimeout(1*time.Second), duoapi.SetInsecure(), duoapi.SetProxy(proxy))
	return New(*base)
}

func getBodyParams(r *http.Request) (url.Values, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return url.Values{}, err
	}
	reqParams, err := url.ParseQuery(string(body))
	return reqParams, err
}

const getUsersResponse = `{
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 1
	},
	"response": [{
		"alias1": "joe.smith",
		"alias2": "jsmith@example.com",
		"alias3": null,
		"alias4": null,
		"created": 1489612729,
		"email": "jsmith@example.com",
		"firstname": "Joe",
		"groups": [{
			"desc": "People with hardware tokens",
			"name": "token_users"
		}],
		"last_directory_sync": 1508789163,
		"last_login": 1343921403,
		"lastname": "Smith",
		"notes": "",
		"phones": [{
			"phone_id": "DPFZRS9FB0D46QFTM899",
			"number": "+15555550100",
			"extension": "",
			"name": "",
			"postdelay": null,
			"predelay": null,
			"type": "Mobile",
			"capabilities": [
				"sms",
				"phone",
				"push"
			],
			"platform": "Apple iOS",
			"activated": false,
			"sms_passcodes_sent": false
		}],
		"realname": "Joe Smith",
		"status": "active",
		"tokens": [{
			"serial": "0",
			"token_id": "DHIZ34ALBA2445ND4AI2",
			"type": "d1"
		}],
		"user_id": "DU3RP9I2WOC59VZX672N",
		"username": "jsmith"
	}]
}`

const getUserResponse = `{
	"stat": "OK",
	"response": {
		"alias1": "joe.smith",
		"alias2": "jsmith@example.com",
		"alias3": null,
		"alias4": null,
		"created": 1489612729,
		"email": "jsmith@example.com",
		"firstname": "Joe",
		"groups": [{
			"desc": "People with hardware tokens",
			"name": "token_users"
		}],
		"last_directory_sync": 1508789163,
		"last_login": 1343921403,
		"lastname": "Smith",
		"notes": "",
		"phones": [{
			"phone_id": "DPFZRS9FB0D46QFTM899",
			"number": "+15555550100",
			"extension": "",
			"name": "",
			"postdelay": null,
			"predelay": null,
			"type": "Mobile",
			"capabilities": [
				"sms",
				"phone",
				"push"
			],
			"platform": "Apple iOS",
			"activated": false,
			"sms_passcodes_sent": false
		}],
		"realname": "Joe Smith",
		"status": "active",
		"tokens": [{
			"serial": "0",
			"token_id": "DHIZ34ALBA2445ND4AI2",
			"type": "d1"
		}],
		"user_id": "DU3RP9I2WOC59VZX672N",
		"username": "jsmith"
	}
}`

func TestGetUsers(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getUsersResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUsers()
	if err != nil {
		t.Errorf("Unexpected error from GetUsers call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 user, but got %d", len(result.Response))
	}
	if result.Response[0].UserID != "DU3RP9I2WOC59VZX672N" {
		t.Errorf("Expected user ID DU3RP9I2WOC59VZX672N, but got %s", result.Response[0].UserID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUsersPage1Response = `{
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": 1,
		"total_objects": 2
	},
	"response": [{
		"alias1": "joe.smith",
		"alias2": "jsmith@example.com",
		"alias3": null,
		"alias4": null,
		"created": 1489612729,
		"email": "jsmith@example.com",
		"firstname": "Joe",
		"groups": [{
			"desc": "People with hardware tokens",
			"name": "token_users"
		}],
		"last_directory_sync": 1508789163,
		"last_login": 1343921403,
		"lastname": "Smith",
		"notes": "",
		"phones": [{
			"phone_id": "DPFZRS9FB0D46QFTM899",
			"number": "+15555550100",
			"extension": "",
			"name": "",
			"postdelay": null,
			"predelay": null,
			"type": "Mobile",
			"capabilities": [
				"sms",
				"phone",
				"push"
			],
			"platform": "Apple iOS",
			"activated": false,
			"sms_passcodes_sent": false
		}],
		"realname": "Joe Smith",
		"status": "active",
		"tokens": [{
			"serial": "0",
			"token_id": "DHIZ34ALBA2445ND4AI2",
			"type": "d1"
		}],
		"user_id": "DU3RP9I2WOC59VZX672N",
		"username": "jsmith"
	}]
}`

const getUsersPage2Response = `{
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 2
	},
	"response": [{
		"alias1": "joe.smith",
		"alias2": "jsmith@example.com",
		"alias3": null,
		"alias4": null,
		"created": 1489612729,
		"email": "jsmith@example.com",
		"firstname": "Joe",
		"groups": [{
			"desc": "People with hardware tokens",
			"name": "token_users"
		}],
		"last_directory_sync": 1508789163,
		"last_login": 1343921403,
		"lastname": "Smith",
		"notes": "",
		"phones": [{
			"phone_id": "DPFZRS9FB0D46QFTM899",
			"number": "+15555550100",
			"extension": "",
			"name": "",
			"postdelay": null,
			"predelay": null,
			"type": "Mobile",
			"capabilities": [
				"sms",
				"phone",
				"push"
			],
			"platform": "Apple iOS",
			"activated": false,
			"sms_passcodes_sent": false
		}],
		"realname": "Joe Smith",
		"status": "active",
		"tokens": [{
			"serial": "0",
			"token_id": "DHIZ34ALBA2445ND4AI2",
			"type": "d1"
		}],
		"user_id": "DU3RP9I2WOC59VZX672N",
		"username": "jsmith"
	}]
}`

func TestGetUsersMultipage(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getUsersPage1Response)
			} else {
				fmt.Fprintln(w, getUsersPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUsers()

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if result.Metadata.TotalObjects != "2" {
		t.Errorf("Expected total obects to be two, found %s", result.Metadata.TotalObjects)
	}

	if len(result.Response) != 2 {
		t.Errorf("Expected two users in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

const getEmptyPageArgsResponse = `{
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": 2,
		"total_objects": 2
	},
	"response": []
}`

func TestGetUserPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetUsers(func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

func TestGetUser(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getUserResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUser("DU3RP9I2WOC59VZX672N")
	if err != nil {
		t.Errorf("Unexpected error from GetUser call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if result.Response.UserID != "DU3RP9I2WOC59VZX672N" {
		t.Errorf("Expected user ID DU3RP9I2WOC59VZX672N, but got %s", result.Response.UserID)
	}
}

const getGroupsResponse = `{
	"response": [{
		"desc": "This is group A",
		"group_id": "DGXXXXXXXXXXXXXXXXXA",
		"name": "Group A",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	},
	{
		"desc": "This is group B",
		"group_id": "DGXXXXXXXXXXXXXXXXXB",
		"name": "Group B",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	}],
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetUserGroups(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getGroupsResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserGroups("DU3RP9I2WOC59VZX672N")
	if err != nil {
		t.Errorf("Unexpected error from GetUserGroups call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 2 {
		t.Errorf("Expected 2 groups, but got %d", len(result.Response))
	}
	if result.Response[0].GroupID != "DGXXXXXXXXXXXXXXXXXA" {
		t.Errorf("Expected group ID DGXXXXXXXXXXXXXXXXXA, but got %s", result.Response[0].GroupID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getGroupsPage1Response = `{
	"response": [{
		"desc": "This is group A",
		"group_id": "DGXXXXXXXXXXXXXXXXXA",
		"name": "Group A",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	},
	{
		"desc": "This is group B",
		"group_id": "DGXXXXXXXXXXXXXXXXXB",
		"name": "Group B",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	}],
	"stat": "OK",
	"metadata": {
		"prev_offset": null,
		"next_offset": 2,
		"total_objects": 4
	}
}`

const getGroupsPage2Response = `{
	"response": [{
		"desc": "This is group C",
		"group_id": "DGXXXXXXXXXXXXXXXXXC",
		"name": "Group C",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	},
	{
		"desc": "This is group D",
		"group_id": "DGXXXXXXXXXXXXXXXXXD",
		"name": "Group D",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	}],
	"stat": "OK",
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 4
	}
}`

func TestGetUserGroupsMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getGroupsPage1Response)
			} else {
				fmt.Fprintln(w, getGroupsPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserGroups("DU3RP9I2WOC59VZX672N")

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 4 {
		t.Errorf("Expected four groups in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetUserGroupsPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetUserGroups("DU3RP9I2WOC59VZX672N", func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUserPhonesResponse = `{
	"stat": "OK",
	"response": [{
		"activated": false,
		"capabilities": [
			"sms",
			"phone",
			"push"
		],
		"extension": "",
		"name": "",
		"number": "+15035550102",
		"phone_id": "DPFZRS9FB0D46QFTM890",
		"platform": "Apple iOS",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Mobile"
	},
	{
		"activated": false,
		"capabilities": [
			"phone"
		],
		"extension": "",
		"name": "",
		"number": "+15035550103",
		"phone_id": "DPFZRS9FB0D46QFTM891",
		"platform": "Unknown",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Landline"
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetUserPhones(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getUserPhonesResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserPhones("DU3RP9I2WOC59VZX672N")
	if err != nil {
		t.Errorf("Unexpected error from GetUserPhones call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 2 {
		t.Errorf("Expected 2 phones, but got %d", len(result.Response))
	}
	if result.Response[0].PhoneID != "DPFZRS9FB0D46QFTM890" {
		t.Errorf("Expected phone ID DPFZRS9FB0D46QFTM890, but got %s", result.Response[0].PhoneID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUserPhonesPage1Response = `{
	"stat": "OK",
	"response": [{
		"activated": false,
		"capabilities": [
			"sms",
			"phone",
			"push"
		],
		"extension": "",
		"name": "",
		"number": "+15035550102",
		"phone_id": "DPFZRS9FB0D46QFTM890",
		"platform": "Apple iOS",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Mobile"
	},
	{
		"activated": false,
		"capabilities": [
			"phone"
		],
		"extension": "",
		"name": "",
		"number": "+15035550103",
		"phone_id": "DPFZRS9FB0D46QFTM891",
		"platform": "Unknown",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Landline"
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 2,
		"total_objects": 4
	}
}`

const getUserPhonesPage2Response = `{
	"stat": "OK",
	"response": [{
		"activated": false,
		"capabilities": [
			"sms",
			"phone",
			"push"
		],
		"extension": "",
		"name": "",
		"number": "+15035550102",
		"phone_id": "DPFZRS9FB0D46QFTM890",
		"platform": "Apple iOS",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Mobile"
	},
	{
		"activated": false,
		"capabilities": [
			"phone"
		],
		"extension": "",
		"name": "",
		"number": "+15035550103",
		"phone_id": "DPFZRS9FB0D46QFTM891",
		"platform": "Unknown",
		"postdelay": null,
		"predelay": null,
		"sms_passcodes_sent": false,
		"type": "Landline"
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 4
	}
}`

func TestGetUserPhonesMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getUserPhonesPage1Response)
			} else {
				fmt.Fprintln(w, getUserPhonesPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserPhones("DU3RP9I2WOC59VZX672N")

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 4 {
		t.Errorf("Expected four phones in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetUserPhonesPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetUserPhones("DU3RP9I2WOC59VZX672N", func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUserTokensResponse = `{
	"stat": "OK",
	"response": [{
		"type": "d1",
		"serial": "0",
		"token_id": "DHEKH0JJIYC1LX3AZWO4"
	},
	{
		"type": "d1",
		"serial": "7",
		"token_id": "DHUNT3ZVS3ACF8AEV2WG",
		"totp_step": null
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetUserTokens(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getUserTokensResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserTokens("DU3RP9I2WOC59VZX672N")
	if err != nil {
		t.Errorf("Unexpected error from GetUserTokens call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 2 {
		t.Errorf("Expected 2 tokens, but got %d", len(result.Response))
	}
	if result.Response[0].TokenID != "DHEKH0JJIYC1LX3AZWO4" {
		t.Errorf("Expected token ID DHEKH0JJIYC1LX3AZWO4, but got %s", result.Response[0].TokenID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUserTokensPage1Response = `{
	"stat": "OK",
	"response": [{
		"type": "d1",
		"serial": "0",
		"token_id": "DHEKH0JJIYC1LX3AZWO4"
	},
	{
		"type": "d1",
		"serial": "7",
		"token_id": "DHUNT3ZVS3ACF8AEV2WG",
		"totp_step": null
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 2,
		"total_objects": 4
	}
}`

const getUserTokensPage2Response = `{
	"stat": "OK",
	"response": [{
		"type": "d1",
		"serial": "0",
		"token_id": "DHEKH0JJIYC1LX3AZWO4"
	},
	{
		"type": "d1",
		"serial": "7",
		"token_id": "DHUNT3ZVS3ACF8AEV2WG",
		"totp_step": null
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 4
	}
}`

func TestGetUserTokensMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getUserTokensPage1Response)
			} else {
				fmt.Fprintln(w, getUserTokensPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserTokens("DU3RP9I2WOC59VZX672N")

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 4 {
		t.Errorf("Expected four tokens in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetUserTokensPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetUserTokens("DU3RP9I2WOC59VZX672N", func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const associateUserTokenResponse = `{
	"stat": "OK",
	"response": ""
}`

func TestAssociateUserToken(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, associateUserTokenResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.AssociateUserToken("DU3RP9I2WOC59VZX672N", "DHEKH0JJIYC1LX3AZWO4")
	if err != nil {
		t.Errorf("Unexpected error from AssociateUserToken call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 0 {
		t.Errorf("Expected empty response, but got %s", result.Response)
	}
}

const getUserU2FTokensResponse = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV"
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 1
	}
}`

func TestGetUserU2FTokens(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getUserU2FTokensResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserU2FTokens("DU3RP9I2WOC59VZX672N")
	if err != nil {
		t.Errorf("Unexpected error from GetUserU2FTokens call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 token, but got %d", len(result.Response))
	}
	if result.Response[0].RegistrationID != "D21RU6X1B1DF5P54B6PV" {
		t.Errorf("Expected registration ID D21RU6X1B1DF5P54B6PV, but got %s", result.Response[0].RegistrationID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getUserU2FTokensPage1Response = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV"
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 1,
		"total_objects": 2
	}
}`

const getUserU2FTokensPage2Response = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV"
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetUserU2FTokensMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getUserU2FTokensPage1Response)
			} else {
				fmt.Fprintln(w, getUserU2FTokensPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetUserU2FTokens("DU3RP9I2WOC59VZX672N")

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 2 {
		t.Errorf("Expected two tokens in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetUserU2FTokensPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetUserU2FTokens("DU3RP9I2WOC59VZX672N", func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

func TestGetGroups(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getGroupsResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetGroups()
	if err != nil {
		t.Errorf("Unexpected error from GetGroups call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 2 {
		t.Errorf("Expected 2 groups, but got %d", len(result.Response))
	}
	if result.Response[0].Name != "Group A" {
		t.Errorf("Expected group name Group A, but got %s", result.Response[0].Name)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

func TestGetGroupsMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getGroupsPage1Response)
			} else {
				fmt.Fprintln(w, getGroupsPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetGroups()

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 4 {
		t.Errorf("Expected four groups in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetGroupsPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetGroups(func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getGroupResponse = `{
	"response": {
		"desc": "Group description",
		"group_id": "DGXXXXXXXXXXXXXXXXXX",
		"name": "Group Name",
		"push_enabled": true,
		"sms_enabled": true,
		"status": "active",
		"voice_enabled": true,
		"mobile_otp_enabled": true
	},
	"stat": "OK"
}`

func TestGetGroup(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getGroupResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetGroup("DGXXXXXXXXXXXXXXXXXX")
	if err != nil {
		t.Errorf("Unexpected error from GetGroups call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if result.Response.GroupID != "DGXXXXXXXXXXXXXXXXXX" {
		t.Errorf("Expected group ID DGXXXXXXXXXXXXXXXXXX, but got %s", result.Response.GroupID)
	}
	if !result.Response.PushEnabled {
		t.Errorf("Expected push to be enabled, but got %v", result.Response.PushEnabled)
	}
}

const getPhonesResponse = `{
	"stat": "OK",
	"response": [{
		"activated": true,
		"capabilities": [
			"push",
			"sms",
			"phone",
			"mobile_otp"
		],
		"encrypted": "Encrypted",
		"extension": "",
		"fingerprint": "Configured",
		"name": "",
		"number": "+15555550100",
		"phone_id": "DPFZRS9FB0D46QFTM899",
		"platform": "Google Android",
		"postdelay": "",
		"predelay": "",
		"screenlock": "Locked",
		"sms_passcodes_sent": false,
		"tampered": "Not tampered",
		"type": "Mobile",
		"users": [{
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_login": 1474399627,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith"
		}]
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 1
	}
}`

func TestGetPhones(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getPhonesResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetPhones()
	if err != nil {
		t.Errorf("Unexpected error from GetPhones call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 phone, but got %d", len(result.Response))
	}
	if result.Response[0].PhoneID != "DPFZRS9FB0D46QFTM899" {
		t.Errorf("Expected phone ID DPFZRS9FB0D46QFTM899, but got %s", result.Response[0].PhoneID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getPhonesPage1Response = `{
	"stat": "OK",
	"response": [{
		"activated": true,
		"capabilities": [
			"push",
			"sms",
			"phone",
			"mobile_otp"
		],
		"encrypted": "Encrypted",
		"extension": "",
		"fingerprint": "Configured",
		"name": "",
		"number": "+15555550100",
		"phone_id": "DPFZRS9FB0D46QFTM899",
		"platform": "Google Android",
		"postdelay": "",
		"predelay": "",
		"screenlock": "Locked",
		"sms_passcodes_sent": false,
		"tampered": "Not tampered",
		"type": "Mobile",
		"users": [{
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_login": 1474399627,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith"
		}]
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 1,
		"total_objects": 2
	}
}`

const getPhonesPage2Response = `{
	"stat": "OK",
	"response": [{
		"activated": true,
		"capabilities": [
			"push",
			"sms",
			"phone",
			"mobile_otp"
		],
		"encrypted": "Encrypted",
		"extension": "",
		"fingerprint": "Configured",
		"name": "",
		"number": "+15555550100",
		"phone_id": "DPFZRS9FB0D46QFTM899",
		"platform": "Google Android",
		"postdelay": "",
		"predelay": "",
		"screenlock": "Locked",
		"sms_passcodes_sent": false,
		"tampered": "Not tampered",
		"type": "Mobile",
		"users": [{
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_login": 1474399627,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith"
		}]
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetPhonesMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getPhonesPage1Response)
			} else {
				fmt.Fprintln(w, getPhonesPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetPhones()

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 2 {
		t.Errorf("Expected two phones in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetPhonesPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetPhones(func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getPhoneResponse = `{
	"stat": "OK",
	"response": {
		"phone_id": "DPFZRS9FB0D46QFTM899",
		"number": "+15555550100",
		"name": "",
		"extension": "",
		"postdelay": null,
		"predelay": null,
		"type": "Mobile",
		"capabilities": [
			"sms",
			"phone",
			"push"
		],
		"platform": "Apple iOS",
		"activated": false,
		"sms_passcodes_sent": false,
		"users": [{
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith",
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"realname": "Joe Smith",
			"email": "jsmith@example.com",
			"status": "active",
			"last_login": 1343921403,
			"notes": ""
		}]
	}
}`

func TestGetPhone(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getPhoneResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetPhone("DPFZRS9FB0D46QFTM899")
	if err != nil {
		t.Errorf("Unexpected error from GetPhone call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if result.Response.PhoneID != "DPFZRS9FB0D46QFTM899" {
		t.Errorf("Expected phone ID DPFZRS9FB0D46QFTM899, but got %s", result.Response.PhoneID)
	}
}

const getTokensResponse = `{
	"stat": "OK",
	"response": [{
		"serial": "0",
		"token_id": "DHIZ34ALBA2445ND4AI2",
		"type": "d1",
		"totp_step": null,
		"users": [{
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith",
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"realname": "Joe Smith",
			"email": "jsmith@example.com",
			"status": "active",
			"last_login": 1343921403,
			"notes": ""
		}]
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 1
	}
}`

func TestGetTokens(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getTokensResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetTokens()
	if err != nil {
		t.Errorf("Unexpected error from GetTokens call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 token, but got %d", len(result.Response))
	}
	if result.Response[0].TokenID != "DHIZ34ALBA2445ND4AI2" {
		t.Errorf("Expected token ID DHIZ34ALBA2445ND4AI2, but got %s", result.Response[0].TokenID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getTokensPage1Response = `{
	"stat": "OK",
	"response": [{
		"serial": "0",
		"token_id": "DHIZ34ALBA2445ND4AI2",
		"type": "d1",
		"totp_step": null,
		"users": [{
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith",
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"realname": "Joe Smith",
			"email": "jsmith@example.com",
			"status": "active",
			"last_login": 1343921403,
			"notes": ""
		}]
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 1,
		"total_objects": 2
	}
}`

const getTokensPage2Response = `{
	"stat": "OK",
	"response": [{
		"serial": "0",
		"token_id": "DHIZ34ALBA2445ND4AI2",
		"type": "d1",
		"totp_step": null,
		"users": [{
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith",
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"realname": "Joe Smith",
			"email": "jsmith@example.com",
			"status": "active",
			"last_login": 1343921403,
			"notes": ""
		}]
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetTokensMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getTokensPage1Response)
			} else {
				fmt.Fprintln(w, getTokensPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetTokens()

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 2 {
		t.Errorf("Expected two tokens in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetTokensPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetTokens(func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getTokenResponse = `{
	"stat": "OK",
	"response": {
		"serial": "0",
		"token_id": "DHIZ34ALBA2445ND4AI2",
		"type": "d1",
		"totp_step": null,
		"users": [{
			"user_id": "DUJZ2U4L80HT45MQ4EOQ",
			"username": "jsmith",
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"realname": "Joe Smith",
			"email": "jsmith@example.com",
			"status": "active",
			"last_login": 1343921403,
			"notes": ""
		}]
	}
}`

func TestGetToken(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getTokenResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetToken("DPFZRS9FB0D46QFTM899")
	if err != nil {
		t.Errorf("Unexpected error from GetToken call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if result.Response.TokenID != "DHIZ34ALBA2445ND4AI2" {
		t.Errorf("Expected token ID DHIZ34ALBA2445ND4AI2, but got %s", result.Response.TokenID)
	}
}

const getU2FTokensResponse = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV",
		"user": {
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"created": 1384275337,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_directory_sync": 1384275337,
			"last_login": 1514922986,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DU3RP9I2WOC59VZX672N",
			"username": "jsmith"
		}
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": null,
		"total_objects": 1
	}
}`

func TestGetU2FTokens(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getU2FTokensResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetU2FTokens()
	if err != nil {
		t.Errorf("Unexpected error from GetU2FTokens call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 token, but got %d", len(result.Response))
	}
	if result.Response[0].RegistrationID != "D21RU6X1B1DF5P54B6PV" {
		t.Errorf("Expected registration ID D21RU6X1B1DF5P54B6PV, but got %s", result.Response[0].RegistrationID)
	}

	request_query := last_request.URL.Query()
	if request_query["limit"][0] != "100" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "0" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

const getU2FTokensPage1Response = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV",
		"user": {
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"created": 1384275337,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_directory_sync": 1384275337,
			"last_login": 1514922986,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DU3RP9I2WOC59VZX672N",
			"username": "jsmith"
		}
	}],
	"metadata": {
		"prev_offset": null,
		"next_offset": 1,
		"total_objects": 2
	}
}`

const getU2FTokensPage2Response = `{
	"stat": "OK",
	"response": [{
		"date_added": 1444678994,
		"registration_id": "D21RU6X1B1DF5P54B6PV",
		"user": {
			"alias1": "joe.smith",
			"alias2": "jsmith@example.com",
			"alias3": null,
			"alias4": null,
			"created": 1384275337,
			"email": "jsmith@example.com",
			"firstname": "Joe",
			"last_directory_sync": 1384275337,
			"last_login": 1514922986,
			"lastname": "Smith",
			"notes": "",
			"realname": "Joe Smith",
			"status": "active",
			"user_id": "DU3RP9I2WOC59VZX672N",
			"username": "jsmith"
		}
	}],
	"metadata": {
		"prev_offset": 0,
		"next_offset": null,
		"total_objects": 2
	}
}`

func TestGetU2fTokensMultiple(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(requests) == 0 {
				fmt.Fprintln(w, getU2FTokensPage1Response)
			} else {
				fmt.Fprintln(w, getU2FTokensPage2Response)
			}
			requests = append(requests, r)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetU2FTokens()

	if len(requests) != 2 {
		t.Errorf("Expected two requets, found %d", len(requests))
	}

	if len(result.Response) != 2 {
		t.Errorf("Expected two tokens in the response, found %d", len(result.Response))
	}

	if err != nil {
		t.Errorf("Expected err to be nil, found %s", err)
	}
}

func TestGetU2FTokensPageArgs(t *testing.T) {
	requests := []*http.Request{}
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getEmptyPageArgsResponse)
			requests = append(requests, r)
		}),
	)

	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	_, err := duo.GetU2FTokens(func(values *url.Values) {
		values.Set("limit", "200")
		values.Set("offset", "1")
		return
	})

	if err != nil {
		t.Errorf("Encountered unexpected error: %s", err)
	}

	if len(requests) != 1 {
		t.Errorf("Expected there to be one request, found %d", len(requests))
	}
	request := requests[0]
	request_query := request.URL.Query()
	if request_query["limit"][0] != "200" {
		t.Errorf("Expected to see a limit of 100 in request, bug got %s", request_query["limit"])
	}
	if request_query["offset"][0] != "1" {
		t.Errorf("Expected to see an offset of 0 in request, bug got %s", request_query["offset"])
	}
}

func TestGetU2FToken(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getU2FTokensResponse)
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	result, err := duo.GetU2FToken("D21RU6X1B1DF5P54B6PV")
	if err != nil {
		t.Errorf("Unexpected error from GetU2FToken call %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if len(result.Response) != 1 {
		t.Errorf("Expected 1 token, but got %d", len(result.Response))
	}
	if result.Response[0].RegistrationID != "D21RU6X1B1DF5P54B6PV" {
		t.Errorf("Expected registration ID D21RU6X1B1DF5P54B6PV, but got %s", result.Response[0].RegistrationID)
	}
}

/*
 * Logs
 */

// TestLogListV2Metadata ensures that the next offset is properly configured based on known metadata.
func TestLogListV2Metadata(t *testing.T) {
	page := LogListV2Metadata{
		NextOffset: []string{"1532951895000", "af0ba235-0b33-23c8-bc23-a31aa0231de8"},
	}
	nextOffset := page.GetNextOffset()
	if nextOffset == nil {
		t.Fatalf("Expected option to configure next offset, got nil")
	}
	params := url.Values{}
	nextOffset(&params)
	if nextParam := params.Get("next_offset"); nextParam != "1532951895000,af0ba235-0b33-23c8-bc23-a31aa0231de8" {
		t.Fatalf("Expected option to configure next offset to be '1532951895000,af0ba235-0b33-23c8-bc23-a31aa0231de8', got %q", nextParam)
	}

	lastPage := LogListV2Metadata{}
	if lastPage.GetNextOffset() != nil {
		t.Errorf("Expected nil option to represent no more available logs, got a non-nil option")
	}
}

// TestLogListV2Metadata ensures that timestamps are properly parsed from log data.
func TestParseLogV1Timestamp(t *testing.T) {
	validLog := map[string]interface{}{
		"timestamp": 1346172820,
	}
	timestamp, err := parseLogV1Timestamp(validLog)
	if err != nil {
		t.Errorf("Failed to parse log timestamp: %v", err)
	}
	if expectedTs := time.Unix(1346172820, 0); !timestamp.Equal(expectedTs) {
		t.Errorf("Parsed incorrect value for log timestamp, expected %v but got: %v", expectedTs, timestamp)
	}
}

// getAuthLogsResponse is an example response from the Duo API documentation example: https://duo.com/docs/adminapi#authentication-logs
const getAuthLogsResponse = `{
    "response": {
        "authlogs": [
            {
                "access_device": {
                    "browser": "Chrome",
                    "browser_version": "67.0.3396.99",
                    "flash_version": "uninstalled",
                    "hostname": "null",
                    "ip": "169.232.89.219",
                    "java_version": "uninstalled",
                    "location": {
                        "city": "Ann Arbor",
                        "country": "United States",
                        "state": "Michigan"
                    },
                    "os": "Mac OS X",
                    "os_version": "10.14.1"
                },
                "application": {
                    "key": "DIY231J8BR23QK4UKBY8",
                    "name": "Microsoft Azure Active Directory"
                },
                "auth_device": {
                    "ip": "192.168.225.254",
                    "location": {
                        "city": "Ann Arbor",
                        "country": "United States",
                        "state": "Michigan"
                    },
                    "name": "My iPhone X (734-555-2342)"
                },
                "event_type": "authentication",
                "factor": "duo_push",
                "reason": "user_approved",
                "result": "success",
                "timestamp": 1532951962,
                "trusted_endpoint_status": "not trusted",
                "txid": "340a23e3-23f3-23c1-87dc-1491a23dfdbb",
                "user": {
                    "key": "DU3KC77WJ06Y5HIV7XKQ",
                    "name": "narroway@example.com"
                }
            }
        ],
        "metadata": {
            "next_offset": [
                "1532951895000",
                "af0ba235-0b33-23c8-bc23-a31aa0231de8"
            ],
            "total_objects": 1
        }
    },
    "stat": "OK"
}`

// TestGetAuthLogs ensures proper functionality of the client.GetAuthLogs method.
func TestGetAuthLogs(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getAuthLogsResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	lastMetadata := LogListV2Metadata{
		NextOffset: []string{"1532951920000", "b40ba235-0b33-23c8-bc23-a31aa0231db4"},
	}
	mintime := time.Unix(1532951960, 0)
	window := 5 * time.Second
	result, err := duo.GetAuthLogs(mintime, window, lastMetadata.GetNextOffset())

	if err != nil {
		t.Errorf("Unexpected error from GetAuthLogs call: %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if length := len(result.Response.Logs); length != 1 {
		t.Errorf("Expected 1 log, but got %d", length)
	}
	if txid := result.Response.Logs[0]["txid"]; txid != "340a23e3-23f3-23c1-87dc-1491a23dfdbb" {
		t.Errorf("Expected txid '340a23e3-23f3-23c1-87dc-1491a23dfdbb', but got %v", txid)
	}
	if next := result.Response.Metadata.GetNextOffset(); next == nil {
		t.Errorf("Expected metadata.GetNextOffset option to configure pagination for next request, got nil")
	}

	request_query := last_request.URL.Query()
	if qMintime := request_query["mintime"][0]; qMintime != "1532951960000" {
		t.Errorf("Expected to see a mintime of 153295196000 in request, but got %q", qMintime)
	}
	if qMaxtime := request_query["maxtime"][0]; qMaxtime != "1532951965000" {
		t.Errorf("Expected to see a maxtime of 153295196500 in request, but got %q", qMaxtime)
	}
	if qNextOffset := request_query["next_offset"][0]; qNextOffset != "1532951920000,b40ba235-0b33-23c8-bc23-a31aa0231db4" {
		t.Errorf("Expected to see a next_offset of 1532951920000,b40ba235-0b33-23c8-bc23-a31aa0231db4 in request, but got %q", qNextOffset)
	}
}

// getAdminLogsResponse is an example response from the Duo API documentation example: https://duo.com/docs/adminapi#administrator-logs
const getAdminLogsResponse = `{
	"stat": "OK",
	"response": [{
		"action": "user_update",
		"description": "{\"notes\": \"Joe asked for their nickname to be displayed instead of Joseph.\", \"realname\": \"Joe Smith\"}",
		"object": "jsmith",
		"timestamp": 1346172820,
		"username": "admin"
	},
	{
		"action": "admin_login_error",
		"description": "{\"ip_address\": \"10.1.23.116\", \"error\": \"SAML login is disabled\", \"email\": \"narroway@example.com\"}",
		"object": null,
		"timestamp": 1446172820,
		"username": ""
	}]
  }`

// TestGetAdminLogs ensures proper functionality of the client.GetAdminLogs method.
func TestGetAdminLogs(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getAdminLogsResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	mintime := time.Unix(1346172815, 0)
	maxtime := mintime.Add(time.Second * 10)
	result, err := duo.GetAdminLogs(mintime)

	if err != nil {
		t.Errorf("Unexpected error from GetAdminLogs call: %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if length := len(result.Logs); length != 2 {
		t.Errorf("Expected 2 logs, but got %d", length)
	}
	timestamp, err := result.Logs[0].Timestamp()
	if err != nil {
		t.Errorf("Failed to parse timestamp timestamp: %v", err)
	}
	if expectedTs := time.Unix(1346172820, 0); !expectedTs.Equal(timestamp) {
		t.Errorf("Expected timestamp %v, but got: %v", expectedTs, timestamp)
	}
	if next := result.Logs.GetNextOffset(maxtime); next != nil {
		t.Errorf("Expected no next page available, got non-nil option")
	}

	request_query := last_request.URL.Query()
	if qMintime := request_query["mintime"][0]; qMintime != "1346172815" {
		t.Errorf("Expected to see a mintime of 1346172815 in request, but got %q", qMintime)
	}
}

// TestAdminLogsNextOffset ensures proper pagination functionality for AdminLogResult
func TestAdminLogsNextOffset(t *testing.T) {
	maxtime := time.Unix(1346172825, 0)

	// Ensure < 1000 logs returns none
	result := &AdminLogResult{}
	if next := result.Logs.GetNextOffset(maxtime); next != nil {
		t.Errorf("Expected no next page available, got non-nil option")
	}

	// Ensure mintime == maxtime returns maxtime + 1
	logs := make([]AdminLog, 0, 1000)
	for i := 0; i < 1000; i++ {
		logs = append(logs, AdminLog{"timestamp": 1346172816})
	}
	result.Logs = AdminLogList(logs)
	params := &url.Values{}
	result.Logs.GetNextOffset(maxtime)(params)
	if newMintime := params.Get("mintime"); newMintime != "1346172817" {
		t.Errorf("Expected new mintime to be 1346172817, got: %v", newMintime)
	}

	// Ensure single maxtime returns maxtime
	result.Logs[0] = AdminLog{"timestamp": 1346172820}
	params = &url.Values{}
	result.Logs.GetNextOffset(maxtime)(params)
	if newMintime := params.Get("mintime"); newMintime != "1346172820" {
		t.Errorf("Expected new mintime to be 1346172820, got: %v", newMintime)
	}
}

// getTelephonyLogsResponse is an example response from the Duo API documentation example: https://duo.com/docs/adminapi#telephony-logs
const getTelephonyLogsResponse = `{
	"stat": "OK",
	"response": [{
		"context": "authentication",
		"credits": 1,
		"phone": "+15035550100",
		"timestamp": 1346172697,
		"type": "sms"
	}]
  }`

// TestGetTelephonyLogs ensures proper functionality of the client.GetTelephonyLogs method.
func TestGetTelephonyLogs(t *testing.T) {
	var last_request *http.Request
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, getTelephonyLogsResponse)
			last_request = r
		}),
	)
	defer ts.Close()

	duo := buildAdminClient(ts.URL, nil)

	mintime := time.Unix(1346172600, 0)
	maxtime := mintime.Add(time.Second * 10)
	result, err := duo.GetTelephonyLogs(mintime)

	if err != nil {
		t.Errorf("Unexpected error from GetTelephonyLogs call: %v", err.Error())
	}
	if result.Stat != "OK" {
		t.Errorf("Expected OK, but got %s", result.Stat)
	}
	if length := len(result.Logs); length != 1 {
		t.Errorf("Expected 1 logs, but got %d", length)
	}
	timestamp, err := result.Logs[0].Timestamp()
	if err != nil {
		t.Errorf("Failed to parse timestamp timestamp: %v", err)
	}
	if expectedTs := time.Unix(1346172697, 0); !expectedTs.Equal(timestamp) {
		t.Errorf("Expected timestamp %v, but got: %v", expectedTs, timestamp)
	}
	if next := result.Logs.GetNextOffset(maxtime); next != nil {
		t.Errorf("Expected no next page available, got non-nil option")
	}

	request_query := last_request.URL.Query()
	if qMintime := request_query["mintime"][0]; qMintime != "1346172600" {
		t.Errorf("Expected to see a mintime of 1346172600 in request, but got %q", qMintime)
	}
}

// TestTelephonyLogsNextOffset ensures proper pagination functionality for TelephonyLogResult
func TestTelephonyLogsNextOffset(t *testing.T) {
	maxtime := time.Unix(1346172825, 0)

	// Ensure < 1000 logs returns none
	result := &TelephonyLogResult{}
	if next := result.Logs.GetNextOffset(maxtime); next != nil {
		t.Errorf("Expected no next page available, got non-nil option")
	}

	// Ensure mintime == maxtime returns maxtime + 1
	logs := make([]TelephonyLog, 0, 1000)
	for i := 0; i < 1000; i++ {
		logs = append(logs, TelephonyLog{"timestamp": 1346172816})
	}
	result.Logs = TelephonyLogList(logs)
	params := &url.Values{}
	result.Logs.GetNextOffset(maxtime)(params)
	if newMintime := params.Get("mintime"); newMintime != "1346172817" {
		t.Errorf("Expected new mintime to be 1346172817, got: %v", newMintime)
	}

	// Ensure single maxtime returns maxtime
	result.Logs[0] = TelephonyLog{"timestamp": 1346172820}
	params = &url.Values{}
	result.Logs.GetNextOffset(maxtime)(params)
	if newMintime := params.Get("mintime"); newMintime != "1346172820" {
		t.Errorf("Expected new mintime to be 1346172820, got: %v", newMintime)
	}
}
