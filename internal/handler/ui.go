package handler

import (
	"net/http"
	"strings"
)

type UIHandler struct {
	version     string
	environment string
	appName     string
}

func NewUIHandler(version, environment, appName string) *UIHandler {
	return &UIHandler{
		version:     version,
		environment: environment,
		appName:     appName,
	}
}

func (u *UIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/dashboard" {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "404 page not found",
			"message": "The requested endpoint does not exist. Visit / for the enterprise dashboard or /api/v1/orders for the API.",
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(u.renderHTML()))
}

func (u *UIHandler) renderHTML() string {
	html := dashboardHTML
	html = strings.ReplaceAll(html, "{{VERSION}}", u.version)
	html = strings.ReplaceAll(html, "{{ENV}}", u.environment)
	html = strings.ReplaceAll(html, "{{APP_NAME}}", u.appName)
	return html
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Enterprise Microservice Control Center</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg: #030712;
      --card-bg: rgba(17, 24, 39, 0.7);
      --card-border: rgba(255, 255, 255, 0.08);
      --card-hover: rgba(99, 102, 241, 0.3);
      --primary: #6366f1;
      --primary-hover: #4f46e5;
      --cyan: #06b6d4;
      --green: #10b981;
      --yellow: #f59e0b;
      --red: #ef4444;
      --text: #f9fafb;
      --muted: #9ca3af;
    }

    * { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background: radial-gradient(circle at 50% 0%, #1e1b4b 0%, #030712 60%);
      background-attachment: fixed;
      color: var(--text);
      font-family: 'Inter', system-ui, sans-serif;
      min-height: 100vh;
      line-height: 1.5;
    }

    header {
      position: sticky;
      top: 0;
      z-index: 100;
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      background: rgba(3, 7, 18, 0.85);
      border-bottom: 1px solid var(--card-border);
      padding: 1rem 2rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .brand {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-weight: 800;
      font-size: 1.25rem;
      background: linear-gradient(135deg, #a5b4fc, #38bdf8);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .brand-icon {
      width: 38px;
      height: 38px;
      background: linear-gradient(135deg, var(--primary), var(--cyan));
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.2rem;
      box-shadow: 0 0 15px rgba(99, 102, 241, 0.5);
    }

    .status-container {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .version-tag {
      font-size: 0.85rem;
      background: rgba(99, 102, 241, 0.15);
      border: 1px solid rgba(99, 102, 241, 0.3);
      color: #a5b4fc;
      padding: 0.25rem 0.65rem;
      border-radius: 6px;
      font-weight: 600;
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(16, 185, 129, 0.1);
      border: 1px solid rgba(16, 185, 129, 0.3);
      color: var(--green);
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      font-size: 0.85rem;
      font-weight: 600;
    }

    .pulse-dot {
      width: 8px;
      height: 8px;
      background-color: var(--green);
      border-radius: 50%;
      box-shadow: 0 0 10px var(--green);
      animation: pulse 2s infinite;
    }

    @keyframes pulse {
      0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
      70% { transform: scale(1.1); box-shadow: 0 0 0 8px rgba(16, 185, 129, 0); }
      100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
    }

    main {
      max-width: 1280px;
      margin: 2rem auto;
      padding: 0 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 2rem;
    }

    .grid-4 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 1.25rem;
    }

    .card {
      background: var(--card-bg);
      backdrop-filter: blur(12px);
      -webkit-backdrop-filter: blur(12px);
      border: 1px solid var(--card-border);
      border-radius: 16px;
      padding: 1.5rem;
      transition: all 0.3s ease;
      box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
    }

    .card:hover {
      border-color: var(--card-hover);
      transform: translateY(-2px);
    }

    .card-title {
      font-size: 0.8rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--muted);
      margin-bottom: 0.4rem;
      display: flex;
      justify-content: space-between;
    }

    .card-value {
      font-size: 1.8rem;
      font-weight: 800;
      color: var(--text);
    }

    /* Toolbar & Search */
    .toolbar {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1rem;
      margin-bottom: 1.25rem;
    }

    .filter-group {
      display: flex;
      gap: 0.5rem;
      align-items: center;
    }

    .filter-btn {
      background: rgba(255, 255, 255, 0.05);
      border: 1px solid var(--card-border);
      color: var(--muted);
      padding: 0.4rem 0.85rem;
      border-radius: 8px;
      font-size: 0.85rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
    }

    .filter-btn.active, .filter-btn:hover {
      background: rgba(99, 102, 241, 0.2);
      border-color: var(--primary);
      color: var(--text);
    }

    .search-input {
      background: rgba(3, 7, 18, 0.6);
      border: 1px solid var(--card-border);
      color: var(--text);
      padding: 0.55rem 1rem;
      border-radius: 10px;
      font-size: 0.9rem;
      width: 240px;
      outline: none;
      transition: all 0.2s;
    }

    .search-input:focus {
      border-color: var(--primary);
      box-shadow: 0 0 10px rgba(99, 102, 241, 0.3);
    }

    .btn {
      background: linear-gradient(135deg, var(--primary), var(--primary-hover));
      color: #fff;
      border: none;
      padding: 0.6rem 1.25rem;
      border-radius: 10px;
      font-weight: 600;
      font-size: 0.9rem;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      transition: all 0.2s ease;
      box-shadow: 0 4px 14px 0 rgba(99, 102, 241, 0.4);
    }

    .btn:hover {
      opacity: 0.95;
      transform: scale(1.02);
    }

    .btn-secondary {
      background: rgba(255, 255, 255, 0.05);
      border: 1px solid var(--card-border);
      color: var(--text);
      padding: 0.5rem 1rem;
      border-radius: 8px;
      font-size: 0.85rem;
      font-weight: 600;
      cursor: pointer;
    }

    .btn-secondary:hover {
      background: rgba(255, 255, 255, 0.1);
    }

    .btn-danger {
      background: rgba(239, 68, 68, 0.15);
      border: 1px solid rgba(239, 68, 68, 0.3);
      color: var(--red);
      padding: 0.3rem 0.6rem;
      border-radius: 6px;
      font-size: 0.75rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn-danger:hover {
      background: rgba(239, 68, 68, 0.3);
    }

    /* Table Component */
    .table-wrapper {
      overflow-x: auto;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
    }

    th {
      padding: 0.85rem 1rem;
      color: var(--muted);
      font-size: 0.75rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      border-bottom: 1px solid var(--card-border);
    }

    td {
      padding: 1rem;
      border-bottom: 1px solid rgba(255, 255, 255, 0.04);
      font-size: 0.9rem;
    }

    tr:hover td {
      background: rgba(255, 255, 255, 0.02);
    }

    .status-select {
      background: #030712;
      border: 1px solid var(--card-border);
      color: var(--text);
      padding: 0.3rem 0.6rem;
      border-radius: 6px;
      font-size: 0.8rem;
      font-weight: 600;
      outline: none;
      cursor: pointer;
    }

    .status-select.CREATED { color: #38bdf8; border-color: rgba(56, 189, 248, 0.4); }
    .status-select.PROCESSING { color: var(--yellow); border-color: rgba(245, 158, 11, 0.4); }
    .status-select.COMPLETED { color: var(--green); border-color: rgba(16, 185, 129, 0.4); }
    .status-select.CANCELLED { color: var(--red); border-color: rgba(239, 68, 68, 0.4); }

    /* Audit Logs Feed */
    .audit-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      max-height: 250px;
      overflow-y: auto;
      margin-top: 1rem;
    }

    .audit-item {
      background: rgba(3, 7, 18, 0.6);
      border: 1px solid var(--card-border);
      padding: 0.75rem 1rem;
      border-radius: 10px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-size: 0.85rem;
    }

    .audit-badge {
      font-size: 0.7rem;
      font-weight: 700;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
    }
    .badge-CREATE { background: rgba(16, 185, 129, 0.2); color: var(--green); }
    .badge-UPDATE_STATUS { background: rgba(56, 189, 248, 0.2); color: #38bdf8; }
    .badge-DELETE { background: rgba(239, 68, 68, 0.2); color: var(--red); }
    .badge-DISCOUNT { background: rgba(245, 158, 11, 0.2); color: var(--yellow); }

    /* API Sandbox Panel */
    .endpoint-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 1rem;
      margin-top: 1rem;
    }

    .endpoint-card {
      background: rgba(3, 7, 18, 0.6);
      border: 1px solid var(--card-border);
      padding: 0.85rem 1rem;
      border-radius: 12px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      cursor: pointer;
      transition: all 0.2s;
    }

    .endpoint-card:hover {
      border-color: var(--primary);
    }

    .method {
      font-weight: 700;
      font-size: 0.75rem;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      margin-right: 0.5rem;
    }
    .method-get { background: rgba(16, 185, 129, 0.2); color: var(--green); }

    pre {
      background: #030712;
      border: 1px solid var(--card-border);
      padding: 1.25rem;
      border-radius: 12px;
      overflow-x: auto;
      font-family: 'Fira Code', monospace;
      color: #38bdf8;
      font-size: 0.85rem;
      max-height: 280px;
    }

    /* Modal Form */
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.75);
      backdrop-filter: blur(8px);
      display: none;
      align-items: center;
      justify-content: center;
      z-index: 200;
    }

    .modal {
      background: #0f172a;
      border: 1px solid var(--card-hover);
      border-radius: 20px;
      width: 100%;
      max-width: 480px;
      padding: 2rem;
      box-shadow: 0 25px 60px rgba(0, 0, 0, 0.8);
    }

    .form-group {
      margin-bottom: 1.25rem;
    }

    label {
      display: block;
      font-size: 0.85rem;
      color: var(--muted);
      margin-bottom: 0.4rem;
    }

    input, select {
      width: 100%;
      background: #030712;
      border: 1px solid var(--card-border);
      color: var(--text);
      padding: 0.75rem 1rem;
      border-radius: 10px;
      font-family: inherit;
      outline: none;
    }

    input:focus { border-color: var(--primary); }

    /* Toast Notification */
    .toast-container {
      position: fixed;
      bottom: 2rem;
      right: 2rem;
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      z-index: 300;
    }

    .toast {
      background: #0f172a;
      border: 1px solid var(--card-border);
      border-left: 4px solid var(--primary);
      padding: 0.85rem 1.25rem;
      border-radius: 10px;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
      display: flex;
      align-items: center;
      gap: 0.75rem;
      font-size: 0.9rem;
      animation: slideIn 0.3s ease;
    }

    @keyframes slideIn {
      from { transform: translateX(100%); opacity: 0; }
      to { transform: translateX(0); opacity: 1; }
    }
  </style>
</head>
<body>

  <header>
    <div class="brand">
      <div class="brand-icon">⚡</div>
      <div>{{APP_NAME}}</div>
    </div>
    <div class="status-container">
      <span class="version-tag">v{{VERSION}} ({{ENV}})</span>
      <div class="status-badge">
        <span class="pulse-dot"></span> LIVE APP OK
      </div>
    </div>
  </header>

  <main>

    <!-- Metrics Header Cards -->
    <div class="grid-4">
      <div class="card">
        <div class="card-title"><span>Total Orders</span> 📦</div>
        <div class="card-value" id="metric-orders">-</div>
      </div>
      <div class="card">
        <div class="card-title"><span>Total Revenue</span> 💵</div>
        <div class="card-value" style="color: var(--green);" id="metric-revenue">$0.00</div>
      </div>
      <div class="card">
        <div class="card-title"><span>Average Order Value</span> 📈</div>
        <div class="card-value" style="color: var(--cyan);" id="metric-aov">$0.00</div>
      </div>
      <div class="card">
        <div class="card-title"><span>Top Customer</span> 🏆</div>
        <div class="card-value" style="font-size: 1.3rem; color: #a855f7;" id="metric-topcust">-</div>
      </div>
    </div>

    <!-- Live Orders Table Section -->
    <div class="card">
      <div class="toolbar">
        <div>
          <h2 style="font-size: 1.3rem; font-weight: 700;">Live Domain Orders</h2>
          <p style="font-size: 0.85rem; color: var(--muted);">Real-time interactive order management & batch operations</p>
        </div>
        <div style="display: flex; gap: 0.5rem; align-items: center; flex-wrap: wrap;">
          <button class="btn-secondary" onclick="window.open('/api/v1/orders/export?format=csv')">📥 Export CSV</button>
          <button class="btn-secondary" onclick="window.open('/api/v1/orders/export?format=json')">📥 Export JSON</button>
          <button class="btn" onclick="openModal()">+ Create Order</button>
        </div>
      </div>

      <div class="toolbar" style="margin-bottom: 1rem;">
        <div class="filter-group">
          <button class="filter-btn active" onclick="setFilter('ALL', this)">All</button>
          <button class="filter-btn" onclick="setFilter('CREATED', this)">Created</button>
          <button class="filter-btn" onclick="setFilter('PROCESSING', this)">Processing</button>
          <button class="filter-btn" onclick="setFilter('COMPLETED', this)">Completed</button>
          <button class="filter-btn" onclick="setFilter('CANCELLED', this)">Cancelled</button>
        </div>
        <div style="display: flex; gap: 0.75rem; align-items: center;">
          <input type="text" class="search-input" id="searchInput" placeholder="🔍 Search customer, item..." oninput="renderTable()" />
          <div style="font-size: 0.8rem; color: var(--muted);">
            Auto-Refresh: <input type="checkbox" id="autoRefresh" checked /> (3s)
          </div>
        </div>
      </div>

      <!-- Promo Coupon Application Widget -->
      <div style="background: rgba(3, 7, 18, 0.5); border: 1px solid var(--card-border); padding: 0.85rem 1.25rem; border-radius: 12px; margin-bottom: 1.25rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.75rem;">
        <div style="font-size: 0.85rem; color: var(--muted); font-weight: 500;">
          🏷️ <strong style="color: var(--text);">Promo Discount Simulator:</strong> Apply global discounts to active orders
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <select id="couponSelect" style="padding: 0.35rem 0.75rem; border-radius: 8px; font-size: 0.85rem; background: #030712; color: var(--text); border: 1px solid var(--card-border);">
            <option value="ENTERPRISE10">ENTERPRISE10 (10% Off)</option>
            <option value="CLOUD20">CLOUD20 (20% Off)</option>
            <option value="DEVOPS30">DEVOPS30 (30% Off)</option>
          </select>
          <button class="btn-secondary" onclick="applyPromo()">Apply Code</button>
        </div>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Customer Name</th>
              <th>Item Description</th>
              <th>Qty</th>
              <th>Unit Price</th>
              <th>Total Amount</th>
              <th>Status (Interactive)</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody id="orders-list">
            <tr><td colspan="8" style="text-align: center; color: var(--muted);">Loading orders...</td></tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Audit Logs & Governance Section -->
    <div class="card">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <div>
          <h2 style="font-size: 1.2rem; font-weight: 700;">Audit & Compliance Feed</h2>
          <p style="font-size: 0.85rem; color: var(--muted);">Immutable activity logging for compliance auditing</p>
        </div>
        <button class="btn-secondary" onclick="fetchAuditLogs()">🔄 Refresh Feed</button>
      </div>

      <div class="audit-list" id="audit-list">
        <div style="text-align: center; color: var(--muted); padding: 1rem;">Loading audit feed...</div>
      </div>
    </div>

    <!-- Interactive API Explorer -->
    <div class="card">
      <h2 style="font-size: 1.25rem; font-weight: 700; margin-bottom: 0.25rem;">Interactive API Sandbox</h2>
      <p style="font-size: 0.85rem; color: var(--muted);">Test live Go microservice endpoints and inspect JSON payloads</p>

      <div class="endpoint-grid">
        <div class="endpoint-card" onclick="testEndpoint('/api/v1/metrics')">
          <div><span class="method method-get">GET</span>/api/v1/metrics</div>
        </div>
        <div class="endpoint-card" onclick="testEndpoint('/api/v1/audit-logs')">
          <div><span class="method method-get">GET</span>/api/v1/audit-logs</div>
        </div>
        <div class="endpoint-card" onclick="testEndpoint('/api/v1/orders')">
          <div><span class="method method-get">GET</span>/api/v1/orders</div>
        </div>
        <div class="endpoint-card" onclick="testEndpoint('/healthz')">
          <div><span class="method method-get">GET</span>/healthz</div>
        </div>
      </div>

      <div style="margin-top: 1.5rem;">
        <pre id="api-output">// Click any API endpoint card above to inspect live JSON response</pre>
      </div>
    </div>

  </main>

  <!-- Create Order Modal -->
  <div class="modal-overlay" id="orderModal">
    <div class="modal">
      <h3 style="font-size: 1.3rem; margin-bottom: 1.25rem; font-weight: 700;">Create New Enterprise Order</h3>
      <form id="createOrderForm" onsubmit="submitOrder(event)">
        <div class="form-group">
          <label>Customer Name</label>
          <input type="text" id="custName" placeholder="Acme Logistics Inc." required />
        </div>
        <div class="form-group">
          <label>Item Description</label>
          <input type="text" id="itemDesc" placeholder="Kubernetes Managed Cluster" required />
        </div>
        <div class="form-group">
          <label>Quantity</label>
          <input type="number" id="itemQty" value="5" min="1" required />
        </div>
        <div class="form-group">
          <label>Unit Price ($)</label>
          <input type="number" step="0.01" id="itemPrice" value="499.00" required />
        </div>
        <div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.75rem;">
          <button type="button" class="btn-secondary" onclick="closeModal()">Cancel</button>
          <button type="submit" class="btn">Create Order</button>
        </div>
      </form>
    </div>
  </div>

  <div class="toast-container" id="toastContainer"></div>

  <script>
    let allOrders = [];
    let currentFilter = 'ALL';

    async function loadData() {
      await Promise.all([fetchOrders(), fetchMetrics(), fetchAuditLogs()]);
    }

    async function fetchMetrics() {
      try {
        const res = await fetch('/api/v1/metrics');
        const m = await res.json();
        document.getElementById('metric-orders').textContent = m.total_orders;
        document.getElementById('metric-revenue').textContent = '$' + m.total_revenue.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2});
        document.getElementById('metric-aov').textContent = '$' + m.avg_order_value.toFixed(2);
        document.getElementById('metric-topcust').textContent = m.top_customer;
      } catch (err) {
        console.error('Failed metrics fetch:', err);
      }
    }

    async function fetchOrders() {
      try {
        const res = await fetch('/api/v1/orders');
        allOrders = await res.json();
        renderTable();
      } catch (err) {
        console.error('Failed orders fetch:', err);
      }
    }

    async function fetchAuditLogs() {
      try {
        const res = await fetch('/api/v1/audit-logs');
        const logs = await res.json();
        const container = document.getElementById('audit-list');
        container.innerHTML = '';

        if (!logs || logs.length === 0) {
          container.innerHTML = '<div style="text-align: center; color: var(--muted);">No audit events recorded yet.</div>';
          return;
        }

        logs.slice(0, 10).forEach(log => {
          const div = document.createElement('div');
          div.className = 'audit-item';
          const timeStr = new Date(log.timestamp).toLocaleTimeString();
          div.innerHTML = 
            '<div>' +
              '<span class="audit-badge badge-' + log.action + '">' + log.action + '</span> ' +
              '<strong style="color: var(--text);">' + escapeHtml(log.details) + '</strong>' +
            '</div>' +
            '<div style="color: var(--muted); font-size: 0.75rem;">' + timeStr + '</div>';
          container.appendChild(div);
        });
      } catch (err) {
        console.error('Failed audit logs fetch:', err);
      }
    }

    function renderTable() {
      const tbody = document.getElementById('orders-list');
      const search = document.getElementById('searchInput').value.toLowerCase();
      tbody.innerHTML = '';

      const filtered = allOrders.filter(o => {
        const matchesFilter = currentFilter === 'ALL' || o.status === currentFilter;
        const matchesSearch = o.customer_name.toLowerCase().includes(search) || 
                              o.item.toLowerCase().includes(search) || 
                              o.id.toLowerCase().includes(search);
        return matchesFilter && matchesSearch;
      });

      if (filtered.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; color: var(--muted); padding: 2rem;">No matching orders found.</td></tr>';
        return;
      }

      filtered.forEach(o => {
        const tr = document.createElement('tr');
        const total = (o.price * o.quantity).toFixed(2);
        
        tr.innerHTML = 
          '<td style="font-weight: 700; color: #38bdf8;">' + o.id + '</td>' +
          '<td>' + escapeHtml(o.customer_name) + '</td>' +
          '<td>' + escapeHtml(o.item) + '</td>' +
          '<td>' + o.quantity + '</td>' +
          '<td>$' + o.price.toFixed(2) + '</td>' +
          '<td style="font-weight: 600;">$' + total + '</td>' +
          '<td>' +
            '<select class="status-select ' + o.status + '" onchange="updateStatus(\'' + o.id + '\', this.value)">' +
              '<option value="CREATED"' + (o.status==='CREATED'?' selected':'') + '>CREATED</option>' +
              '<option value="PROCESSING"' + (o.status==='PROCESSING'?' selected':'') + '>PROCESSING</option>' +
              '<option value="COMPLETED"' + (o.status==='COMPLETED'?' selected':'') + '>COMPLETED</option>' +
              '<option value="CANCELLED"' + (o.status==='CANCELLED'?' selected':'') + '>CANCELLED</option>' +
            '</select>' +
          '</td>' +
          '<td><button class="btn-danger" onclick="deleteOrder(\'' + o.id + '\')">Delete</button></td>';

        tbody.appendChild(tr);
      });
    }

    function setFilter(filter, btn) {
      currentFilter = filter;
      document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      renderTable();
    }

    async function applyPromo() {
      const code = document.getElementById('couponSelect').value;
      try {
        const res = await fetch('/api/v1/orders/discount', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ coupon_code: code })
        });
        const data = await res.json();
        if (res.ok) {
          showToast('Applied ' + code + '! Reduced prices on ' + data.affected_orders + ' active orders');
          loadData();
        } else {
          showToast(data.message || 'Failed applying coupon', true);
        }
      } catch (err) {
        showToast('Error applying coupon', true);
      }
    }

    async function updateStatus(id, newStatus) {
      try {
        const res = await fetch('/api/v1/orders/' + id + '/status', {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ status: newStatus })
        });
        if (res.ok) {
          showToast('Updated ' + id + ' status to ' + newStatus);
          loadData();
        }
      } catch (err) {
        showToast('Error updating status', true);
      }
    }

    async function deleteOrder(id) {
      if (!confirm('Are you sure you want to delete order ' + id + '?')) return;
      try {
        const res = await fetch('/api/v1/orders/' + id, { method: 'DELETE' });
        if (res.ok) {
          showToast('Order ' + id + ' deleted successfully');
          loadData();
        }
      } catch (err) {
        showToast('Error deleting order', true);
      }
    }

    async function submitOrder(e) {
      e.preventDefault();
      const payload = {
        customer_name: document.getElementById('custName').value,
        item: document.getElementById('itemDesc').value,
        quantity: parseInt(document.getElementById('itemQty').value),
        price: parseFloat(document.getElementById('itemPrice').value)
      };

      try {
        const res = await fetch('/api/v1/orders', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });

        if (res.ok) {
          closeModal();
          showToast('New order created successfully!');
          loadData();
          document.getElementById('createOrderForm').reset();
        }
      } catch (err) {
        showToast('Error creating order', true);
      }
    }

    async function testEndpoint(path) {
      const output = document.getElementById('api-output');
      output.textContent = 'Fetching ' + path + '...';
      try {
        const res = await fetch(path);
        const json = await res.json();
        output.textContent = JSON.stringify(json, null, 2);
      } catch (err) {
        output.textContent = 'Error: ' + err.message;
      }
    }

    function showToast(msg, isError = false) {
      const container = document.getElementById('toastContainer');
      const toast = document.createElement('div');
      toast.className = 'toast';
      if (isError) toast.style.borderLeftColor = 'var(--red)';
      toast.innerHTML = (isError ? '⚠️ ' : '✅ ') + escapeHtml(msg);
      container.appendChild(toast);
      setTimeout(() => toast.remove(), 3500);
    }

    function openModal() { document.getElementById('orderModal').style.display = 'flex'; }
    function closeModal() { document.getElementById('orderModal').style.display = 'none'; }
    function escapeHtml(str) { return str.replace(/[&<>"']/g, m => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[m]); }

    // Initial load & 3-second Auto Refresh
    loadData();
    setInterval(() => {
      if (document.getElementById('autoRefresh').checked) {
        loadData();
      }
    }, 3000);
  </script>
</body>
</html>
`
