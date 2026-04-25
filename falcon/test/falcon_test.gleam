import gleeunit
import gleeunit/should
import gleam/string
import domain/reasoning
import domain/planning
import domain/stats
import domain/teammem
import domain/security
import simplifile

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

pub fn stats_test() {
  let result = stats.get_stats()
  string.contains(result, "System Stats:")
  |> should.be_true()
  string.contains(result, "CPU")
  |> should.be_true()
  string.contains(result, "RAM")
  |> should.be_true()
}

pub fn teammem_test() {
  let _ = simplifile.write("teammem.md", "Real Team Member Data")
  let result = teammem.get_teammem()
  string.contains(result, "Real Team Member Data")
  |> should.be_true()
  let _ = simplifile.delete("teammem.md")
}

pub fn security_safe_test() {
  let result = security.bash_classifier("ls -la")
  string.contains(result, "ALLOW")
  |> should.be_true()
}

pub fn security_unsafe_rm_test() {
  let result = security.bash_classifier("rm -rf /")
  string.contains(result, "BLOCK")
  |> should.be_true()
}

pub fn security_unsafe_sudo_test() {
  let result = security.bash_classifier("sudo su")
  string.contains(result, "BLOCK")
  |> should.be_true()
}
