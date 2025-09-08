本项目选用postgres而非mysql
1.postgres支持原生JSONB,可索引JSON，高性能
  适合本项目代码提交结果、战报分析、技能卡效果这类结构不固定的数据

2.MVCC(多版本并发控制)事务隔离级别更严格，避免并发写问题，适合高并发
  Mysql要调锁策略

3.postgres查询接近oracle，支持窗口函数、CTE(with子句)、全文检索、地理数据

4.postgres扩展强，严格遵循SQL标准，数据一致性更可靠