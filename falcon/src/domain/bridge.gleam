import gleam/string
import gleam/list

pub type BridgeConnection {
  BridgeConnection(service: String, connected: Bool, latency_ms: Int)
}

pub fn bridge_mode() -> String {
  let connections = [
    BridgeConnection("Stripe", True, 45),
    BridgeConnection("SendGrid", True, 120),
    BridgeConnection("LegacyDB", False, 0)
  ]
  
  let conn_str = connections
    |> list.map(fn(c) { 
      let status = case c.connected {
        True -> "UP (" <> string.inspect(c.latency_ms) <> "ms)"
        False -> "DOWN"
      }
      c.service <> ": " <> status
    })
    |> string.join("\n")

  "Bridge Mode Status:\n" <> conn_str
}
