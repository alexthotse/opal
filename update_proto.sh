sed -i '/^}/i \
  rpc GetTokenBudget(BudgetRequest) returns (BudgetResponse) {}\
  rpc GetTeamMem(TeamMemRequest) returns (TeamMemResponse) {}\
  rpc StartBridgeMode(BridgeModeRequest) returns (BridgeModeResponse) {}\
  rpc SetAgentTrigger(TriggerRequest) returns (TriggerResponse) {}\
  rpc SetAgentTriggerRemote(TriggerRemoteRequest) returns (TriggerRemoteResponse) {}\
  rpc ExtractMemories(ExtractMemoriesRequest) returns (ExtractMemoriesResponse) {}\
  rpc GetCompactionReminders(CompactionRequest) returns (CompactionResponse) {}\
  rpc CachedMicrocompact(MicrocompactRequest) returns (MicrocompactResponse) {}\
  rpc GetShotStats(ShotStatsRequest) returns (ShotStatsResponse) {}\
  rpc StartVerification(VerificationRequest) returns (VerificationResponse) {}\
  rpc ClassifyBash(BashClassifyRequest) returns (BashClassifyResponse) {}\
  rpc GetVoiceMode(VoiceModeRequest) returns (VoiceModeResponse) {}\
  rpc GetHistoryPicker(HistoryPickerRequest) returns (HistoryPickerResponse) {}\
  rpc MessageActions(MessageActionsRequest) returns (MessageActionsResponse) {}' /workspace/proto/falcon/v1/falcon.proto

cat << 'PROTO' >> /workspace/proto/falcon/v1/falcon.proto

message BudgetRequest { string id = 1; }
message BudgetResponse { string result = 1; }

message TeamMemRequest { string id = 1; }
message TeamMemResponse { string result = 1; }

message BridgeModeRequest { string id = 1; }
message BridgeModeResponse { string result = 1; }

message TriggerRequest { string id = 1; }
message TriggerResponse { string result = 1; }

message TriggerRemoteRequest { string id = 1; }
message TriggerRemoteResponse { string result = 1; }

message ExtractMemoriesRequest { string id = 1; }
message ExtractMemoriesResponse { string result = 1; }

message CompactionRequest { string id = 1; }
message CompactionResponse { string result = 1; }

message MicrocompactRequest { string id = 1; }
message MicrocompactResponse { string result = 1; }

message ShotStatsRequest { string id = 1; }
message ShotStatsResponse { string result = 1; }

message VerificationRequest { string id = 1; }
message VerificationResponse { string result = 1; }

message BashClassifyRequest { string id = 1; }
message BashClassifyResponse { string result = 1; }

message VoiceModeRequest { string id = 1; }
message VoiceModeResponse { string result = 1; }

message HistoryPickerRequest { string id = 1; }
message HistoryPickerResponse { string result = 1; }

message MessageActionsRequest { string id = 1; }
message MessageActionsResponse { string result = 1; }
PROTO
