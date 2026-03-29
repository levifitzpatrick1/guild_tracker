-- +goose Up
CREATE TABLE materials (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
) WITHOUT ROWID;

CREATE TABLE recipes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    profession TEXT NOT NULL
) WITHOUT ROWID;

CREATE TABLE characters (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    server TEXT NOT NULL,
    guild TEXT,
    score REAL,
    level INTEGER NOT NULL
) WITHOUT ROWID;

CREATE TABLE recipe_crafting_slots (
    recipe_id INTEGER NOT NULL,
    slot_name TEXT NOT NULL,
    display_order INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, display_order),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
) WITHOUT ROWID;

CREATE TABLE recipe_materials (
    recipe_id INTEGER NOT NULL,
    material_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, material_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
) WITHOUT ROWID;

CREATE TABLE character_recipes (
    character_id INTEGER NOT NULL,
    recipe_id INTEGER NOT NULL,
    PRIMARY KEY (character_id, recipe_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
) WITHOUT ROWID;

-- +goose Down
DROP TABLE character_recipes;
DROP TABLE recipe_materials;
DROP TABLE recipe_crafting_slots;
DROP TABLE characters;
DROP TABLE recipes;
DROP TABLE materials;
