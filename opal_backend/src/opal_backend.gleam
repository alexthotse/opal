import gleam/io
import gleam/json
import gleam/string
import gleam/dynamic/decode
import gleam/result
import gleam/dynamic

pub type Request {
  Request(id: String, method: String)
}

pub fn main() {
  io.println("{\"jsonrpc\": \"2.0\", \"method\": \"system.started\", \"params\": {}}")
  loop()
}

fn loop() {
  case gleam_erlang_ffi_read_line() {
    Ok(line) -> {
      let line = string.trim(line)
      case string.length(line) > 0 {
        True -> handle_request(line)
        False -> Nil
      }
      loop()
    }
    Error(_) -> {
      io.println("{\"jsonrpc\": \"2.0\", \"method\": \"system.stopped\", \"params\": {}}")
    }
  }
}

fn handle_request(raw_json: String) {
  let decoder = {
    use id <- decode.field("id", decode.string)
    use method <- decode.field("method", decode.string)
    decode.success(Request(id, method))
  }

  case json.parse(raw_json, decoder) {
    Ok(req) -> {
      let res = execute_method(req.method)
      let response = json.object([
        #("jsonrpc", json.string("2.0")),
        #("id", json.string(req.id)),
        #("result", json.string(res))
      ])
      io.println(json.to_string(response))
    }
    Error(_) -> {
      io.println("{\"jsonrpc\": \"2.0\", \"error\": {\"code\": -32600, \"message\": \"Invalid Request\"}}")
    }
  }
}

pub fn execute_method(method: String) -> String {
  case method {
    "ping" -> "pong"
    "agent.start" -> "agent_started_successfully"
    "agent.status" -> "idle"
    "ultrathink.start" -> "ultrathink_mode_activated"
    "ultraplan.start" -> "ultraplan_mode_activated"
    _ -> "unknown_method"
  }
}

@external(erlang, "io", "get_line")
fn gleam_erlang_ffi_read_line_raw(prompt: String) -> dynamic.Dynamic

fn gleam_erlang_ffi_read_line() -> Result(String, Nil) {
  let res = gleam_erlang_ffi_read_line_raw("")
  decode.run(res, decode.string)
  |> result.map_error(fn(_) { Nil })
}
