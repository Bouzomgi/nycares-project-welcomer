package s3service

// Covered by e2e tests.
// GetMessageContent has no pure-logic paths that skip the S3 client call:
// the s3:// URI prefix stripping is a silent transformation with no
// validation error, so every code path requires a real or mocked client.
// Integration coverage lives in the end-to-end test suite.
