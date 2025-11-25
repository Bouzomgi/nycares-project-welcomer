# Project Notifier System

## Problem Statement

The goal is to build an automated **Project Notifier** system to manage communication for projects. For each project, the system should automatically send a **welcome message** (no earlier than one week before the project date) and a **reminder message** (no earlier than two days before the project date) to all attendees. Each project should receive only one of each message, and every time a message is sent, it should be **pinned** immediately afterward.

All actions—checking projects, sending messages, and pinning—require authentication, which is handled through a **login** process that returns a cookie used for future requests. Before any message is sent, the system should issue a **notification** requesting approval; messages proceed only after the approval is given.

The system should allow the ability to **disable messaging** for specific projects if they are canceled. Projects may share names, but the combination of **name and date** uniquely identifies each one. Each project will have a **custom message template**, and multiple projects may occur per week.

The system is intended to run **once daily**, handle all of this automatically, and send a **follow-up notification** once messages have been sent. The solution should use **AWS cloud-native services**.

---

## Project Notifier: Implementation Options

### 1. Single AWS Lambda Orchestrator (Simple + Low Cost)

**Architecture:**

- A **Lambda** function runs once daily (triggered by an **EventBridge rule**).
- It:

  1. Logs in and retrieves the auth cookie.
  2. Fetches upcoming projects.
  3. Determines which messages (welcome/reminder) need to be sent.
  4. Sends a **notification** (via SNS, email, or Slack webhook) listing pending approvals.
  5. Waits for approval (possibly stored in **DynamoDB**).
  6. Sends and pins messages once approved.
  7. Sends a **follow-up notification** after completion.

**Pros:**

- Simple and easy to maintain.
- Minimal AWS components.
- Very low operational cost.

**Cons:**

- Limited control flow — harder to pause execution while waiting for approval.
- Requires **custom retry logic**: if one message fails while others succeed, partial retries must be handled manually (tracking which messages were sent, pinned, or failed).
- Harder to isolate and resume from failed steps compared to Step Functions.

---

### 2. AWS Step Functions Workflow (Moderate Complexity, Better Control)

**Architecture:**

- A **Step Functions** state machine orchestrates the process:

  1. **Login Step** → get auth cookie.
  2. **Fetch Projects Step** → identify upcoming projects.
  3. **Check Message Timing Step** → decide whether to send welcome/reminder messages.
  4. **Send Notification Step** → alert for approval via SNS or Slack.
  5. **Wait State** → pauses execution until approval is received (can resume via Lambda callback).
  6. **Send & Pin Message Step**.
  7. **Record Message** → record message data.
  8. **Notify Completion Step**.

**Pros:**

- Built-in retries per state with exponential backoff.
- Can continue processing other projects even if one fails.
- Visual workflow in Step Functions console.
- Clear error handling and monitoring.

**Cons:**

- Slightly higher cost than a single Lambda.
- More setup and maintenance.
- Slightly more complex for small-scale workloads.

---

### 3. Event-Driven Modular System (Scalable + Extensible)

**Architecture:**

- **EventBridge** schedules a daily event that triggers a **“Fetch Projects” Lambda**.
- That Lambda emits **Project Events** (e.g., “WelcomeNeeded”, “ReminderNeeded”) into **EventBridge**.
- Separate **Notifier Lambdas** handle:

  - Sending approval notifications
  - Sending messages
  - Pinning messages
  - Sending follow-up notifications

**Pros:**

- Highly modular and scalable.
- Easy to extend to new message types or rules.
- Asynchronous approvals naturally fit this model.

**Cons:**

- More moving parts to maintain.
- Debugging is slightly more complex.
- Careful event design is required to avoid duplication or missed events.

---

### 4. Containerized Workflow (ECS Fargate)

**Architecture:**

- Run a **daily container task** on **ECS Fargate** that performs the full workflow:

  1. Logs in and fetches projects.
  2. Checks which messages are needed.
  3. Sends approval notifications.
  4. Waits for approval (via API Gateway/DynamoDB).
  5. Sends & pins messages.
  6. Sends follow-up notification.

**Pros:**

- Full control over runtime and libraries.
- Easier to integrate complex scripts or web automation.
- Single container can encapsulate all dependencies.

**Cons:**

- Higher cost than serverless Lambda approaches.
- Requires management of task definitions, networking, and security.
- Overkill unless workflow is very complex.

---

## Estimated Monthly Cost Breakdown (All Options)

| Service                 | Option 1: Lambda                    | Option 2: Step Functions       | Option 3: Event-Driven         | Option 4: ECS Fargate             | Notes                                      |
| ----------------------- | ----------------------------------- | ------------------------------ | ------------------------------ | --------------------------------- | ------------------------------------------ |
| **AWS Lambda**          | 1 daily invocation (~2–3 s, 512 MB) | ~7 steps/day × 30 days         | 5 Lambdas × 30 invocations/day | N/A (Fargate container)           | Memory/runtime estimates                   |
| **Cost**                | ~$0.00                              | ~$0.01                         | ~$0.10                         | N/A                               | Free tier covers most workloads            |
| **EventBridge**         | 1 daily trigger                     | 1 daily trigger                | ~150 events/month              | 1 daily trigger                   | $1 per 1M events; negligible at this scale |
| **DynamoDB**            | ~100 reads/writes/month             | ~200 reads/writes/month        | ~500 reads/writes/month        | ~200 reads/writes/month           | Track approvals & message state            |
| **SNS / Notifications** | ~30 notifications/month             | ~50 notifications/month        | ~150 notifications/month       | ~50 notifications/month           | Email/webhook; SMS may add cost            |
| **CloudWatch Logs**     | 50 MB                               | 200 MB                         | 200 MB                         | 100 MB                            | $0.50 per GB ingested                      |
| **Step Functions**      | N/A                                 | ~1,500 state transitions/month | N/A                            | N/A                               | $0.000025 per transition                   |
| **ECS Fargate**         | N/A                                 | N/A                            | N/A                            | 1 task/day, 512 MB, 1 vCPU, 5 min | $0.04048 per vCPU-h + $0.004445 per GB-h   |

**Approximate Monthly Totals:**

- **Option 1 (Lambda Orchestrator):** ~$0.03 – $0.10
- **Option 2 (Step Functions Workflow):** ~$0.05 – $0.10
- **Option 3 (Event-Driven Modular):** ~$0.20 – $0.30
- **Option 4 (ECS Fargate):** ~$0.27 – $0.30

**Notes:**

- All estimates assume low to moderate project volume (a few dozen per month).
- Costs scale linearly with additional projects, notifications, or logging.
- All options can leverage AWS Free Tier for minimal cost in small-scale usage.

---

## Chosen Architecture: AWS Step Functions Workflow

**Reasons for this choice:**

- **Built-in retries and error handling:** Each step in the workflow can automatically retry on failure with configurable backoff and maximum attempts, ensuring reliable message delivery.
- **Centralized orchestration:** The state machine provides a clear, visual representation of the workflow, making it easier to understand and maintain.
- **Human approval integration:** Step Functions allows a built-in **pause/wait state** for approval before messages are sent.
- **Fault isolation:** Failures in one step or for one project do not block other steps or projects from completing.
- **Extensibility:** New steps or logic can be added to the workflow without redesigning the entire system.

**Summary of Workflow Steps:**

1. **Login Step** → authenticate and obtain an auth cookie.
2. **Fetch Projects Step** → retrieve upcoming projects.
3. **Check Message Timing Step** → determine if welcome or reminder messages should be sent.
4. **Send Notification Step** → notify for approval.
5. **Wait State** → pause workflow until approval is received.
6. **Send & Pin Message Step** → send messages and pin them for each project.
7. **Record Message** → record message data in database.
8. **Notify Completion Step** → send a downstream notification confirming messages were sent.

This architecture balances **automation, reliability, and control**, ensuring multiple projects are managed efficiently and notifications are delivered correctly.

---

## Notification State Tracking in DynamoDB

A **single DynamoDB table** is used to track which messages have been sent for each project. Each row represents a unique project, identified by **name and date**, allowing the system to manage message sending, prevent duplicates, and stop notifications for canceled projects.

### Table Schema

| Attribute            | Type    | Description                                          |
| -------------------- | ------- | ---------------------------------------------------- |
| **ProjectName**      | String  | Name of the project.                                 |
| **ProjectDate**      | String  | Date of the project (ISO format).                    |
| **HasSentWelcome**   | Boolean | Indicates if the welcome message has been sent.      |
| **HasSentReminder**  | Boolean | Indicates if the reminder message has been sent.     |
| **ShouldStopNotify** | Boolean | Flag to disable notifications for canceled projects. |
| **LastUpdated**      | String  | ISO timestamp of the last update to the row.         |

### Usage Notes

- **Primary Key:** A composite key, such as `ProjectName#ProjectDate`, ensures uniqueness for each project instance.
- **Upserts:** Each Lambda performs an upsert (`PutItem` or `UpdateItem`) when messages are sent or a project is canceled, ensuring atomic updates.

---

## Handling Auth Cookie Expiration Mid-Workflow

To ensure reliability if the **auth cookie expires during execution**, the workflow integrates an **auth refresh mechanism**:

### How It Works

- **Lambda Steps** that perform authenticated actions (fetch projects, send/pin messages) attempt their requests using the current cookie.
- If the request fails due to authentication (e.g., 401 Unauthorized), the Lambda **calls the login API internally to refresh the cookie** and retries the action automatically.
- Step Functions continues to handle other failures (network errors, API timeouts) using its built-in retry mechanisms.

---

## Step Functions Input/Output Table

| Step                      | Notes                                                                                          |
| ------------------------- | ---------------------------------------------------------------------------------------------- |
| **Login**                 | Returns auth cookie for subsequent requests.                                                   |
| **FetchProjects**         | Returns list of projects with **template reference** and message type.                         |
| **ComputeMessageToSend**  | Reads DynamoDB table and determines **single eligible message** to send.                       |
| **RequestApprovalToSend** | Notifier Lambda receives project info **and template reference**; approval triggers next step. |
| **WaitForApproval**       | Passes input unchanged until approval.                                                         |
| **SendAndPinMessage**     | Lambda fetches template from S3, sends message, pins it. Throws error if send/pin fails.       |
| **RecordMessage**         | Updates DynamoDB table to track that the message has been sent for the project.                |
| **NotifyCompletion**      | Sends completion notification with **project info, message type, and template reference**.     |
| **DLQNotifier**           | Sends failure notification for manual resolution.                                              |

## Step Functions Inputs and Outputs

### Login

**Input:** `{}`

**Output:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  }
}
```

### FetchProjects

**Input:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  }
}
```

**Output:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "projects": [
    {
      "name": "Project A",
      "date": "2025-11-10"
    }
  ]
}
```

### ComputeMessageToSend

**Input:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "project": {
    "name": "Project A",
    "date": "2025-11-10"
  }
}
```

**Output:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "project": {
    "name": "Project A",
    "date": "2025-11-10"
  },
  "messageToSend": {
    "type": "welcome",
    "templateRef": "s3://bucket/projectA/welcome.md"
  }
}
```

### RequestApprovalToSend

**Input:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "project": {
    "name": "Project A",
    "date": "2025-11-10"
  },
  "messageToSend": {
    "type": "welcome",
    "templateRef": "s3://bucket/projectA/welcome.md"
  }
}
```

**Output:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "approvedToSend": "true",
  "project": {
    "name": "Project A",
    "date": "2025-11-10"
  },
  "messageToSend": {
    "type": "welcome",
    "templateRef": "s3://bucket/projectA/welcome.md"
  }
}
```

### WaitForApproval

**Input/Output:** Unchanged from previous step.

### SendAndPinMessage

**Input:**

**Input/Output:** Unchanged from previous step.

### NotifyCompletion

**Input:**

```json
{
  "auth": {
    "cookie": {
      "name": "session_id",
      "value": "abc123-session",
      "domain": "example.com",
      "path": "/"
    }
  },
  "project": {
    "name": "Project A",
    "date": "2025-11-10"
  },
  "messageToSend": {
    "type": "welcome",
    "templateRef": "s3://bucket/projectA/welcome.md"
  }
}
```

**Output:** `n/a`

### DLQNotifier

**Input:**

```json
{
  "errorType": "SendPinFailed",
  "errorMessage": "Failed to send or pin",
  "projectName": "Project A",
  "projectDate": "2025-11-10",
  "messageToSend": {
    "type": "welcome",
    "templateRef": "s3://bucket/projectA/welcome.md"
  }
}
```

**Output:** `n/a`
