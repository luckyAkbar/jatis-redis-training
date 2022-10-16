-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "data" (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "menu" (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id BIGINT default NULL,
    url TEXT default NULL
);

INSERT INTO "menu" (id, name) VALUES (1, 'menu satu');
INSERT INTO "menu" (id, name, parent_id) VALUES (2, 'menu dua', 1);
INSERT INTO "menu" (id, name) VALUES (3, 'menu tiga');
INSERT INTO "menu" (id, name) VALUES (4, 'menu empat');
INSERT INTO "menu" (id, name) VALUES (5, 'menu lima');
INSERT INTO "menu" (id, name) VALUES (6, 'menu enam');
INSERT INTO "menu" (id, name, parent_id) VALUES (7, 'menu tujuh', 4);
INSERT INTO "menu" (id, name) VALUES (8,'menu delapan');
INSERT INTO "menu" (id, name) VALUES (9, 'menu sembilan');
INSERT INTO "menu" (id, name) VALUES (10, 'menu sepuluh');

INSERT INTO "data" (id, name) VALUES (1, 'data satu');
INSERT INTO "data" (id, name) VALUES (2, 'data dua');
INSERT INTO "data" (id, name) VALUES (3, 'data tiga');
INSERT INTO "data" (id, name) VALUES (4, 'data empat');
INSERT INTO "data" (id, name) VALUES (5, 'data lima');
INSERT INTO "data" (id, name) VALUES (6, 'data enam');
INSERT INTO "data" (id, name) VALUES (7, 'data tujuh');
INSERT INTO "data" (id, name) VALUES (8, 'data delapan');
INSERT INTO "data" (id, name) VALUES (9, 'data sembilan');
INSERT INTO "data" (id, name) VALUES (10, 'data sepuluh');

-- +migrate Down
DROP TABLE IF EXISTS "data";
DROP TABLE IF EXISTS "menus";