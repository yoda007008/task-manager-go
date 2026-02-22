-- чтобы находить быстро выполненные и невыполненные задачи, вешаем индекс на колонку completed
CREATE INDEX IF NOT EXISTS completed_tasks_idx ON tasks(completed);
