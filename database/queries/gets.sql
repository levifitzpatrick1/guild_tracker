-- name: SearchRecipes :many
SELECT * FROM recipes
WHERE name = ?
ORDER BY name ASC;

-- name: GetCharactersForRecipe :many
SELECT c.* FROM characters c
JOIN character_recipes cr ON c.guid = cr.character_id
WHERE cr.recipe_id = ?
ORDER BY c.name ASC;

-- name: GetMaterialsForRecipe :many
SELECT
    m.guid,
    m.wow_id,
    m.name,
    rm.quantity
FROM materials m
JOIN recipe_materials rm ON m.guid = rm.material_id
WHERE rm.recipe_id = ?
ORDER BY m.name ASC;
