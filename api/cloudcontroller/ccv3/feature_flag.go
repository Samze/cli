package ccv3

import (
	"encoding/json"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
)

// FeatureFlag represents a Cloud Controller V3 Feature Flag.
type FeatureFlag struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (f FeatureFlag) MarshalJSON() ([]byte, error) {
	var ccBodyFlag struct {
		Enabled bool `json:"enabled"`
	}

	ccBodyFlag.Enabled = f.Enabled

	return json.Marshal(ccBodyFlag)
}

func (client *Client) GetFeatureFlag(flagName string) (FeatureFlag, Warnings, error) {
	var responseBody FeatureFlag

	_, warnings, err := client.makeRequest(requestParams{
		RequestName:  internal.GetFeatureFlagRequest,
		URIParams:    internal.Params{"name": flagName},
		ResponseBody: &responseBody,
	})

	return responseBody, warnings, err
}

// GetFeatureFlags lists feature flags.
func (client *Client) GetFeatureFlags() ([]FeatureFlag, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.GetFeatureFlagsRequest,
	})

	if err != nil {
		return nil, nil, err
	}

	var fullFeatureFlagList []FeatureFlag
	warnings, err := client.paginate(request, FeatureFlag{}, func(item interface{}) error {
		if featureFlag, ok := item.(FeatureFlag); ok {
			fullFeatureFlagList = append(fullFeatureFlagList, featureFlag)
		} else {
			return ccerror.UnknownObjectInListError{
				Expected:   FeatureFlag{},
				Unexpected: item,
			}
		}
		return nil
	})

	return fullFeatureFlagList, warnings, err

}

func (client *Client) UpdateFeatureFlag(flag FeatureFlag) (FeatureFlag, Warnings, error) {
	var responseBody FeatureFlag

	_, warnings, err := client.makeRequest(requestParams{
		RequestName:  internal.PatchFeatureFlagRequest,
		URIParams:    internal.Params{"name": flag.Name},
		RequestBody:  flag,
		ResponseBody: &responseBody,
	})

	return responseBody, warnings, err
}
