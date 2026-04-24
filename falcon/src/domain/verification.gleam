import gleam/string
import gleam/list

pub type TestResult {
  TestResult(name: String, passed: Bool, error: String)
}

pub type VerificationSuite {
  VerificationSuite(suite_name: String, tests: List(TestResult))
}

pub fn verify_task() -> String {
  let suite = VerificationSuite(
    suite_name: "Backend Auth",
    tests: [
      TestResult("Login Valid", True, ""),
      TestResult("Login Invalid", True, ""),
      TestResult("Token Expiry", False, "Token did not expire after 1 hour")
    ]
  )

  let tests_str = suite.tests
    |> list.map(fn(t) { 
      let status = case t.passed {
        True -> "PASS"
        False -> "FAIL - " <> t.error
      }
      t.name <> ": " <> status 
    })
    |> string.join("\n")

  "Verification for suite '" <> suite.suite_name <> "':\n" <> tests_str
}
