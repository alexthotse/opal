import simplifile

pub type TeamMember {
  TeamMember(name: String, role: String, status: String)
}

pub fn get_teammem() -> String {
  case simplifile.read("teammem.md") {
    Ok(content) -> "Team Members Context:\n" <> content
    Error(_) -> "Team Members Context:\nCould not read teammem.md"
  }
}
