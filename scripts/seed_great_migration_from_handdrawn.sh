#!/usr/bin/env bash
set -euo pipefail

# Seed "The Great Migration" backend from the handwritten schedule.
#
# Usage:
#   chmod +x seed_great_migration_from_handdrawn.sh
#   BASE_URL=http://localhost:3000 ./seed_great_migration_from_handdrawn.sh
#
# Or:
#   ./seed_great_migration_from_handdrawn.sh http://localhost:3000
#
# Requirements:
#   - bash
#   - curl
#   - jq
#
# Notes:
#   - This script assumes an empty/dev database. The API has create endpoints but no
#     person/task upsert-by-name endpoints, so re-running may create duplicates.
#   - Dates are configurable. The handwritten plan only says 2–5 July, so this defaults
#     to 2026 because your API description mentions a 2026 planning window.
#   - Set START_YEAR=2024 if you want the mockup-style 2024 dates instead.

BASE_URL="${1:-${BASE_URL:-http://localhost:3000}}"
BASE_URL="${BASE_URL%/}"

START_YEAR="${START_YEAR:-2026}"
START_DATE="${START_DATE:-${START_YEAR}-07-02}"
END_DATE="${END_DATE:-${START_YEAR}-07-05}"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing required command: $1" >&2
    exit 1
  }
}

require_cmd curl
require_cmd jq

api_json() {
  local method="$1"
  local path="$2"
  local body="${3:-}"

  if [[ -n "$body" ]]; then
    curl --silent --show-error --fail-with-body \
      -X "$method" \
      -H "Content-Type: application/json" \
      -H "Accept: application/json" \
      --data "$body" \
      "$BASE_URL$path"
  else
    curl --silent --show-error --fail-with-body \
      -X "$method" \
      -H "Accept: application/json" \
      "$BASE_URL$path"
  fi
}

date_for_day() {
  local day="$1"
  printf "%s-07-%02d" "$START_YEAR" "$day"
}

declare -A PERSON_ID
declare -A TASK_ID

echo "Seeding The Great Migration into: $BASE_URL"
echo "Planning window: $START_DATE to $END_DATE"
echo

echo "1/6 Updating planning window..."
api_json PUT "/api/planning-window" \
  "$(jq -n --arg startDate "$START_DATE" --arg endDate "$END_DATE" \
    '{startDate:$startDate,endDate:$endDate}')" >/dev/null

create_person() {
  local key="$1"
  local name="$2"
  local initials="$3"

  local response
  response="$(api_json POST "/api/people" \
    "$(jq -n --arg name "$name" --arg initials "$initials" \
      '{name:$name,initials:$initials}')")"

  PERSON_ID["$key"]="$(jq -r '.id' <<<"$response")"
  printf "  created person %-10s -> %s\n" "$name" "${PERSON_ID[$key]}"
}

echo "2/6 Creating people..."
create_person "estella"  "Estella"  "E"
create_person "meryl"    "Meryl"    "M"
create_person "matthijs" "Matthijs" "M"
create_person "thomas"   "Thomas"   "T"
create_person "esther"   "Esther"   "E"
create_person "peter"    "Peter"    "P"
create_person "beau"     "Beau"     "B"
create_person "kees"     "Kees"     "K"

create_room() {
  local name="$1"
  local type="$2"

  api_json POST "/api/rooms" \
    "$(jq -n --arg name "$name" --arg type "$type" \
      '{name:$name,type:$type}')" >/dev/null

  printf "  created %-4s %s\n" "$type" "$name"
}

echo "3/6 Creating rooms and areas..."
create_room "Hal" "area"
create_room "Woonkamer" "room"
create_room "Keuken" "room"
create_room "Slaapkamer" "room"
create_room "Tweede verdieping" "area"
create_room "Algemeen" "area"

set_availability() {
  local person_key="$1"
  local date="$2"
  local status="$3"
  local person_id="${PERSON_ID[$person_key]}"

  api_json PUT "/api/people/${person_id}/availability/${date}" \
    "$(jq -n --arg status "$status" '{status:$status}')" >/dev/null
}

echo "4/6 Setting availability from handwritten schedule..."

# 2 July: Estella, Meryl, Matthijs, Thomas, Esther, Peter?, Beau.
set_availability estella  "$(date_for_day 2)" available
set_availability meryl    "$(date_for_day 2)" available
set_availability matthijs "$(date_for_day 2)" available
set_availability thomas   "$(date_for_day 2)" available
set_availability esther   "$(date_for_day 2)" available
set_availability peter    "$(date_for_day 2)" partial
set_availability beau     "$(date_for_day 2)" available
set_availability kees     "$(date_for_day 2)" off

# 3 July: Thomas, Matthijs, Estella, Beau?
set_availability estella  "$(date_for_day 3)" available
set_availability meryl    "$(date_for_day 3)" off
set_availability matthijs "$(date_for_day 3)" available
set_availability thomas   "$(date_for_day 3)" available
set_availability esther   "$(date_for_day 3)" off
set_availability peter    "$(date_for_day 3)" off
set_availability beau     "$(date_for_day 3)" partial
set_availability kees     "$(date_for_day 3)" off

# 4 July: Thomas, Matthijs, Estella, Esther, Peter, Kees.
set_availability estella  "$(date_for_day 4)" available
set_availability meryl    "$(date_for_day 4)" off
set_availability matthijs "$(date_for_day 4)" available
set_availability thomas   "$(date_for_day 4)" available
set_availability esther   "$(date_for_day 4)" available
set_availability peter    "$(date_for_day 4)" available
set_availability beau     "$(date_for_day 4)" off
set_availability kees     "$(date_for_day 4)" available

# 5 July: Esther, Thomas, Matthijs, Kees, Estella, Peter.
set_availability estella  "$(date_for_day 5)" available
set_availability meryl    "$(date_for_day 5)" off
set_availability matthijs "$(date_for_day 5)" available
set_availability thomas   "$(date_for_day 5)" available
set_availability esther   "$(date_for_day 5)" available
set_availability peter    "$(date_for_day 5)" available
set_availability beau     "$(date_for_day 5)" off
set_availability kees     "$(date_for_day 5)" available

create_task() {
  local key="$1"
  local title="$2"
  local priority="$3"
  local people_needed="$4"
  local room="$5"
  local status="${6:-ready}"

  local response
  response="$(api_json POST "/api/tasks" \
    "$(jq -n \
      --arg title "$title" \
      --arg priority "$priority" \
      --argjson peopleNeeded "$people_needed" \
      --arg room "$room" \
      --arg status "$status" \
      --argjson assignedTo '[]' \
      '{
        title:$title,
        priority:$priority,
        peopleNeeded:$peopleNeeded,
        room:$room,
        status:$status,
        assignedTo:$assignedTo
      }')")"

  TASK_ID["$key"]="$(jq -r '.id' <<<"$response")"
  printf "  created task %-38s -> %s\n" "$title" "${TASK_ID[$key]}"
}

echo "5/6 Creating backlog tasks..."
create_task "afplakken"                          "Afplakken"                              "medium" 1 "Algemeen"
create_task "schuren-hal"                        "Schuren"                                "medium" 1 "Hal"
create_task "schilderen-hal"                     "Schilderen hal"                         "high"   2 "Hal"
create_task "schoonmaken"                        "Schoonmaken"                            "low"    1 "Algemeen"
create_task "stoom-muren"                        "Stoom muren"                            "high"   1 "Algemeen"
create_task "schilderen-3-juli"                  "Schilderen"                             "high"   2 "Hal"
create_task "schuren-schilderen-keuken"          "Schuren / schilderen keuken"            "medium" 1 "Keuken"
create_task "schilderen-hal-woonkamer"           "Schilderen hal / woonkamer"             "high"   2 "Hal"
create_task "stoom-muren-woonkamer"              "Stoom muren woonkamer"                  "high"   1 "Woonkamer"
create_task "schoonmaken-eten"                   "Schoonmaken / eten"                     "low"    1 "Algemeen"
create_task "slaapkamer-schilderen"              "Slaapkamer schilderen"                  "high"   2 "Slaapkamer"
create_task "tweede-verdieping-muren-schuren"    "2e verdieping muren schuren"            "medium" 1 "Tweede verdieping"
create_task "hal-woonkamer-afronden"             "Hal / woonkamer afronden"               "high"   2 "Hal"
create_task "tweede-verdieping-muren-schilderen" "2e verdieping muren schilderen"         "high"   1 "Tweede verdieping"

person_ids_json() {
  local ids='[]'
  local key

  for key in "$@"; do
    ids="$(jq -c --arg id "${PERSON_ID[$key]}" '. + [$id]' <<<"$ids")"
  done

  printf "%s" "$ids"
}

create_schedule_card() {
  local date="$1"
  local task_key="$2"
  shift 2

  local assigned_to
  assigned_to="$(person_ids_json "$@")"

  local response
  response="$(api_json POST "/api/schedule/cards" \
    "$(jq -n \
      --arg scheduledDate "$date" \
      --arg taskId "${TASK_ID[$task_key]}" \
      --argjson assignedTo "$assigned_to" \
      '{
        scheduledDate:$scheduledDate,
        taskId:$taskId,
        assignedTo:$assignedTo
      }')")"

  local id title count needed
  id="$(jq -r '.id' <<<"$response")"
  title="$(jq -r '.title' <<<"$response")"
  count="$(jq -r '.assignedCount' <<<"$response")"
  needed="$(jq -r '.peopleNeeded' <<<"$response")"

  printf "  scheduled %-34s on %s -> %s (%s/%s)\n" "$title" "$date" "$id" "$count" "$needed"
}

echo "6/6 Creating schedule cards..."

# 2 July
create_schedule_card "$(date_for_day 2)" "afplakken"                 estella
create_schedule_card "$(date_for_day 2)" "schuren-hal"               meryl
create_schedule_card "$(date_for_day 2)" "schilderen-hal"            matthijs thomas
create_schedule_card "$(date_for_day 2)" "schoonmaken"               esther
create_schedule_card "$(date_for_day 2)" "stoom-muren"               beau

# 3 July
create_schedule_card "$(date_for_day 3)" "schilderen-3-juli"         thomas matthijs
create_schedule_card "$(date_for_day 3)" "stoom-muren"               estella
create_schedule_card "$(date_for_day 3)" "schuren-schilderen-keuken" beau

# 4 July
create_schedule_card "$(date_for_day 4)" "schilderen-hal-woonkamer"  thomas matthijs
create_schedule_card "$(date_for_day 4)" "stoom-muren-woonkamer"     estella
create_schedule_card "$(date_for_day 4)" "schoonmaken-eten"          esther
create_schedule_card "$(date_for_day 4)" "slaapkamer-schilderen"     peter kees

# 5 July
create_schedule_card "$(date_for_day 5)" "tweede-verdieping-muren-schuren"    esther
create_schedule_card "$(date_for_day 5)" "hal-woonkamer-afronden"             thomas matthijs
create_schedule_card "$(date_for_day 5)" "slaapkamer-schilderen"              kees peter
create_schedule_card "$(date_for_day 5)" "tweede-verdieping-muren-schilderen" estella

echo
echo "Seed complete."
echo
echo "Useful checks:"
echo "  curl -s '$BASE_URL/api/planning-window' | jq"
echo "  curl -s '$BASE_URL/api/dashboard/people-availability?start=$START_DATE&days=4' | jq"
echo "  curl -s '$BASE_URL/api/tasks/backlog' | jq"
echo "  curl -s '$BASE_URL/api/dashboard/daily-schedule?start=$START_DATE&days=4' | jq"
