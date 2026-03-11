ALTER TABLE attendees DROP INDEX idx_attendees_unique;
ALTER TABLE attendees DROP INDEX idx_attendees_event_id;
ALTER TABLE attendees DROP INDEX idx_attendees_user_id;
ALTER TABLE events DROP INDEX idx_events_owner_id;
