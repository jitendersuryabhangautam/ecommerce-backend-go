# E-Commerce Backend - Architecture Diagrams

## System Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │  Web Client  │  │ Mobile App   │  │  Admin Panel │           │
│  │ (React/Vue)  │  │ (iOS/Android)│  │   (React)    │           │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘           │
│         │                  │                  │                   │
│         └──────────────────┴──────────────────┘                   │
│                            │                                      │
│                   HTTP/HTTPS (REST API)                           │
└────────────────────────────┼──────────────────────────────────────┘
                             │
┌────────────────────────────┼──────────────────────────────────────┐
│                   API GATEWAY (Port 8080)                         │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │            Gin Web Framework Router                        │   │
│  └───────────────────────────────────────────────────────────┘   │
└────────────────────────────┼──────────────────────────────────────┘
                             │
┌────────────────────────────┼──────────────────────────────────────┐
│                    MIDDLEWARE LAYER                               │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐            │
│  │   CORS   │→│ Request  │→│ Logging  │→│ Recovery │            │
│  │          │ │    ID    │ │          │ │          │            │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘            │
│  ┌──────────┐ ┌──────────┐                                       │
│  │   Auth   │→│  Admin   │  (Conditional)                        │
│  │   JWT    │ │  Check   │                                       │
│  └──────────┘ └──────────┘                                       │
└────────────────────────────┼──────────────────────────────────────┘
                             │
┌────────────────────────────┼──────────────────────────────────────┐
│                    HANDLER LAYER                                  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐            │
│  │   Auth   │ │ Products │ │   Cart   │ │  Orders  │            │
│  │ Handler  │ │ Handler  │ │ Handler  │ │ Handler  │            │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘            │
│       │            │            │            │                   │
│  ┌────┴─────┐ ┌────┴─────┐ ┌────┴─────┐                          │
│  │ Payment  │ │ Returns  │ │  Health  │                          │
│  │ Handler  │ │ Handler  │ │ Handler  │                          │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘                          │
│       │            │            │                                 │
│       │ Request Validation & Response Formatting                 │
└───────┼────────────┼────────────┼──────────────────────────────────┘
        │            │            │
┌───────┼────────────┼────────────┼──────────────────────────────────┐
│               SERVICE LAYER (Business Logic)                      │
│  ┌────┴─────┐ ┌────┴─────┐ ┌────┴─────┐ ┌──────────┐            │
│  │   Auth   │ │ Product  │ │   Cart   │ │  Order   │            │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │            │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘            │
│       │            │            │            │                   │
│  ┌────┴─────┐ ┌────┴─────┐                                       │
│  │ Payment  │ │ Return   │                                       │
│  │ Service  │ │ Service  │                                       │
│  └────┬─────┘ └────┬─────┘                                       │
│       │            │                                             │
│       │ Orchestration | Validation | Transaction Management      │
└───────┼────────────┼──────────────────────────────────────────────┘
        │            │
┌───────┼────────────┼──────────────────────────────────────────────┐
│            REPOSITORY LAYER (Data Access)                         │
│  ┌────┴─────┐ ┌────┴─────┐ ┌──────────┐ ┌──────────┐            │
│  │   User   │ │ Product  │ │   Cart   │ │  Order   │            │
│  │   Repo   │ │   Repo   │ │   Repo   │ │   Repo   │            │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘            │
│       │            │            │            │                   │
│  ┌────┴─────┐ ┌────┴─────┐                                       │
│  │ Payment  │ │ Return   │                                       │
│  │   Repo   │ │   Repo   │                                       │
│  └────┬─────┘ └────┬─────┘                                       │
│       │            │                                             │
│       │ SQL Queries | Connection Pool | Transaction Mgmt         │
└───────┼────────────┼──────────────────────────────────────────────┘
        │            │
        └────────────┴───────────────────┐
                                        │
┌────────────────────────────────────────┼──────────────────────────┐
│               DATABASE LAYER                                      │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │         PostgreSQL 15 (Port 5433)                          │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │   │
│  │  │  users   │ │ products │ │  carts   │ │  orders  │     │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘     │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐                  │   │
│  │  │ payments │ │ returns  │ │  stock_  │                  │   │
│  │  │          │ │          │ │reservatns│                  │   │
│  │  └──────────┘ └──────────┘ └──────────┘                  │   │
│  │                                                            │   │
│  │  Connection Pool | Transactions | Indexes | Constraints   │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                   │
│  Volume: postgres_data (Persistent Storage)                      │
└───────────────────────────────────────────────────────────────────┘
```

## Request Flow Diagram - Create Order

```
┌─────────┐
│ Client  │
└────┬────┘
     │ POST /api/v1/orders
     │ Authorization: Bearer <JWT>
     │ {shipping_address, billing_address, payment_method}
     ↓
┌────────────────────────────────────────────────────────────────┐
│ MIDDLEWARE CHAIN                                               │
├────────────────────────────────────────────────────────────────┤
│ 1. CORS          → Set CORS headers                            │
│ 2. RequestID     → Generate UUID                               │
│ 3. Logging       → Log request                                 │
│ 4. Recovery      → Panic handler                               │
│ 5. Auth          → Validate JWT, extract user_id & role        │
└────┬───────────────────────────────────────────────────────────┘
     ↓ c.Set("userID", ...) | c.Set("userRole", ...)
┌────────────────────────────────────────────────────────────────┐
│ HANDLER: OrderHandler.CreateOrder                              │
├────────────────────────────────────────────────────────────────┤
│ 1. Parse JSON body                                             │
│ 2. Validate input (validator library)                          │
│ 3. Extract userID from context                                 │
│ 4. Call OrderService.CreateOrder(userID, request)              │
└────┬───────────────────────────────────────────────────────────┘
     ↓
┌────────────────────────────────────────────────────────────────┐
│ SERVICE: OrderService.CreateOrder                              │
├────────────────────────────────────────────────────────────────┤
│ 1. Get user's cart (CartRepository.GetByUserID)                │
│ 2. Validate cart not empty                                     │
│ 3. Validate cart (CartService.ValidateCart)                    │
│    ├─ Check each product exists                                │
│    ├─ Check stock availability                                 │
│    └─ Return errors if any                                     │
│                                                                 │
│ 4. BEGIN TRANSACTION                                            │
│    ├─ Calculate total amount                                   │
│    ├─ Create order items                                       │
│    ├─ Deduct stock (ProductRepo.UpdateStockWithTx)             │
│    ├─ Create order (OrderRepo.CreateWithTx)                    │
│    └─ COMMIT or ROLLBACK                                       │
│                                                                 │
│ 5. Clear cart (CartService.ClearCart)                          │
│    └─ Delete stock reservations                                │
│                                                                 │
│ 6. If CC/DC payment:                                            │
│    └─ Create payment (PaymentService.CreatePaymentForOrder)    │
│                                                                 │
│ 7. Return created order                                        │
└────┬───────────────────────────────────────────────────────────┘
     ↓
┌────────────────────────────────────────────────────────────────┐
│ REPOSITORY: Multiple repositories called                       │
├────────────────────────────────────────────────────────────────┤
│ CartRepository.GetByUserID(userID)                             │
│   → SELECT * FROM carts WHERE user_id = $1                     │
│                                                                 │
│ ProductRepository.UpdateStockWithTx(tx, productID, -quantity)  │
│   → UPDATE products SET stock_quantity = stock_quantity + $2   │
│     WHERE id = $1                                              │
│                                                                 │
│ OrderRepository.CreateWithTx(tx, order)                        │
│   → INSERT INTO orders (...) VALUES (...)                      │
│   → INSERT INTO order_items (...) VALUES (...)                 │
│                                                                 │
│ CartRepository.ClearCart(cartID)                               │
│   → DELETE FROM cart_items WHERE cart_id = $1                  │
│   → DELETE FROM stock_reservations WHERE cart_id = $1          │
│                                                                 │
│ PaymentRepository.Create(payment)                              │
│   → INSERT INTO payments (...) VALUES (...)                    │
└────┬───────────────────────────────────────────────────────────┘
     ↓
┌────────────────────────────────────────────────────────────────┐
│ DATABASE: PostgreSQL Transaction                               │
├────────────────────────────────────────────────────────────────┤
│ BEGIN TRANSACTION;                                             │
│   UPDATE products SET stock_quantity = ... WHERE id = ...;     │
│   INSERT INTO orders (...) VALUES (...);                       │
│   INSERT INTO order_items (...) VALUES (...);                  │
│ COMMIT;                                                         │
│                                                                 │
│ DELETE FROM cart_items WHERE cart_id = ...;                    │
│ DELETE FROM stock_reservations WHERE cart_id = ...;            │
│ INSERT INTO payments (...) VALUES (...);                       │
└────┬───────────────────────────────────────────────────────────┘
     ↓ Return Results
┌────────────────────────────────────────────────────────────────┐
│ HANDLER: Format Response                                       │
├────────────────────────────────────────────────────────────────┤
│ {                                                              │
│   "success": true,                                             │
│   "message": "Order created successfully",                     │
│   "data": {                                                    │
│     "id": "...",                                               │
│     "order_number": "ORD-...",                                 │
│     "total_amount": 2999.98,                                   │
│     "status": "pending",                                       │
│     "items": [...],                                            │
│     ...                                                        │
│   }                                                            │
│ }                                                              │
└────┬───────────────────────────────────────────────────────────┘
     ↓ HTTP 201 Created
┌─────────┐
│ Client  │ ← Response
└─────────┘
```

## Authentication & Authorization Flow

```
┌──────────────────────────────────────────────────────────────────┐
│                    REGISTRATION FLOW                              │
└──────────────────────────────────────────────────────────────────┘

POST /api/v1/auth/register
{email, password, first_name, last_name}
           ↓
   [AuthHandler.Register]
           ↓
   Validate input
           ↓
   [AuthService.Register]
           ↓
   Check email exists? ──YES──→ Error: "User already exists"
           │
           NO
           ↓
   Hash password (bcrypt)
           ↓
   Create user with role="customer"
           ↓
   [UserRepository.Create]
           ↓
   INSERT INTO users
           ↓
   Generate JWT token
   {user_id, email, role="customer", exp, iat}
           ↓
   Return {user, access_token}


┌──────────────────────────────────────────────────────────────────┐
│                       LOGIN FLOW                                  │
└──────────────────────────────────────────────────────────────────┘

POST /api/v1/auth/login
{email, password}
           ↓
   [AuthHandler.Login]
           ↓
   Validate input
           ↓
   [AuthService.Login]
           ↓
   [UserRepository.GetByEmail]
           ↓
   User exists? ──NO──→ Error: "Invalid email or password"
           │
           YES
           ↓
   Verify password hash (bcrypt)
   bcrypt.CompareHashAndPassword(hash, password)
           ↓
   Match? ──NO──→ Error: "Invalid email or password"
           │
           YES
           ↓
   Generate JWT token
   {user_id, email, role, exp, iat}
   Signed with JWT_SECRET
           ↓
   Return {user, access_token}


┌──────────────────────────────────────────────────────────────────┐
│              AUTHENTICATED REQUEST FLOW                           │
└──────────────────────────────────────────────────────────────────┘

GET /api/v1/cart
Authorization: Bearer <JWT_TOKEN>
           ↓
   [GinAuthMiddleware]
           ↓
   Extract "Authorization" header
           ↓
   Format valid? ──NO──→ 401 Unauthorized
   "Bearer <token>"
           │
           YES
           ↓
   Parse JWT token
   jwt.Parse(token, secretFunc)
           ↓
   Valid signature? ──NO──→ 401 Invalid token
           │
           YES
           ↓
   Token expired? ──YES──→ 401 Token expired
           │
           NO
           ↓
   Extract claims:
   - user_id
   - email
   - role
           ↓
   Store in context:
   c.Set("userID", user_id)
   c.Set("userRole", role)
           ↓
   Continue to handler


┌──────────────────────────────────────────────────────────────────┐
│                 ADMIN REQUEST FLOW                                │
└──────────────────────────────────────────────────────────────────┘

GET /api/v1/admin/orders
Authorization: Bearer <JWT_TOKEN>
           ↓
   [GinAuthMiddleware]
   (validates token, sets userID & userRole)
           ↓
   [GinAdminMiddleware]
           ↓
   Get role from context
   role := c.Get("userRole")
           ↓
   Role exists? ──NO──→ 403 Forbidden
           │
           YES
           ↓
   role == "admin"? ──NO──→ 403 Admin access required
           │
           YES
           ↓
   Continue to admin handler
```

## Data Flow - Shopping Cart to Order

```
┌──────────────────────────────────────────────────────────────────┐
│                     SHOPPING JOURNEY                              │
└──────────────────────────────────────────────────────────────────┘

1. Add Product to Cart
   ─────────────────────
   POST /api/v1/cart/items
   {product_id, quantity: 2}
           ↓
   Check product exists
           ↓
   Check available stock:
   actual_stock - reserved_stock >= quantity?
           ↓
   YES → Reserve stock
   INSERT INTO stock_reservations
   (product_id, cart_id, quantity, expires_at=NOW()+10min)
           ↓
   Add to cart
   INSERT INTO cart_items
   (cart_id, product_id, quantity)
           ↓
   ┌─────────────────────────────────┐
   │ Cart State:                      │
   │ - Laptop x2                      │
   │ - Stock reserved for 10 minutes  │
   └─────────────────────────────────┘


2. Update Cart Item
   ─────────────────
   PUT /api/v1/cart/items/:itemId
   {quantity: 3}
           ↓
   Calculate diff: new - old = 3 - 2 = +1
           ↓
   Check additional stock available?
           ↓
   YES → Reserve additional stock (+1)
   UPDATE stock_reservations
   SET quantity = quantity + 1
           ↓
   Update cart item
   UPDATE cart_items
   SET quantity = 3
           ↓
   ┌─────────────────────────────────┐
   │ Cart State:                      │
   │ - Laptop x3                      │
   │ - Stock reserved for 10 minutes  │
   └─────────────────────────────────┘


3. Validate Cart Before Checkout
   ──────────────────────────────
   GET /api/v1/cart/validate
           ↓
   For each cart item:
   - Product still exists?
   - Stock still available?
   - Price changed? (optional warning)
           ↓
   Return validation status
   {
     "valid": true,
     "errors": []
   }


4. Create Order (Checkout)
   ────────────────────────
   POST /api/v1/orders
   {shipping_address, billing_address, payment_method: "cc"}
           ↓
   BEGIN TRANSACTION
           ↓
   Calculate total: 3 × $1499.99 = $4499.97
           ↓
   Deduct stock from inventory:
   UPDATE products
   SET stock_quantity = stock_quantity - 3
   WHERE id = laptop_id
           ↓
   Create order:
   INSERT INTO orders
   (user_id, order_number, total_amount, status="pending", ...)
           ↓
   Create order items:
   INSERT INTO order_items
   (order_id, product_id, quantity=3, price_at_time=1499.99)
           ↓
   Clear cart & reservations:
   DELETE FROM cart_items WHERE cart_id = ...
   DELETE FROM stock_reservations WHERE cart_id = ...
           ↓
   COMMIT TRANSACTION
           ↓
   Create payment (for CC/DC):
   INSERT INTO payments
   (order_id, amount, status="completed", method="cc")
           ↓
   ┌─────────────────────────────────┐
   │ Order Created:                   │
   │ - Order #ORD-1707302400-a1b2c3d4 │
   │ - Total: $4499.97                │
   │ - Status: pending                │
   │ - Payment: completed             │
   │ - Stock deducted                 │
   │ - Cart cleared                   │
   └─────────────────────────────────┘
```

## Stock Reservation System

```
┌──────────────────────────────────────────────────────────────────┐
│               STOCK RESERVATION MECHANISM                         │
└──────────────────────────────────────────────────────────────────┘

Product: Gaming Laptop
Actual Stock: 100 units

Timeline:
────────────────────────────────────────────────────────────────────

T=0: Initial State
─────────────────
┌─────────────────────────────┐
│ Product Stock: 100           │
│ Reserved: 0                  │
│ Available: 100 - 0 = 100     │
└─────────────────────────────┘


T=1: User A adds 30 to cart
────────────────────────────
POST /api/v1/cart/items {product_id, quantity: 30}
        ↓
Check available: 100 - 0 = 100 ≥ 30 ✓
        ↓
INSERT INTO stock_reservations
(product_id, cart_id_A, quantity=30, expires_at=T+10min)
        ↓
┌─────────────────────────────┐
│ Product Stock: 100           │
│ Reserved: 30 (User A)        │
│ Available: 100 - 30 = 70     │
└─────────────────────────────┘


T=2: User B adds 50 to cart (1 second later)
─────────────────────────────────────────────
POST /api/v1/cart/items {product_id, quantity: 50}
        ↓
Check available: 100 - 30 = 70 ≥ 50 ✓
        ↓
INSERT INTO stock_reservations
(product_id, cart_id_B, quantity=50, expires_at=T+10min)
        ↓
┌─────────────────────────────┐
│ Product Stock: 100           │
│ Reserved: 30 (A) + 50 (B) = 80│
│ Available: 100 - 80 = 20     │
└─────────────────────────────┘


T=3: User C tries to add 25 to cart
────────────────────────────────────
POST /api/v1/cart/items {product_id, quantity: 25}
        ↓
Check available: 100 - 80 = 20 ≥ 25 ✗
        ↓
Error: "Insufficient stock"
        ↓
┌─────────────────────────────┐
│ Product Stock: 100           │
│ Reserved: 80                 │
│ Available: 20                │
│ User C: REJECTED ✗           │
└─────────────────────────────┘


T=5: User A completes checkout
───────────────────────────────
POST /api/v1/orders
        ↓
BEGIN TRANSACTION
        ↓
UPDATE products
SET stock_quantity = stock_quantity - 30
WHERE id = product_id
        ↓
DELETE FROM stock_reservations
WHERE cart_id = cart_id_A
        ↓
COMMIT
        ↓
┌─────────────────────────────┐
│ Product Stock: 70 (updated)  │
│ Reserved: 50 (B only)        │
│ Available: 70 - 50 = 20      │
└─────────────────────────────┘


T=12: Reservation B expires (10 minutes)
─────────────────────────────────────────
Automatic expiry (expires_at < NOW())
        ↓
Query filters out expired reservations:
LEFT JOIN stock_reservations sr
ON p.id = sr.product_id
AND sr.expires_at > NOW()
        ↓
┌─────────────────────────────┐
│ Product Stock: 70            │
│ Reserved: 0 (B expired)      │
│ Available: 70 - 0 = 70       │
└─────────────────────────────┘

Note: Background job can clean up expired rows:
DELETE FROM stock_reservations
WHERE expires_at < NOW();
```

## Order State Machine

```
┌──────────────────────────────────────────────────────────────────┐
│                    ORDER STATUS LIFECYCLE                         │
└──────────────────────────────────────────────────────────────────┘

                    ┌──────────────┐
        ┌───────────┤   PENDING    ├───────────┐
        │           └──────┬───────┘           │
        │                  │                   │
        │                  │ Admin:            │ Customer/Admin:
        │                  │ Update Status     │ Cancel Order
        │                  ↓                   ↓
        │           ┌──────────────┐     ┌─────────────┐
        │           │  PROCESSING  │     │  CANCELLED  │
        │           └──────┬───────┘     └─────────────┘
        │                  │                   (Terminal)
        │                  │ Admin:
        │                  │ Update Status
        │                  ↓
        │           ┌──────────────┐
        │           │   SHIPPED    │
        │           └──────┬───────┘
        │                  │
        │                  │ Admin:
        │                  │ Confirm Delivery
        │                  ↓
        │           ┌──────────────┐
        │           │  DELIVERED   │────┐
        │           └──────┬───────┘    │
        │                  │            │ Customer:
        │                  │            │ Request Return
        │                  │            ↓
        │                  │      ┌──────────────────┐
        │                  │      │ RETURN_REQUESTED │
        │                  │      └──────────────────┘
        │                  │            │
        │                  │            │ Admin: Approve
        │                  │            ↓
        │                  │      ┌──────────────┐
        │                  │      │   REFUNDED   │
        │                  │      └──────────────┘
        │                  │            (Terminal)
        │                  ↓
        │           ┌──────────────┐
        └───────────┤  COMPLETED   │
                    └──────────────┘
                     (Terminal)


Valid Transitions:
──────────────────
pending      → processing, cancelled
processing   → shipped, delivered, completed, cancelled
shipped      → delivered, completed
delivered    → completed
completed    → (none - terminal state)
cancelled    → (none - terminal state)
refunded     → (none - terminal state)


Payment Handling:
─────────────────
- CC/DC: Payment created immediately on order creation
         Status = "completed" (mock)

- COD:   Payment created when status changes to "delivered"
         Status = "completed" after cash collected


Return Flow:
────────────
1. Customer: Create return request
   → Order status: delivered/completed → return_requested

2. Admin: Review and approve/reject
   → If approved: Process refund, restore stock
   → Order status: return_requested → refunded

3. Admin: If rejected
   → Order status: return_requested → delivered
```

## Database Transaction Example

```
┌──────────────────────────────────────────────────────────────────┐
│           DATABASE TRANSACTION - ORDER CREATION                   │
└──────────────────────────────────────────────────────────────────┘

BEGIN TRANSACTION;
─────────────────

1. Deduct Stock
   ────────────
   UPDATE products
   SET stock_quantity = stock_quantity - 2
   WHERE id = 'laptop-uuid'
   AND stock_quantity >= 2;  -- Ensure sufficient stock

   Rows affected: 1 ✓


2. Create Order
   ────────────
   INSERT INTO orders (
       id, user_id, order_number, total_amount,
       status, payment_method, shipping_address,
       billing_address, created_at, updated_at
   ) VALUES (
       '550e8400-...', '123e4567-...', 'ORD-1707302400-a1b2c3d4',
       2999.98, 'pending', 'cc', '{"street": "123 Main St", ...}',
       '{"street": "123 Main St", ...}', NOW(), NOW()
   )
   RETURNING id;

   Order ID: 550e8400-... ✓


3. Create Order Items
   ──────────────────
   INSERT INTO order_items (
       id, order_id, product_id, quantity, price_at_time
   ) VALUES (
       'item-uuid-1', '550e8400-...', 'laptop-uuid', 2, 1499.99
   );

   Rows affected: 1 ✓


4. Clear Cart Items
   ────────────────
   DELETE FROM cart_items
   WHERE cart_id = 'cart-uuid';

   Rows affected: 1 ✓


5. Release Stock Reservations
   ──────────────────────────
   DELETE FROM stock_reservations
   WHERE cart_id = 'cart-uuid';

   Rows affected: 1 ✓


COMMIT;
───────

Result: ✓ Order created successfully
        ✓ Stock deducted
        ✓ Cart cleared
        ✓ Reservations released


IF ANY STEP FAILS:
─────────────────
ROLLBACK;
         ↓
All changes reverted:
- Stock quantity restored
- No order created
- Cart remains unchanged
- Reservations preserved
         ↓
Return Error to Client:
{
  "success": false,
  "message": "Failed to create order",
  "error": "Insufficient stock"
}
```

---

**End of Architecture Diagrams**
