import gleam/io
import gleam/string
import gleam/dynamic
import msgpack
import wisp
import gleam/http/response

import domain/reasoning
import domain/planning
import domain/memory
import adapters/jido_agent

pub fn handle_request(req: wisp.Request) -> wisp.Response {
  // Ensure application/msgpack
  use <- wisp.require_content_type(req, "application/msgpack")
  use body <- wisp.require_bit_array_body(req)

  // Unpack the MsgPack to route the message (ignoring parsing for simple mock demonstration,
  // but a real implementation would decode body into the protobuf struct).
  let path = req.path
  
  // A naive routing logic for Connect RPC:
  let result = case path {
    "/falcon.v1.FalconService/Ping" -> "pong"
    "/falcon.v1.FalconService/StartUltrathink" -> reasoning.start_ultrathink()
    "/falcon.v1.FalconService/StartUltraplan" -> planning.start_ultraplan()
    "/falcon.v1.FalconService/DispatchAction" -> jido_agent.dispatch_action("test_action")
    "/falcon.v1.FalconService/QuickSearch" -> "search_completed"
    _ -> "not_found"
  }

  // ConnectRPC expects a Protobuf-like struct representation. 
  // In MessagePack, Go's struct is usually mapped to an array of field values or map.
  // We'll pack a simple map for now that maps to `message { string result = 1; }`
  // Actually, msgpack/v5 in Go decodes structs as Maps by default, keyed by struct field names.
  let response_payload = msgpack.pack(msgpack.Map([
    #(msgpack.Str("Result"), msgpack.Str(result)),
    #(msgpack.Str("Message"), msgpack.Str(result))
  ]))

  wisp.response(200)
  |> wisp.set_header("content-type", "application/msgpack")
  |> wisp.set_body(wisp.BitArrayBody(response_payload))
}
