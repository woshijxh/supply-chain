-- =====================================================
-- 供应链管理系统 - 权限模拟数据 (完整版)
-- 生成时间: 2026-03-24
-- 执行方式: mysql -u root -p supply_chain < permission_data.sql
-- =====================================================

USE supply_chain;

-- =====================================================
-- 1. 清理旧数据
-- =====================================================
SET FOREIGN_KEY_CHECKS = 0;
DELETE FROM role_permissions;
DELETE FROM user_roles;
DELETE FROM user_permissions;
-- Casbin 规则由 Casbin 自动管理，如有 casbin_rule 表可取消注释
-- DELETE FROM casbin_rule;
DELETE FROM permissions;
DELETE FROM roles WHERE id > 0;
DELETE FROM user_profiles;
DELETE FROM users WHERE id > 0;
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 2. 创建权限 (根据路由中的 RBAC 权限定义)
-- =====================================================
-- 权限类型说明:
-- - menu: 菜单权限，控制侧边栏菜单显示
-- - button: 按钮权限，控制操作按钮显示
-- 权限操作:
-- - read: 查看列表、详情、导出等只读操作
-- - create: 新增操作
-- - update: 编辑、审批等修改操作
-- - delete: 删除操作

INSERT INTO permissions (id, name, code, type, description, status, created_at) VALUES
-- ========== 仪表盘 (id: 1) ==========
(1, '仪表盘查看', 'dashboard:read', 'menu', '查看仪表盘统计数据、图表分析', 1, NOW()),

-- ========== 供应商管理 (id: 2-5) ==========
(2, '供应商查看', 'supplier:read', 'menu', '查看供应商列表和详情', 1, NOW()),
(3, '供应商新增', 'supplier:create', 'button', '新增供应商', 1, NOW()),
(4, '供应商编辑', 'supplier:update', 'button', '编辑供应商信息', 1, NOW()),
(5, '供应商删除', 'supplier:delete', 'button', '删除供应商', 1, NOW()),

-- ========== 产品管理 (id: 6-9) ==========
(6, '产品查看', 'product:read', 'menu', '查看产品列表和详情', 1, NOW()),
(7, '产品新增', 'product:create', 'button', '新增产品', 1, NOW()),
(8, '产品编辑', 'product:update', 'button', '编辑产品信息', 1, NOW()),
(9, '产品删除', 'product:delete', 'button', '删除产品', 1, NOW()),

-- ========== 库存管理 (id: 10-13) ==========
(10, '库存查看', 'inventory:read', 'menu', '查看库存列表、统计、流水记录', 1, NOW()),
(11, '库存入库', 'inventory:create', 'button', '库存入库操作', 1, NOW()),
(12, '库存出库', 'inventory:update', 'button', '库存出库操作', 1, NOW()),
(13, '库存调整', 'inventory:delete', 'button', '库存盘点调整', 1, NOW()),

-- ========== 客户管理 (id: 14-17) ==========
(14, '客户查看', 'customer:read', 'menu', '查看客户列表和详情', 1, NOW()),
(15, '客户新增', 'customer:create', 'button', '新增客户', 1, NOW()),
(16, '客户编辑', 'customer:update', 'button', '编辑客户信息', 1, NOW()),
(17, '客户删除', 'customer:delete', 'button', '删除客户', 1, NOW()),

-- ========== 采购管理 (id: 18-21) ==========
(18, '采购查看', 'procurement:read', 'menu', '查看采购订单、采购退货列表和详情', 1, NOW()),
(19, '采购新增', 'procurement:create', 'button', '新增采购订单、采购退货', 1, NOW()),
(20, '采购编辑', 'procurement:update', 'button', '编辑、审批采购订单，采购退货操作', 1, NOW()),
(21, '采购删除', 'procurement:delete', 'button', '删除采购订单、采购退货', 1, NOW()),

-- ========== 销售管理 (id: 22-25) ==========
(22, '销售查看', 'sales:read', 'menu', '查看销售订单、销售退货列表和详情', 1, NOW()),
(23, '销售新增', 'sales:create', 'button', '新增销售订单、销售退货', 1, NOW()),
(24, '销售编辑', 'sales:update', 'button', '确认、取消销售订单，销售退货操作', 1, NOW()),
(25, '销售删除', 'sales:delete', 'button', '删除销售订单、销售退货', 1, NOW()),

-- ========== 物流管理 (id: 26-29) ==========
(26, '物流查看', 'logistics:read', 'menu', '查看物流订单列表、详情、轨迹', 1, NOW()),
(27, '物流新增', 'logistics:create', 'button', '新增物流订单', 1, NOW()),
(28, '物流编辑', 'logistics:update', 'button', '更新物流状态、轨迹', 1, NOW()),
(29, '物流删除', 'logistics:delete', 'button', '删除物流订单', 1, NOW()),

-- ========== 用户管理 (id: 30-33) ==========
(30, '用户查看', 'user:read', 'menu', '查看用户列表和详情', 1, NOW()),
(31, '用户新增', 'user:create', 'button', '新增用户', 1, NOW()),
(32, '用户编辑', 'user:update', 'button', '编辑用户、重置密码、分配角色权限', 1, NOW()),
(33, '用户删除', 'user:delete', 'button', '删除用户', 1, NOW()),

-- ========== 角色管理 (id: 34-37) ==========
(34, '角色查看', 'role:read', 'menu', '查看角色列表和详情', 1, NOW()),
(35, '角色新增', 'role:create', 'button', '新增角色', 1, NOW()),
(36, '角色编辑', 'role:update', 'button', '编辑角色、分配权限', 1, NOW()),
(37, '角色删除', 'role:delete', 'button', '删除角色', 1, NOW()),

-- ========== 权限管理 (id: 38-41) ==========
(38, '权限查看', 'permission:read', 'menu', '查看权限列表和详情', 1, NOW()),
(39, '权限新增', 'permission:create', 'button', '新增权限', 1, NOW()),
(40, '权限编辑', 'permission:update', 'button', '编辑权限', 1, NOW()),
(41, '权限删除', 'permission:delete', 'button', '删除权限', 1, NOW()),

-- ========== 商品追溯 (id: 42) ==========
(42, '商品追溯', 'trace:read', 'menu', '追溯商品流转全链路', 1, NOW());

-- =====================================================
-- 3. 创建角色
-- =====================================================
INSERT INTO roles (id, name, code, description, status, created_at) VALUES
(1, '超级管理员', 'admin',     '拥有系统所有权限，可管理用户、角色、权限和所有业务模块', 1, NOW()),
(2, '业务经理',   'manager',   '可管理供应商、产品、库存、客户、采购、销售、物流等业务模块', 1, NOW()),
(3, '业务员',     'operator',  '只能查看各模块数据，无编辑删除权限', 1, NOW()),
(4, '采购专员',   'buyer',     '负责采购模块的查看和操作，以及相关供应商和库存查看', 1, NOW()),
(5, '销售专员',   'seller',    '负责销售、物流模块的查看和操作，以及相关客户和库存查看', 1, NOW()),
(6, '仓库管理员', 'warehouse', '负责库存模块的查看和操作，以及采购和销售的查看', 1, NOW()),
(7, '财务专员',   'finance',   '负责查看所有订单数据，审批相关操作', 1, NOW());

-- =====================================================
-- 4. 创建用户 (密码统一为: 123456)
-- =====================================================
INSERT INTO users (id, username, password, email, phone, role, department, position, status, created_at, updated_at) VALUES
(1, 'admin',     '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'admin@supply.com',     '13800000001', 'admin',     '管理层', '系统管理员',     1, NOW(), NOW()),
(2, 'manager',   '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'manager@supply.com',   '13800000002', 'manager',   '业务部', '业务经理',       1, NOW(), NOW()),
(3, 'operator',  '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'operator@supply.com',  '13800000003', 'operator',  '业务部', '业务员',         1, NOW(), NOW()),
(4, 'buyer01',   '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'buyer01@supply.com',   '13800000004', 'buyer',     '采购部', '采购专员',       1, NOW(), NOW()),
(5, 'seller01',  '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'seller01@supply.com',  '13800000005', 'seller',    '销售部', '销售专员',       1, NOW(), NOW()),
(6, 'warehouse', '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'warehouse@supply.com', '13800000006', 'warehouse', '仓储部', '仓库管理员',     1, NOW(), NOW()),
(7, 'finance',   '$2a$10$vRXaknqQ0/esJRg2liVHOOKMC5kcDKaaSmCVCJ6FwRfb0KuS8s8Ui', 'finance@supply.com',   '13800000007', 'finance',   '财务部', '财务专员',       1, NOW(), NOW());

-- =====================================================
-- 5. 用户扩展信息
-- =====================================================
INSERT INTO user_profiles (user_id, department, position, created_at, updated_at) VALUES
(1, '管理层', '系统管理员', NOW(), NOW()),
(2, '业务部', '业务经理', NOW(), NOW()),
(3, '业务部', '业务员', NOW(), NOW()),
(4, '采购部', '采购专员', NOW(), NOW()),
(5, '销售部', '销售专员', NOW(), NOW()),
(6, '仓储部', '仓库管理员', NOW(), NOW()),
(7, '财务部', '财务专员', NOW(), NOW());

-- =====================================================
-- 6. 用户角色关联
-- =====================================================
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),  -- admin -> 超级管理员
(2, 2),  -- manager -> 业务经理
(3, 3),  -- operator -> 业务员
(4, 4),  -- buyer01 -> 采购专员
(5, 5),  -- seller01 -> 销售专员
(6, 6),  -- warehouse -> 仓库管理员
(7, 7);  -- finance -> 财务专员

-- =====================================================
-- 7. 角色权限关联
-- =====================================================

-- ========== 超级管理员 (admin) - 拥有所有权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- 仪表盘
(1, 1),
-- 供应商
(1, 2), (1, 3), (1, 4), (1, 5),
-- 产品
(1, 6), (1, 7), (1, 8), (1, 9),
-- 库存
(1, 10), (1, 11), (1, 12), (1, 13),
-- 客户
(1, 14), (1, 15), (1, 16), (1, 17),
-- 采购
(1, 18), (1, 19), (1, 20), (1, 21),
-- 销售
(1, 22), (1, 23), (1, 24), (1, 25),
-- 物流
(1, 26), (1, 27), (1, 28), (1, 29),
-- 用户
(1, 30), (1, 31), (1, 32), (1, 33),
-- 角色
(1, 34), (1, 35), (1, 36), (1, 37),
-- 权限
(1, 38), (1, 39), (1, 40), (1, 41),
-- 商品追溯
(1, 42);

-- ========== 业务经理 (manager) - 业务模块权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- 仪表盘
(2, 1),
-- 供应商
(2, 2), (2, 3), (2, 4), (2, 5),
-- 产品
(2, 6), (2, 7), (2, 8), (2, 9),
-- 库存
(2, 10), (2, 11), (2, 12), (2, 13),
-- 客户
(2, 14), (2, 15), (2, 16), (2, 17),
-- 采购
(2, 18), (2, 19), (2, 20), (2, 21),
-- 销售
(2, 22), (2, 23), (2, 24), (2, 25),
-- 物流
(2, 26), (2, 27), (2, 28), (2, 29),
-- 商品追溯
(2, 42);

-- ========== 业务员 (operator) - 只读权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
(3, 1),   -- 仪表盘
(3, 2),   -- 供应商查看
(3, 6),   -- 产品查看
(3, 10),  -- 库存查看
(3, 14),  -- 客户查看
(3, 18),  -- 采购查看
(3, 22),  -- 销售查看
(3, 26),  -- 物流查看
(3, 42);  -- 商品追溯

-- ========== 采购专员 (buyer) - 采购相关权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
(4, 1),   -- 仪表盘
(4, 2),   -- 供应商查看
(4, 3),   -- 供应商新增
(4, 4),   -- 供应商编辑
(4, 6),   -- 产品查看
(4, 10),  -- 库存查看
(4, 11),  -- 库存入库
-- 采购（全部操作）
(4, 18), (4, 19), (4, 20), (4, 21),
-- 商品追溯
(4, 42);

-- ========== 销售专员 (seller) - 销售相关权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
(5, 1),   -- 仪表盘
(5, 6),   -- 产品查看
(5, 10),  -- 库存查看
-- 客户（全部操作）
(5, 14), (5, 15), (5, 16), (5, 17),
-- 销售（全部操作）
(5, 22), (5, 23), (5, 24), (5, 25),
-- 物流（全部操作）
(5, 26), (5, 27), (5, 28), (5, 29),
-- 商品追溯
(5, 42);

-- ========== 仓库管理员 (warehouse) - 库存相关权限 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
(6, 1),   -- 仪表盘
(6, 6),   -- 产品查看
-- 库存（全部操作）
(6, 10), (6, 11), (6, 12), (6, 13),
(6, 18),  -- 采购查看
(6, 22),  -- 销售查看
(6, 26),  -- 物流查看
-- 商品追溯
(6, 42);

-- ========== 财务专员 (finance) - 查看所有订单 ==========
INSERT INTO role_permissions (role_id, permission_id) VALUES
(7, 1),   -- 仪表盘
(7, 2),   -- 供应商查看
(7, 6),   -- 产品查看
(7, 10),  -- 库存查看
(7, 14),  -- 客户查看
(7, 18),  -- 采购查看
(7, 20),  -- 采购编辑（审批）
(7, 22),  -- 销售查看
(7, 24),  -- 销售编辑（审批）
(7, 26),  -- 物流查看
(7, 42);  -- 商品追溯

-- =====================================================
-- 8. 初始化 Casbin 规则 (RBAC 策略)
-- 注意: Casbin 规则通常由应用启动时自动初始化
-- 如需手动初始化，请确保 casbin_rule 表已存在
-- =====================================================
-- 以下 SQL 仅在 casbin_rule 表存在时执行
-- INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ...

-- =====================================================
-- 完成
-- =====================================================
SELECT '======================================' AS '';
SELECT '权限数据初始化完成！' AS message;
SELECT '======================================' AS '';

SELECT '权限统计' AS type, COUNT(*) AS count FROM permissions
UNION ALL SELECT '角色统计', COUNT(*) FROM roles
UNION ALL SELECT '用户统计', COUNT(*) FROM users
UNION ALL SELECT '角色权限关联', COUNT(*) FROM role_permissions
UNION ALL SELECT '用户角色关联', COUNT(*) FROM user_roles;

SELECT '======================================' AS '';
SELECT '用户账号列表 (密码统一为: 123456)' AS '';
SELECT '======================================' AS '';
SELECT id, username, role, department, position FROM users ORDER BY id;

SELECT '======================================' AS '';
SELECT '角色权限统计' AS '';
SELECT '======================================' AS '';
SELECT r.name AS role_name, COUNT(rp.permission_id) AS permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.id, r.name
ORDER BY r.id;