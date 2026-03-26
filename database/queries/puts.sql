-- name: CreateMaterial :one
INSERT INTO materials (
    guid, wow_id, name
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: CreateRecipe :one
INSERT INTO recipes (
    guid, wow_id, name, profession
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: CreateCharacter :one
INSERT INTO characters (
    guid, name, server, guild, score, level
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: AddMaterialToRecipe :exec
INSERT INTO recipe_materials (
    recipe_id, material_id, quantity
) VALUES (
    ?, ?, ?
);

-- name: AddRecipeToCharacter :exec
INSERT INTO character_recipes (
    character_id, recipe_id
) VALUES (
    ?, ?
);
