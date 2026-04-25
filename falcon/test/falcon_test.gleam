import gleeunit
import gleeunit/should
import gleam/string
import domain/reasoning
import domain/planning

pub fn main() {
  gleeunit.main()
}

pub fn hello_world_test() {
  1
  |> should.equal(1)
}

pub fn ultrathink_test() {
  let result = reasoning.start_ultrathink()
  string.contains(result, "ULTRATHINK initialized")
  |> should.be_true()
}

pub fn ultraplan_test() {
  let result = planning.start_ultraplan()
  string.contains(result, "ULTRAPLAN initialized")
  |> should.be_true()
}
