ONE
I want to have the option to set a custom current date in the lambdas.
this will enable future testing

TWO
I want my docker compose to handle s3 seeding and standing up my mockserver

THREE
I want an integration test suite that will run against the stood up stack. This may require small changes to the mockserver. I want to test:

calls (in series) when...

- my projects are all too far away
- one project is in range for a welcome notification and it is denied
- one project is in range for a welcome notification and it is permitted (and all the following will be permitted)
- that project is in range for welcome and already got the welcome (nothing happens)
- that project is in range for a reminder notification
- that project is in range for reminder and already got the warning (nothing happens)

another viable test may be ensuring that no messages are sent once the project is in the past

FOUR
Would a make file be appropriate to call the integration test? I may also want to use it to run unit tests and/or formatting? IDK

FIVE
I want GHA CICD for this repo. The CICD should, on PR to main, sure the code is formatted and that it passes these integration tests.

On push to main, I want to deploy the CDK to AWS.

SIX
I want to uphaul my documenation. Maybe a mermaid diagram of how it works?
