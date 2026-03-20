-- RBAC 权限管理系统
-- 执行前确保已创建主数据库

USE supply_chain;

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    code VARCHAR(50) COMMENT '角色编码',
    description VARCHAR(255) COMMENT '角色描述',
    status TINYINT DEFAULT 1 COMMENT '状态: 1-启用, 0-禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '权限标识',
    type VARCHAR(20) DEFAULT 'api' COMMENT '权限类型: api, menu, button, data',
    description VARCHAR(255) COMMENT '权限描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 用户-角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 用户扩展信息（用于行级权限）
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id BIGINT UNSIGNED PRIMARY KEY,
    department VARCHAR(50) COMMENT '部门',
    position VARCHAR(50) COMMENT '职位',
    created_by BIGINT UNSIGNED COMMENT '创建者用户ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户扩展信息表';

-- 添加部门字段到 users 表（如果不存在）
-- ALTER TABLE users ADD COLUMN department VARCHAR(50) DEFAULT NULL COMMENT '部门' AFTER role;

-- 添加 created_by 和 department 字段到 suppliers 表（如果不存在）
-- ALTER TABLE suppliers ADD COLUMN department VARCHAR(50) DEFAULT NULL COMMENT '部门' AFTER remark;
-- ALTER TABLE suppliers ADD COLUMN created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID';

-- 添加 created_by 和 department 字段到 procurement_orders 表（如果不存在）
-- ALTER TABLE procurement_orders ADD COLUMN department VARCHAR(50) DEFAULT NULL COMMENT '部门' AFTER attachment_url;
-- ALTER TABLE procurement_orders ADD COLUMN created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID';

-- 添加 created_by 和 department 字段到 sales_orders 表（如果不存在）
-- ALTER TABLE sales_orders ADD COLUMN department VARCHAR(50) DEFAULT NULL COMMENT '部门' AFTER remark;
-- ALTER TABLE sales_orders ADD COLUMN created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID';

-- 添加 created_by 和 department 字段到 logistics_orders 表（如果不存在）
-- ALTER TABLE logistics_orders ADD COLUMN department VARCHAR(50) DEFAULT NULL COMMENT '部门' AFTER actual_delivery;
-- ALTER TABLE logistics_orders ADD COLUMN created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID';

-- 插入预定义角色
INSERT INTO roles (name, code, description, status) VALUES
('admin', 'admin', '管理员，拥有所有权限', 1),
('manager', 'manager', '经理，管理本部门数据', 1),
('operator', 'operator', '操作员，只能看自己创建的数据', 1)
ON DUPLICATE KEY UPDATE description = VALUES(description), code = VALUES(code);

-- 插入预定义权限（示例）
INSERT INTO permissions (name, type, description) VALUES
('supplier:read', 'api', '查看供应商'),
('supplier:write', 'api', '管理供应商'),
('inventory:read', 'api', '查看库存'),
('inventory:write', 'api', '管理库存'),
('procurement:read', 'api', '查看采购订单'),
('procurement:write', 'api', '管理采购订单'),
('sales:read', 'api', '查看销售订单'),
('sales:write', 'api', '管理销售订单'),
('logistics:read', 'api', '查看物流'),
('logistics:write', 'api', '管理物流'),
('user:read', 'api', '查看用户'),
('user:write', 'api', '管理用户'),
('role:read', 'api', '查看角色'),
('role:write', 'api', '管理角色'),
('permission:read', 'api', '查看权限'),
('permission:write', 'api', '管理权限'),
('dashboard:read', 'api', '查看仪表盘')
ON DUPLICATE KEY UPDATE description = VALUES(description), type = VALUES(type);

-- 给 admin 角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 给 manager 角色分配部分权限（示例）
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE name IN (
    'supplier:read', 'supplier:write',
    'inventory:read', 'inventory:write',
    'procurement:read', 'procurement:write',
    'sales:read', 'sales:write',
    'logistics:read', 'logistics:write'
)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 给 operator 角色分配部分权限（示例）
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE name IN (
    'supplier:read',
    'inventory:read',
    'procurement:read',
    'sales:read',
    'logistics:read',
    'dashboard:read'
)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 给用户分配角色（假设 admin 用户 id=1, manager 用户 id=2, operator 用户 id=3）
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1), (2, 2)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 用户-权限直接关联表（不通过角色）
CREATE TABLE IF NOT EXISTS user_permissions (
    user_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限直接关联表';