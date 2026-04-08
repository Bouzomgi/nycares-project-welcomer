# Plan: Thank-You Message (1 Hour After Project Ends)

## Context

The workflow currently sends welcome (7 days before) and reminder (2 days before) messages. This adds a third notification type — a thank-you message sent 1 hour after the project ends, with the same human-approval + pin flow as the other types.

The NYC Cares API already returns session end time data (`Start_DateTime__c` + `Duration__c`), which is the right source for computing end datetime in UTC. These fields exist in the raw DTO response but aren't currently mapped into the domain.

## Approach

Add `ThankYou` as a first-class `NotificationType`. Thread `EndDateTime` through the data pipeline. Add a Step Functions **Choice + Wait** pattern after `ComputeMessageToSend` so thank-you executions pause until `endDateTime + 1 hour` before proceeding into the existing approval→send→pin→record flow.

The daily noon trigger is sufficient — the Wait state handles projects that end after noon (pauses until target time), and for projects that ended before noon, Step Functions passes a past timestamp and continues immediately.

---

## Files to Modify

### Domain layer

**`internal/domain/notifications.go`**
- Add `ThankYou NotificationType = 2`
- Add `HasSentThankYou bool` to `ProjectNotification`
- Add `ThankYou` case to `String()` and `ParseNotificationType()`

**`internal/domain/project.go`**
- Add `EndDateTime time.Time`

### DTO / HTTP layer

**`internal/platform/http/dto/upcoming.go`**
- Add `StartDateTimeUTC string` (maps `Start_DateTime__c`) and `DurationHours float64` (maps `Duration__c`) to `UpcomingSession`
- In `ToDomainProjects()`: compute `EndDateTime = parse(Start_DateTime__c) + Duration__c hours` and map to `domain.Project.EndDateTime`
- Use `time.Parse(time.RFC3339, s.StartDateTimeUTC)` — `Start_DateTime__c` is a proper UTC string (e.g. `"2026-04-11T14:00:00.000+0000"`)

### Models layer (Step Functions I/O)

**`internal/models/fetchprojects.go`**
- Add `EndDateTime string` (RFC3339) to the `project` struct
- Thread through `BuildDomainProject` and `buildModelProject`

**`internal/models/computemessage.go`**
- Add `TargetSendTime string` (RFC3339, empty for welcome/reminder) to `ComputeMessageOutput`
- Add `HasSentThankYou bool` to `projectNotification` model
- Update `ConvertDomainProjectNotification`, `ConvertModelProjectNotification`, `ConvertProjectNotificationToDomainProject`

### DynamoDB layer

**`internal/platform/dynamo/dto/storednotificationdto.go`**
- Add `HasSentThankYou bool \`json:"hasSentThankYou"\``

**`internal/platform/dynamo/service/storednotificationservice.go`**
- Thread `HasSentThankYou` through `GetProjectNotification` and `UpsertProjectNotification`

### Compute message use case

**`internal/app/computemessage/usecase.go`**
- Change `computeNotificationType` to accept `endDateTime time.Time` and return `(NotificationType, time.Time, error)` where the second return is `targetSendTime`
- Restructure logic:
  1. `!isTeamLeader` → `NotTeamLeader`
  2. `status == "Canceled"` → `ProjectCancelled`
  3. `now.Before(projectStartDate)`: existing welcome/reminder/skip logic (unchanged)
  4. `now.After(projectStartDate)` (project started or ended): if `!HasSentThankYou` → return `ThankYou`, `targetSendTime = endDateTime.Add(1 * time.Hour)`; if `HasSentThankYou` → `AllNotificationsSent`
- Update `Execute` to pass `project.EndDateTime` and populate `TargetSendTime` in output

> **Why restructure?** Currently `!now.Before(projectStartDate)` triggers `ProjectPassed`, which would fire on the day of the project and block the thank-you from ever being computed.

**`internal/app/computemessage/usecase_test.go`**
- Update test call signatures
- Add test cases for `ThankYou` path and updated skip conditions

### Record message use case

**`internal/app/recordmessage/usecase.go`**
- Add `updatedHasSentThankYou := existing.HasSentThankYou || (sentMessageType == domain.ThankYou)`
- Include in `updatedProjectNotification`

### Step Functions workflow

**`infra/DailyProjectNotificationWorkflow.json`**

Change `ComputeMessageToSend.Next` from `"RequestApprovalToSend"` to `"ScheduleThankYou"`. Insert two new states:

```json
"ScheduleThankYou": {
  "Type": "Choice",
  "Choices": [
    {
      "Variable": "$.message.type",
      "StringEquals": "thankYou",
      "Next": "WaitForThankYouWindow"
    }
  ],
  "Default": "RequestApprovalToSend"
},
"WaitForThankYouWindow": {
  "Type": "Wait",
  "TimestampPath": "$.targetSendTime",
  "Next": "RequestApprovalToSend"
}
```

No changes needed to `RequestApprovalToSend`, `SendAndPinMessage`, `RecordMessage`, or `NotifyCompletion` — they work as-is for the ThankYou type.

### Mock server

**`internal/mockserver/routes/mockresponses/schedulemock.go`**
- Add `StartDateTimeUTC` and `DurationHours` to `ProjectConfig`
- Default `DurationHours` to `2.0`; default `StartDateTimeUTC` to project `Date` at `14:00:00Z`
- Populate both fields in mock `UpcomingSession`

### Thank-you message template

**`seed/s3Items/thankyou.md`** (single global file, not per-project)

Content:
```
Thank you for leading {{projectName}} today! Your volunteers really appreciate your dedication.
```

The `{{projectName}}` placeholder is replaced at runtime with the actual project name. For welcome/reminder templates (which don't contain the placeholder), the replacement is a no-op.

### Interpolation points

**`internal/app/requestapproval/usecase.go`**
- After fetching S3 content, add: `messageContent = strings.ReplaceAll(messageContent, "{{projectName}}", projectName)`
- `projectName` is already a parameter of `Execute()` — no signature change needed

**`internal/app/sendandpinmessage/usecase.go`**
- Add `projectName string` to `Execute()` signature
- After `GetMessageContent`, add: `messageContent = strings.ReplaceAll(messageContent, "{{projectName}}", projectName)`
- Update handler (`lambda/sendandpinmessage/handler.go`) to pass `input.ExistingProjectNotification.Name`

### S3 path for thank-you

**`internal/app/computemessage/usecase.go`** — `computeS3MessageRefPath`
- Currently all types use `s3://{bucket}/{kebab-project-name}/{type}.md`
- For `ThankYou`: return `s3://{bucket}/thankyou.md` (global, no per-project subdirectory)

---

## Key Decisions

| Decision | Rationale |
|---|---|
| **End time source**: `Start_DateTime__c` + `Duration__c` | More reliable than `Session_End_Time__c`, which uses a misleading `Z` suffix for local time. `Start_DateTime__c` is already a proper UTC RFC3339 timestamp. |
| **Wait state with past timestamp** | Step Functions `Wait` with a past `TimestampPath` completes immediately — no special casing needed for projects that ended before the noon trigger fires. |
| **No new Lambda** | The thank-you reuses `SendAndPinMessage`, `RecordMessage`, and `NotifyCompletion` unchanged. |
| **`ProjectPassed` restructure** | Current logic fires `ProjectPassed` on the project day, blocking thank-you. New logic splits on `now.Before(startDate)` vs. after. |
| **Single global template** | One `thankyou.md` in S3 with `{{projectName}}` interpolated at runtime. No per-project authoring required. |

---

## Verification

1. `go test ./...` — unit tests pass including new ThankYou cases in `computemessage`
2. Add integration test scenario: mock project dated today ending shortly after "now", approve thank-you, confirm `HasSentThankYou = true` in DynamoDB
3. Manual smoke test: trigger workflow with today's mock project, observe `WaitForThankYouWindow` state in Step Functions console, approve via email, confirm message posted and pinned in Chime
