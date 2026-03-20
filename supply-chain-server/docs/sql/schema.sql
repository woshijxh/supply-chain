-- 供应链系统数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS supply_chain DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE supply_chain;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    role VARCHAR(20) DEFAULT 'operator' COMMENT '角色: admin, manager, operator',
    status TINYINT DEFAULT 1 COMMENT '状态: 1-启用, 0-禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入默认管理员
INSERT INTO users (username, password, email, role, status) VALUES
('admin', 'admin123', 'admin@supply.com', 'admin', 1);

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE COMMENT '供应商编码',
    name VARCHAR(100) NOT NULL COMMENT '供应商名称',
    contact VARCHAR(50) COMMENT '联系人',
    phone VARCHAR(20) COMMENT '电话',
    email VARCHAR(100) COMMENT '邮箱',
    address VARCHAR(255) COMMENT '地址',
    level VARCHAR(1) DEFAULT 'B' COMMENT '等级: A, B, C',
    category VARCHAR(50) COMMENT '品类',
    payment_terms VARCHAR(100) COMMENT '付款条款',
    bank_name VARCHAR(100) COMMENT '开户银行',
    bank_account VARCHAR(50) COMMENT '银行账号',
    tax_number VARCHAR(50) COMMENT '税号',
    rating DECIMAL(2,1) DEFAULT 4.0 COMMENT '评分',
    status TINYINT DEFAULT 1 COMMENT '状态',
    remark TEXT COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='供应商表';

-- 产品表
CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE COMMENT '产品编码',
    name VARCHAR(100) NOT NULL COMMENT '产品名称',
    sku VARCHAR(50) UNIQUE COMMENT 'SKU',
    category VARCHAR(50) COMMENT '分类',
    unit VARCHAR(20) COMMENT '单位',
    cost_price DECIMAL(12,2) COMMENT '成本价',
    sale_price DECIMAL(12,2) COMMENT '销售价',
    min_stock INT DEFAULT 0 COMMENT '最低库存',
    max_stock INT DEFAULT 0 COMMENT '最高库存',
    description TEXT COMMENT '描述',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品表';

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
    warehouse VARCHAR(50) NOT NULL COMMENT '仓库',
    quantity INT DEFAULT 0 COMMENT '库存数量',
    available_qty INT DEFAULT 0 COMMENT '可用数量',
    locked_qty INT DEFAULT 0 COMMENT '锁定数量',
    location VARCHAR(50) COMMENT '库位',
    batch_no VARCHAR(50) COMMENT '批次号',
    expiry_date DATE COMMENT '有效期',
    status VARCHAR(20) DEFAULT 'normal' COMMENT '状态: normal, low, over, locked',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_product_id (product_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 采购订单表
CREATE TABLE IF NOT EXISTS procurement_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(30) NOT NULL UNIQUE COMMENT '订单号',
    supplier_id BIGINT UNSIGNED NOT NULL COMMENT '供应商ID',
    total_amount DECIMAL(12,2) COMMENT '总金额',
    order_date DATETIME COMMENT '订单日期',
    expected_date DATE COMMENT '预计到货日期',
    actual_date DATE COMMENT '实际到货日期',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending, approved, purchasing, received, cancelled',
    warehouse VARCHAR(50) COMMENT '仓库',
    remark TEXT COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_supplier_id (supplier_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='采购订单表';

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS procurement_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    product_id BIGINT UNSIGNED COMMENT '产品ID',
    product_name VARCHAR(100) COMMENT '产品名称',
    quantity INT COMMENT '数量',
    unit VARCHAR(20) COMMENT '单位',
    unit_price DECIMAL(12,2) COMMENT '单价',
    amount DECIMAL(12,2) COMMENT '金额',
    received_qty INT DEFAULT 0 COMMENT '已收货数量',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id),
    FOREIGN KEY (order_id) REFERENCES procurement_orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='采购订单明细表';

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(30) NOT NULL UNIQUE COMMENT '订单号',
    customer_id BIGINT UNSIGNED COMMENT '客户ID',
    customer_name VARCHAR(100) NOT NULL COMMENT '客户名称',
    customer_phone VARCHAR(20) COMMENT '客户电话',
    customer_address VARCHAR(255) COMMENT '客户地址',
    total_amount DECIMAL(12,2) COMMENT '总金额',
    discount DECIMAL(12,2) DEFAULT 0 COMMENT '折扣',
    tax DECIMAL(12,2) DEFAULT 0 COMMENT '税费',
    shipping_fee DECIMAL(12,2) DEFAULT 0 COMMENT '运费',
    order_date DATETIME COMMENT '订单日期',
    delivery_date DATE COMMENT '发货日期',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending, confirmed, shipping, completed, cancelled, refunded',
    payment_method VARCHAR(20) COMMENT '支付方式',
    payment_status VARCHAR(20) DEFAULT 'pending' COMMENT '支付状态: pending, paid, refunded',
    remark TEXT COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='销售订单表';

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    product_id BIGINT UNSIGNED COMMENT '产品ID',
    product_name VARCHAR(100) COMMENT '产品名称',
    quantity INT COMMENT '数量',
    unit VARCHAR(20) COMMENT '单位',
    unit_price DECIMAL(12,2) COMMENT '单价',
    amount DECIMAL(12,2) COMMENT '金额',
    discount DECIMAL(12,2) DEFAULT 0 COMMENT '折扣',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id),
    FOREIGN KEY (order_id) REFERENCES sales_orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='销售订单明细表';

-- 物流订单表
CREATE TABLE IF NOT EXISTS logistics_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tracking_no VARCHAR(30) NOT NULL UNIQUE COMMENT '运单号',
    carrier VARCHAR(50) COMMENT '承运商',
    sales_order_id BIGINT UNSIGNED COMMENT '销售订单ID',
    sales_order_no VARCHAR(30) COMMENT '销售订单号',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态: pending, picked, in_transit, delivering, delivered, returned',
    sender_name VARCHAR(50) COMMENT '发件人',
    sender_phone VARCHAR(20) COMMENT '发件人电话',
    sender_address VARCHAR(255) COMMENT '发件地址',
    receiver_name VARCHAR(50) COMMENT '收件人',
    receiver_phone VARCHAR(20) COMMENT '收件人电话',
    receiver_address VARCHAR(255) COMMENT '收件地址',
    weight DECIMAL(10,2) COMMENT '重量(kg)',
    shipping_fee DECIMAL(10,2) COMMENT '运费',
    estimated_delivery DATE COMMENT '预计送达',
    actual_delivery DATE COMMENT '实际送达',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_sales_order_id (sales_order_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流订单表';

-- 物流轨迹表
CREATE TABLE IF NOT EXISTS logistics_timelines (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    logistics_id BIGINT UNSIGNED NOT NULL COMMENT '物流订单ID',
    time DATETIME COMMENT '时间',
    status VARCHAR(50) COMMENT '状态',
    location VARCHAR(100) COMMENT '地点',
    description VARCHAR(255) COMMENT '描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_logistics_id (logistics_id),
    FOREIGN KEY (logistics_id) REFERENCES logistics_orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流轨迹表';