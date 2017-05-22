ALTER TABLE dashboard_graph
  ADD COLUMN creator varchar(50) DEFAULT 'root';

ALTER TABLE dashboard_screen
  ADD COLUMN creator varchar(50) DEFAULT 'root';
