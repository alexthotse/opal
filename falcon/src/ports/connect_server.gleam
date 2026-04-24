import gleam/dynamic/decode
import gleam/http/response
import gleam/io
import gleam/json
import gleam/string
import msgpack
import wisp

pub fn handle_request(req: wisp.Request) -> wisp.Response {
  // We expect application/msgpack per the Connect RPC spec over custom codec
  use <- wisp.require_content_type(req, "application/msgpack")
  use body <- wisp.require_bit_array_body(req)

  let path = req.path

  // A naive routing logic for Connect RPC:
  // POST /falcon.v1.FalconService/Ping
  let result = case path {
    "/falcon.v1.FalconService/Ping" -> "pong"
    "/falcon.v1.FalconService/StartUltrathink" -> "ultrathink_started"
    "/falcon.v1.FalconService/StartUltraplan" -> "ultraplan_started"
    "/falcon.v1.FalconService/DispatchAction" -> "action_dispatched"
    "/falcon.v1.FalconService/QuickSearch" -> "search_completed"
    _ -> "not_found"
  }

  // Create a minimal Map-like structure or just pack a string/tuple depending on what Connect expects
  // Usually Connect RPC expects the exact Protobuf message structure encoded via MsgPack.
  // For PingResponse { string message = 1; }, MsgPack map {"message": "pong"} or array.
  // We'll pack a simple map for demonstration.
  let response_payload =
    msgpack.pack(
      msgpack.Map([
        #(msgpack.Str("result"), msgpack.Str(result)),
        #(msgpack.Str("message"), msgpack.Str(result)),
      ]),
    )

  wisp.response(200)
  |> wisp.set_header("content-type", "application/msgpack")
  |> wisp.set_body(wisp.BitArrayBody(response_payload))
}
