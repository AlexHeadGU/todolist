-- Откат: удаляем таблицу tasks (триггер и функция удалятся каскадно)
DROP TABLE IF EXISTS tasks;
DROP FUNCTION IF EXISTS update_updated_at_column();