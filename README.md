A simple-ish server implimentation for the [Chrome Verified Access API](https://developers.google.com/chrome/verified-access/overview). See also: the [client](https://developer.chrome.com/docs/extensions/reference/enterprise_platformKeys/#method-challengeMachineKey) [API](https://developer.chrome.com/docs/extensions/reference/enterprise_platformKeys/#method-challengeUserKey) and [setup guide](https://support.google.com/chrome/a/answer/7156268).

**Prerequisites:** The server must already be authenticated via OAuth to Google and have the `https://www.googleapis.com/auth/verifiedaccess` scope. Additionally the user corresponding to the OAuth token must be configured under "Service accounts that are allowed to receive device ID" and/or "Service accounts that are allowed to receive user data" in Google Admin. **Read the [setup guide](https://support.google.com/chrome/a/answer/7156268)**.

The exchange goeth thusly:
* The server gets a challenge string by calling `GetChallenge()`.
* This string is passed to the client.
* The string is de-serialized in to an ArrayBuffer. See [util.js](util.js) for a function to do this.
* Your extension passes the array buffer to [chrome.enterprise.platformKeys.challengeMachineKey](https://developer.chrome.com/docs/extensions/reference/enterprise_platformKeys/#method-challengeMachineKey) or [chrome.enterprise.platformKeys.challengeUserKey](https://developer.chrome.com/docs/extensions/reference/enterprise_platformKeys/#method-challengeUserKey).
* An ArrayBuffer is returned to your callback function.
* This array buffer is serialized to a string. See [util.js](util.js) for a function to do this.
* This string is passed to the server.
* The server calls `VerifyResponse()`, optionally specifying the expected identity of the client.

This entire exchange must take place in less than 60 seconds.

You might have noticed that the Google-recommended flow has the client requesting a challenge from Google directly. We don't do that.

Here is an example of what the client-side code might look like:
```javascript
chrome.enterprise.platformKeys.challengeUserKey(
	deserializeArrayBuffer('p6qOaLxKAAmFup8HRDYex08i...'),
	false,
	resp => console.log(serializeArrayBuffer(resp)),
)
```
