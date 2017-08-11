INSERT INTO users (uuid, email) VALUES ('some-user', 'foo@example.com');
INSERT INTO locations (uuid, userid, name, lat, long) VALUES ('some-location', 'some-user', 'Buckingham Palace', 51.501364, -0.141890);
INSERT INTO tasks (uuid, userid, locationid, type, title, description, time) VALUES ( 'some-task', 'some-user', 'some-location', 1, 'Take a quick photo', 'Get a photo for mum', '');
