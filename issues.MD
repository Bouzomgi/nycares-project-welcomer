Critical (Blocks Deployment)

1. CDK stack is missing most AWS resources and Lambda configuration  
   infra/stack/projectnotifierstack.go creates lambdas but:

- No environment variables set on any lambda (they all need NYCARES\_\* env vars)
- No DynamoDB table, S3 bucket, or SNS topic created
- No IAM permissions for ComputeMessageToSend (DynamoDB read), SendAndPinMessage (S3 read),
  RecordMessage (DynamoDB write), NotifyCompletion/DLQNotifier (SNS publish)

2. Workflow JSON uses placeholder ARNs that won't resolve
   infra/DailyProjectNotificationWorkflow.json references
   arn:aws:lambda:REGION:ACCOUNT_ID:function:LoginLambda — these are never substituted with real
   values, and the naming doesn't match CDK output (CDK creates login kebab-case, workflow
   references LoginLambda).

3. DynamoDB call ignores passed context
   internal/platform/dynamo/service/storednotificationservice.go uses context.Background() instead
   of the passed ctx, breaking timeout/cancellation propagation from Lambda and Step Functions.

---

High (Runtime Failures)

4. SNS publish errors silently swallowed
   internal/app/requestapproval/usecase.go:32 and internal/app/notifycompletion/usecase.go call
   PublishNotification without checking the returned error.

5. No HTTP client timeout
   internal/platform/http/service/httpservice.go creates http.Client without a Timeout field. HTTP
   calls to NYC Cares API could hang forever.

6. HTML injection in approval callback
   lambda/approvalcallback/handler.go:43 interpolates err.Error() directly into HTML without
   escaping — an attacker could craft a token that triggers XSS.

7. No approval callback authentication
   The GET /callback endpoint accepts any request with no validation. Anyone with the URL can
   approve/reject messages.

---

Medium (Bugs / Code Quality)

8. Dead code

- internal/config/toremove.go — file with a struct called Nope
- internal/platform/http/service/messageservice.go — GetFirstMessageId defined but never called
- TODO comment at messageservice.go:72: // TODO fix all these endpoints

9. Date type inconsistency across layers
   domain.Project.Date is time.Time, but domain.ProjectNotification.Date is string. Models use
   string dates. This causes repeated parsing and potential timezone bugs.

10. Hardcoded 10-second timeout in every handler
    All 9 handlers use context.WithTimeout(ctx, 10\*time.Second). Some operations (login, HTTP
    calls) may legitimately need more time. This should be configurable or at least differentiated.

11. No logging anywhere
    Zero structured logging in the entire codebase. When something fails in production, the only
    signal is a panic or a Step Functions error state.

12. No tests
    internal/service/service_test.go exists but contains only TODO comments. No unit tests for
    domain logic, use cases, model conversions, or config loading.

---

Low (Technical Debt)

13. Inconsistent handler signatures — Login handler takes no input,
    NotifyCompletion/DLQNotifier return only error (no output), DLQNotifier takes
    map[string]interface{} instead of a typed model.

14. Magic numbers — 10 seconds timeout, 86400 seconds workflow timeout, -7 and -2 day offsets
    for welcome/reminder logic are all hardcoded.

15. Cookie jar thread safety — cookiejar.CookieJar isn't safe for concurrent use, though in
    practice each Lambda invocation gets its own instance so this is low risk.

16. No retry logic — AWS SDK calls (SNS, DynamoDB, S3) have no application-level retry beyond
    SDK defaults.

---

The architecture and layered design are solid. The biggest gap is the CDK stack — it creates
lambdas but doesn't provision the resources they depend on or wire them with env vars and
permissions. The workflow JSON placeholder ARNs are a related blocker. Want me to tackle any of
these?
