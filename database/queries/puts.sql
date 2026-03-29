-- name: UpsertMaterial :one
INSERT INTO materials (
    id, name
) VALUES (
    ?, ?
)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name
RETURNING *;

-- name: UpsertRecipe :one
INSERT INTO recipes (
    id, name, profession
) VALUES (
    ?, ?, ?
)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    profession = excluded.profession
RETURNING *;

-- name: UpsertCharacter :one
INSERT INTO characters (
    id, name, server, guild, score, level
) VALUES (
    ?, ?, ?, ?, ?, ?
)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    server = excluded.server,
    guild = excluded.guild,
    score = excluded.score,
    level = excluded.level
RETURNING *;

-- name: UpsertRecipeMaterial :exec
INSERT INTO recipe_materials (
    recipe_id, material_id, quantity
) VALUES (
    ?, ?, ?
)
ON CONFLICT(recipe_id, material_id) DO UPDATE SET
    quantity = excluded.quantity;

-- name: UpsertCharacterRecipe :exec
INSERT INTO character_recipes (
    character_id, recipe_id
) VALUES (
    ?, ?
)
ON CONFLICT(character_id, recipe_id) DO NOTHING;

-- name: UpsertRecipeCraftingSlot :exec
INSERT INTO recipe_crafting_slots (recipe_id, slot_name, display_order)
VALUES (?, ?, ?)
ON CONFLICT(recipe_id, display_order) DO UPDATE SET
    slot_name = excluded.slot_name;
