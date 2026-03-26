-- +goose Up
CREATE TABLE materials (
    guid TEXT PRIMARY KEY,
    wow_id INTEGER NOT NULL,
    name TEXT NOT NULL
) WITHOUT ROWID;

CREATE TABLE recipes (
    guid TEXT PRIMARY KEY,
    wow_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    profession TEXT NOT NULL
) WITHOUT ROWID;

CREATE TABLE characters (
    guid TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    server TEXT NOT NULL,
    guild TEXT,
    score REAL,
    level INTEGER NOT NULL
) WITHOUT ROWID;

CREATE TABLE recipe_materials (
    recipe_id TEXT,
    material_id TEXT,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, material_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(guid) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(guid) ON DELETE CASCADE
) WITHOUT ROWID;

CREATE TABLE character_recipes (
    character_id TEXT,
    recipe_id TEXT,
    PRIMARY KEY (character_id, recipe_id),
    FOREIGN KEY (character_id) REFERENCES characters(guid) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes(guid) ON DELETE CASCADE
) WITHOUT ROWID;

-- +goose Down
DROP TABLE character_recipes;
DROP TABLE recipe_materials;
DROP TABLE characters;
DROP TABLE recipes;
DROP TABLE materials;
