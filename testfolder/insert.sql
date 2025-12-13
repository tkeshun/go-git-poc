INSERT INTO demo_items (name, note)
VALUES
  ('apple',  'first seed'),
  ('banana', 'second seed')
ON CONFLICT (name) DO UPDATE
  SET note = EXCLUDED.note;