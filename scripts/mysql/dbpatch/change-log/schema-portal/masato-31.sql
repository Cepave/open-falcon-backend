ALTER TABLE `event_cases`
  CHANGE timestamp timestamp Timestamp NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE `event_cases`
  CHANGE update_at update_at Timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

ALTER TABLE `event_note`
  CHANGE timestamp timestamp Timestamp NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE `events`
  CHANGE timestamp timestamp Timestamp NULL DEFAULT CURRENT_TIMESTAMP;
