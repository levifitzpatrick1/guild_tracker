-- name: SearchRecipes :many
SELECT * FROM recipes
WHERE name = ?
ORDER BY name ASC;

-- name: GetCharactersForRecipe :many
SELECT c.* FROM characters c
JOIN character_recipes cr ON c.id = cr.character_id
WHERE cr.recipe_id = ?
ORDER BY c.name ASC;

-- name: GetMaterialsForRecipe :many
SELECT
    m.id,
    m.name,
    rm.quantity
FROM materials m
JOIN recipe_materials rm ON m.id = rm.material_id
WHERE rm.recipe_id = ?
ORDER BY m.name ASC;

-- name: GetRecipeByID :one
SELECT
    *
FROM recipes
WHERE id = ?;

-- name: GetRecipesForCharacter :many
SELECT
    r.*
FROM recipes r
JOIN character_recipes cr ON r.id = cr.recipe_id
WHERE cr.character_id = ?
ORDER BY r.name ASC;
