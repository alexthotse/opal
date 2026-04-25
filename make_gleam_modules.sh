#!/bin/bash
cd opal_backend/src

cat << 'INNER' > memory.gleam
pub fn extract_memories() -> String { "memories_extracted" }
pub fn compaction_reminders() -> String { "compaction_reminders_sent" }
INNER

cat << 'INNER' > cache.gleam
pub fn cached_microcompact() -> String { "microcompact_cached" }
INNER

cat << 'INNER' > teammem.gleam
pub fn get_teammem() -> String { "teammem_retrieved" }
INNER

cat << 'INNER' > verification.gleam
pub fn verify_task() -> String { "task_verified" }
INNER

cat << 'INNER' > triggers.gleam
pub fn run_triggers() -> String { "triggers_executed" }
INNER

cat << 'INNER' > search.gleam
pub fn quick_search() -> String { "search_completed" }
INNER

cat << 'INNER' > budget.gleam
pub fn check_budget() -> String { "budget_ok" }
INNER

cat << 'INNER' > stats.gleam
pub fn get_stats() -> String { "stats_retrieved" }
INNER

cat << 'INNER' > bridge.gleam
pub fn bridge_mode() -> String { "bridge_active" }
INNER

cat << 'INNER' > security.gleam
pub fn bash_classifier() -> String { "bash_classified_safe" }
INNER

cat << 'INNER' > planning.gleam
pub fn start_ultraplan() -> String { "ultraplan_mode_activated" }
INNER

cat << 'INNER' > reasoning.gleam
pub fn start_ultrathink() -> String { "ultrathink_mode_activated" }
INNER
