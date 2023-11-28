CREATE TABLE `ips` (
   `proxy_id` bigint NOT NULL AUTO_INCREMENT,
   `proxy_host` varchar(255) NOT NULL UNIQUE COMMENT 'ip 地址',
   `proxy_port` int(11) NOT NULL COMMENT 'ip 端口',
   `proxy_type` varchar(64) NOT NULL COMMENT '代理类型 http or https',
   `proxy_location` varchar(255) COMMENT 'ip位于地址',
   `proxy_speed` int(20) NOT NULL COMMENT '测试连通速度',
   `proxy_source` varchar(64) NOT NULL COMMENT '代理源',
   `create_time` datetime NOT NULL default CURRENT_TIMESTAMP COMMENT '创建时间',
   `update_time` datetime NOT NULL default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次修改时间',
   PRIMARY KEY (`proxy_id`)
)ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8mb4_general_ci;