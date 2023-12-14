/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `block_uncles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `parent_hash` varchar(255) DEFAULT '',
  `uncle_hash` varchar(255) DEFAULT '',
  `coinbase` varchar(255) DEFAULT '',
  `block_root` varchar(255) DEFAULT '',
  `tx_hash` varchar(255) DEFAULT '',
  `receipt_hash` varchar(255) DEFAULT '',
  `bloom` blob DEFAULT '',
  `difficulty` bigint(20) unsigned DEFAULT 0,
  `block_number` bigint(20) unsigned DEFAULT 0,
  `gas_limit` bigint(20) unsigned DEFAULT 0,
  `gas_used` bigint(20) unsigned DEFAULT 0,
  `block_time` bigint(20) unsigned DEFAULT 0,
  `extra` blob DEFAULT '',
  `mix_digest` varchar(255) DEFAULT '',
  `block_nonce` bigint(20) unsigned DEFAULT 0,
  `base_fee` bigint(20) unsigned DEFAULT 0,
  `withdrawals_hash` varchar(255) DEFAULT '',
  `blob_gas_used` bigint(20) unsigned DEFAULT 0,
  `excess_blob_gas` bigint(20) unsigned DEFAULT 0,
  `parent_beacon_root` varchar(255) DEFAULT '',
  `block_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_block_uncles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blocks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `parent_hash` varchar(255) DEFAULT '',
  `uncle_hash` varchar(255) DEFAULT '',
  `coinbase` varchar(255) DEFAULT '',
  `block_root` varchar(255) DEFAULT '',
  `tx_hash` varchar(255) DEFAULT '',
  `receipt_hash` varchar(255) DEFAULT '',
  `bloom` blob DEFAULT '',
  `difficulty` bigint(20) unsigned DEFAULT 0,
  `block_number` bigint(20) unsigned DEFAULT 0,
  `gas_limit` bigint(20) unsigned DEFAULT 0,
  `gas_used` bigint(20) unsigned DEFAULT 0,
  `block_time` bigint(20) unsigned DEFAULT 0,
  `extra` blob DEFAULT '',
  `mix_digest` varchar(255) DEFAULT '',
  `block_nonce` bigint(20) unsigned DEFAULT 0,
  `base_fee` bigint(20) unsigned DEFAULT 0,
  `withdrawals_hash` varchar(255) DEFAULT '',
  `blob_gas_used` bigint(20) unsigned DEFAULT 0,
  `excess_blob_gas` bigint(20) unsigned DEFAULT 0,
  `parent_beacon_root` varchar(255) DEFAULT '',
  `block_hash` varchar(255) DEFAULT '',
  `block_size` bigint(20) unsigned DEFAULT 0,
  `received_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `idx_blocks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_log_topics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `topic` varchar(255) DEFAULT '',
  `block_id` int(10) unsigned DEFAULT 0,
  `transaction_id` int(10) unsigned DEFAULT 0,
  `transaction_receipt_id` int(10) unsigned DEFAULT 0,
  `transaction_log_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_log_topics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=175 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_logs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `address` varchar(255) DEFAULT '',
  `log_data` blob DEFAULT '',
  `block_number` bigint(20) unsigned DEFAULT 0,
  `tx_hash` varchar(255) DEFAULT '',
  `tx_index` int(10) unsigned DEFAULT 0,
  `block_hash` varchar(255) DEFAULT '',
  `log_index` int(10) unsigned DEFAULT 0,
  `removed` tinyint(1) DEFAULT 0,
  `block_id` int(10) unsigned DEFAULT 0,
  `transaction_id` int(10) unsigned DEFAULT 0,
  `transaction_receipt_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_receipts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `receipt_type` tinyint(3) unsigned DEFAULT 0,
  `post_state` blob DEFAULT '',
  `tx_status` bigint(20) unsigned DEFAULT 0,
  `cumulative_gas_used` bigint(20) unsigned DEFAULT 0,
  `bloom` blob DEFAULT '',
  `tx_hash` varchar(255) DEFAULT '',
  `contract_address` varchar(255) DEFAULT '',
  `gas_used` bigint(20) unsigned DEFAULT 0,
  `effective_gas_price` bigint(20) unsigned DEFAULT 0,
  `blob_gas_used` bigint(20) unsigned DEFAULT 0,
  `blob_gas_price` bigint(20) unsigned DEFAULT 0,
  `block_hash` varchar(255) DEFAULT '',
  `block_number` bigint(20) unsigned DEFAULT 0,
  `transaction_index` int(10) unsigned DEFAULT 0,
  `block_id` int(10) unsigned DEFAULT 0,
  `transaction_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_receipts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `tx_type` tinyint(3) unsigned DEFAULT 0,
  `chain_id` bigint(20) unsigned DEFAULT 0,
  `tx_data` blob DEFAULT '',
  `gas` bigint(20) unsigned DEFAULT 0,
  `gas_price` bigint(20) unsigned DEFAULT 0,
  `gas_tip_cap` bigint(20) unsigned DEFAULT 0,
  `gas_fee_cap` bigint(20) unsigned DEFAULT 0,
  `tx_value` bigint(20) unsigned DEFAULT 0,
  `account_nonce` bigint(20) unsigned DEFAULT 0,
  `tx_to` varchar(255) DEFAULT '',
  `tx_v` bigint(20) unsigned DEFAULT 0,
  `tx_r` bigint(20) unsigned DEFAULT 0,
  `tx_s` bigint(20) unsigned DEFAULT 0,
  `block_number` bigint(20) unsigned DEFAULT 0,
  `block_hash` varchar(255) DEFAULT '',
  `block_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
