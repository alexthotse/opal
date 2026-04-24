import gleeunit
import gleeunit/should
import opal_backend

pub fn main() {
  gleeunit.main()
}

// gleeunit test functions end in `_test`
pub fn hello_world_test() {
  1
  |> should.equal(1)
}

pub fn ultrathink_test() {
  opal_backend.execute_method("ultrathink.start")
  |> should.equal("ultrathink_mode_activated")
}

pub fn ultraplan_test() {
  opal_backend.execute_method("ultraplan.start")
  |> should.equal("ultraplan_mode_activated")
}
