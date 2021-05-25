// VAS is a simple implimentation of the Chrome Verified Access API, which
// does not support the certificate flow. Pay special attention to
// https://support.google.com/chrome/a/answer/7156268
// when setting this up. There are several things you have to configure in
// user settings, device settings, and extension settings.
package vas

import (
	"encoding/base64"
	"net/http"

	"google.golang.org/api/verifiedaccess/v1"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=. vas.proto

// VAS represents authentication information necessary to make calls to the
// Chrome Verified Access Service.
type VAS struct {
	cs *verifiedaccess.ChallengeService
}

// NewFromChallengeService returns a VAS from a
// verifiedaccess.ChallengeService you construct yourself. This allows you
// to use an API key if you want, though note that it appears API key
// authenticated clients can only create and not verify requests.
func NewFromChallengeService(cs *verifiedaccess.ChallengeService) VAS {
	return VAS{cs}
}

// NewFromOAuthClient is probably the easiest way to create a VAS. Get an
// OAuth client by calling `.Client()` on an `oauth2.Config` struct. You can
// find many examples of how to do this online, as this is a common flow for
// using any Google API.
func NewFromOAuthClient(c *http.Client) (VAS, error) {
	service, err := verifiedaccess.New(c)
	if err != nil {
		return VAS{}, err
	}

	cs := verifiedaccess.NewChallengeService(service)
	return NewFromChallengeService(cs), nil
}

// GetChallenge returns a base64 encoded string which should be
// de-serialized in to an ArrayBuffer (see util.js) and passed to
// chrome.enterprise.platformKeys.challengeMachineKey or
// chrome.enterprise.platformKeys.challengeUserKey
func (vas VAS) GetChallenge() (string, error) {
	// the challenge comes to us as 2 base64 strings. We need to decode
	// those, join them with protobuf, then base64 encode the newly merged
	// data

	resp, err := vas.cs.Create(new(verifiedaccess.Empty)).Do()
	if err != nil {
		return "", err
	}

	var challenge SignedData
	challenge.Data, err =
		base64.StdEncoding.DecodeString(resp.Challenge.Data)
	if err != nil {
		return "", err
	}
	challenge.Signature, err =
		base64.StdEncoding.DecodeString(resp.Challenge.Signature)
	if err != nil {
		return "", err
	}

	bytes, err := proto.Marshal(&challenge)
	if err != nil {
		return "", err
	}

	b64Challenge := base64.StdEncoding.EncodeToString(bytes)

	return b64Challenge, nil
}

// VerifyResponse verifies a serialized ArrayBuffer response from
// chrome.enterprise.platformKeys.challengeMachineKey or
// chrome.enterprise.platformKeys.challengeUserKey. If verifying a user,
// you should put the expected email address of the user in the
// `expectedIdentity` field. If verifying a device you may put the domain
// in the `expectedIdentity` field, though this is optional. If verifying a
// device, the device's device ID will be returned as a string. Failures of
// any kind (authentication failures or application issues) are returned as
// errors.
func (vas VAS) VerifyResponse(
	respFromChromebook string, expectedIdentity string,
) (
	deviceID string, err error,
) {
	// decode the SignedData, which we will later split in to data and
	// signature, then re-encode.
	bytesFromChromebook, err :=
		base64.StdEncoding.DecodeString(respFromChromebook)
	if err != nil {
		return "", err
	}

	var signedData SignedData
	proto.Unmarshal(bytesFromChromebook, &signedData)

	var verifyReq verifiedaccess.VerifyChallengeResponseRequest

	verifyReq.ChallengeResponse = new(verifiedaccess.SignedData)
	verifyReq.ChallengeResponse.Data =
		base64.StdEncoding.EncodeToString(signedData.Data)
	verifyReq.ChallengeResponse.Signature =
		base64.StdEncoding.EncodeToString(signedData.Signature)
	verifyReq.ExpectedIdentity = expectedIdentity

	resp, err := vas.cs.Verify(&verifyReq).Do()
	if err != nil {
		return "", err
	}

	return resp.DevicePermanentId, nil
}
