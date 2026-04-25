import gleam/string
import gleam/list

pub type TriggerEvent {
  TriggerEvent(name: String, condition: String, active: Bool)
}

pub fn run_triggers() -> String {
  let triggers = [
    TriggerEvent("OnCommit", "branch == main", True),
    TriggerEvent("OnDeploy", "env == prod", False)
  ]
  let trigs_str = triggers
    |> list.map(fn(t) { 
      let active_str = case t.active {
        True -> "Active"
        False -> "Inactive"
      }
      "Trigger " <> t.name <> " [" <> active_str <> "]: " <> t.condition
    })
    |> string.join("\n")

  "Executed Triggers:\n" <> trigs_str
}
