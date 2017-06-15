ALTER TABLE dashboard_graph
  ADD COLUMN `time_range` varchar(50) DEFAULT '3h',
  ADD COLUMN `y_scale` varchar(50) DEFAULT NULL,
  ADD COLUMN `sample_method` varchar(20) DEFAULT 'AVERAGE',
  ADD COLUMN `sort_by` varchar(30) DEFAULT 'a-z';

