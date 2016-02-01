UPDATE nqm_agent # 台湾中华电信 TAIWANZHONGHUADIANXIN 台湾 TAIWAN 台湾 TAIWAN
SET ag_connection_id='nqm-agent-1@10.20.30.40', ag_hostname='nqm-agent-1', ag_ip_address=0x0A141E28, ag_isp_id=23, ag_pv_id=34, ag_ct_id=291
WHERE ag_id=1;

UPDATE nqm_agent # 广州e家宽 GUANGZHOUeJIAKUAN 广东 GUANGDONG 广州市 GUANGZHOUSHI
SET ag_connection_id='nqm-agent-2@10.20.30.40', ag_hostname='nqm-agent-2', ag_ip_address=0x0A141E28, ag_isp_id=12, ag_pv_id=20, ag_ct_id=5
WHERE ag_id=2;

INSERT INTO nqm_target(
	tg_id, tg_name, tg_host,
	tg_isp_id, tg_pv_id, tg_ct_id, tg_probed_by_all, tg_name_tag
)
VALUES
	(1, 'tgn-1', '203.208.150.145', 1, 17, 26, true, 'Room-1'), # 北京三信时代 BEIJINGSANXINSHIDAI , 湖南 HUNAN, 长沙市 [ZHANG CHANG]SHASHI
	(2, 'tgn-2', '203.208.146.33', 1, 17, 26, true, 'Room-1'), # 北京三信时代 BEIJINGSANXINSHIDAI , 湖南 HUNAN, 长沙市 [ZHANG CHANG]SHASHI
	(3, 'tgn-3', '203.208.232.53', 1, 17, 26, true, 'Room-2'), # 北京三信时代 BEIJINGSANXINSHIDAI , 湖南 HUNAN, 长沙市 [ZHANG CHANG]SHASHI
	(4, 'tgn-4', '211.160.177.145', 2, 17, 27, true, 'Room-1'), # 教育网 JIAOYUWANG, 湖南 HUNAN, 岳阳市 YUEYANGSHI
	(5, 'tgn-5', '211.160.177.150', 2, 17, 27, true, 'Room-2'), # 教育网 JIAOYUWANG, 湖南 HUNAN, 岳阳市 YUEYANGSHI
	(6, 'tgn-6', '211.160.177.70', 2, 17, 27, true, 'Room-2'), # 教育网 JIAOYUWANG, 湖南 HUNAN, 岳阳市 YUEYANGSHI
	(7, 'tgn-7', '61.190.194.110', 1, 20, 10, true, 'Room-1'), # 北京三信时代 BEIJINGSANXINSHIDAI 安徽 ANHUI 佛山市 FUSHANSHI
	(8, 'tgn-8', '61.190.194.114', 1, 20, 10, true, 'Room-2'), # 北京三信时代 BEIJINGSANXINSHIDAI 安徽 ANHUI 佛山市 FUSHANSHI
	(9, 'tgn-9', '61.190.194.222', 1, 20, 10, true, 'Room-2'), # 北京三信时代 BEIJINGSANXINSHIDAI 安徽 ANHUI 佛山市 FUSHANSHI
	(10, 'tgn-10', '219.72.0.82', 3, 17, 28, false, 'Room-0'), # 移动 YIDONG 湖南 HUNAN 张家界市 ZHANGJIAJIESHI
	(11, 'tgn-11', '219.72.224.69', 3, 17, 28, false, 'Room-0'), # 移动 YIDONG 湖南 HUNAN 张家界市 ZHANGJIAJIESHI
	(12, 'tgn-12', '219.72.67.2', 3, 17, 28, false, 'Room-0'), # 移动 YIDONG 湖南 HUNAN 张家界市 ZHANGJIAJIESHI
	(13, 'tgn-13', '221.122.0.18', 4, 20, 13, false, 'Room-0'), # 铁通 TIETONG 安徽 ANHUI 梅州市 MEIZHOUSHI
	(14, 'tgn-14', '221.122.0.22', 4, 20, 13, false, 'Room-0'), # 铁通 TIETONG 安徽 ANHUI 梅州市 MEIZHOUSHI
	(15, 'tgn-15', '221.122.0.78', 4, 20, 13, false, 'Room-0') # 铁通 TIETONG 安徽 ANHUI 梅州市 MEIZHOUSHI
ON DUPLICATE KEY UPDATE
	tg_name = VALUES(tg_name),
	tg_host = VALUES(tg_host),
	tg_isp_id = VALUES(tg_isp_id),
	tg_pv_id = VALUES(tg_pv_id),
	tg_ct_id = VALUES(tg_ct_id),
	tg_probed_by_all = VALUES(tg_probed_by_all),
	tg_name_tag = VALUES(tg_name_tag);

INSERT INTO nqm_ping_task(pt_ag_id, pt_period)
VALUES(1, 1)
ON DUPLICATE KEY UPDATE
	pt_period = VALUES(pt_period);

INSERT INTO nqm_ping_task(pt_ag_id, pt_period)
VALUES(2, 1)
ON DUPLICATE KEY UPDATE
	pt_period = VALUES(pt_period);

INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_ag_id, tfisp_isp_id) # 移动 YIDONG
VALUES(1, 3)
ON DUPLICATE KEY UPDATE
	tfisp_isp_id = VALUES(tfisp_isp_id);

INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_ag_id, tfisp_isp_id) # 铁通 TIETONG
VALUES(2, 4)
ON DUPLICATE KEY UPDATE
	tfisp_isp_id = VALUES(tfisp_isp_id);