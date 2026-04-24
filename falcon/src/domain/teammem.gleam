import gleam/string
import gleam/list

pub type TeamMember {
  TeamMember(name: String, role: String, status: String)
}

pub fn get_teammem() -> String {
  let members = [
    TeamMember("Alice", "Frontend", "Active"),
    TeamMember("Bob", "Backend", "Reviewing PR")
  ]
  let members_str = members
    |> list.map(fn(m) { m.name <> " (" <> m.role <> "): " <> m.status })
    |> string.join("\n")
  "Team Members Context:\n" <> members_str
}
