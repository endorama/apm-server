# Functional tests

This folder contains functional tests. These are the broadest set of tests we have, to test end to end entire functionalities related to APM Server.

Depedencies are isolated with a separate go module, `go.mod` and `go.sum`.

Each test should have its own folder (folder naming is still unclear, as they could get quite verbose).

Within each folder, a `infra/` folder is expected, containing Terraform code required to setup the **initial** state of the test. Further state changes will not be made through Terraform but through API calls.
Each folder should contain a `main_test.go`, the entrypoint for all the test cases for the specific scenario.
Additional folders are Go modules to be leveraged for composing test functionalities without repeting too much code. They should be big enough to be reusable but small enough to be easily composable.

// TODO: maybe "scenario" is a better name here? Would help clarify the goal for a specific folder.
