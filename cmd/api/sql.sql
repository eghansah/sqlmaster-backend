CREATE TABLE
  `sql_query` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `updated_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `title` varchar(255) NOT NULL,
    `description` longtext NOT NULL,
    `query` longtext NOT NULL,
    `datasource`  int(10) unsigned,
    `enabled` tinyint(1) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq_title` (`title`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;



  CREATE TABLE
  `sql_query_param` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `datatype` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `label` varchar(255) NOT NULL,
    `multi_value` tinyint(1) NOT NULL DEFAULT 0,
    `sql_query` int(10) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unq1` (`sql_query` DESC, `name` DESC),
    CONSTRAINT `sql_query_param_relation_1` FOREIGN KEY (`sql_query`) REFERENCES `sql_query` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;


  CREATE TABLE
  `report_request` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `updated_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `sql_query_id` int(10) unsigned NOT NULL,
    `query` longtext NOT NULL,
    `params` longtext NOT NULL,
    `status` varchar(255) NOT NULL DEFAULT 'new',
    `comment` longtext NOT NULL DEFAULT ' ',
    PRIMARY KEY (`id`),
    KEY `report_request_relation_1` (`sql_query_id`),
    CONSTRAINT `report_request_relation_1` FOREIGN KEY (`sql_query_id`) REFERENCES `sql_query` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;




  CREATE TABLE
  `datasources` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
    `type` varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    `dsn` longtext NOT NULL,
    `enabled` tinyint(1) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `datasources_index_2` (`name`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

  