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
	// Serve root dashboard only for root path '/' or '/dashboard'
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
  <title>Enterprise Order Microservice | Dashboard</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
  <style>
    :root {
      --bg-gradient: radial-gradient(circle at 50% 0%, #1e1b4b 0%, #0f172a 50%, #020617 100%);
      --card-bg: rgba(15, 23, 42, 0.75);
      --card-border: rgba(255, 255, 255, 0.08);
      --card-border-hover: rgba(99, 102, 241, 0.4);
      --primary: #6366f1;
      --primary-hover: #4f46e5;
      --accent-cyan: #06b6d4;
      --accent-green: #10b981;
      --text-main: #f8fafc;
      --text-muted: #94a3b8;
      --font-family: 'Inter', system-ui, -apple-system, sans-serif;
    }

    * { box-sizing: border-box; margin: 0; padding: 0; }
    
    body {
      background: var(--bg-gradient);
      background-attachment: fixed;
      color: var(--text-main);
      font-family: var(--font-family);
      min-height: 100vh;
      line-height: 1.5;
    }

    header {
      position: sticky;
      top: 0;
      z-index: 100;
      backdrop-filter: blur(16px);
      -webkit-backdrop-filter: blur(16px);
      background: rgba(2, 6, 23, 0.8);
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
      font-weight: 700;
      font-size: 1.25rem;
      background: linear-gradient(135deg, #a5b4fc, #38bdf8);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .brand-icon {
      width: 36px;
      height: 36px;
      background: linear-gradient(135deg, var(--primary), var(--accent-cyan));
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.1rem;
      box-shadow: 0 0 15px rgba(99, 102, 241, 0.5);
    }

    .status-badge {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(16, 185, 129, 0.1);
      border: 1px solid rgba(16, 185, 129, 0.3);
      color: var(--accent-green);
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      font-size: 0.85rem;
      font-weight: 600;
    }

    .pulse-dot {
      width: 8px;
      height: 8px;
      background-color: var(--accent-green);
      border-radius: 50%;
      box-shadow: 0 0 10px var(--accent-green);
      animation: pulse 2s infinite;
    }

    @keyframes pulse {
      0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
      70% { transform: scale(1.1); box-shadow: 0 0 0 8px rgba(16, 185, 129, 0); }
      100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
    }

    main {
      max-width: 1200px;
      margin: 2rem auto;
      padding: 0 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 2rem;
    }

    .grid-3 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 1.5rem;
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
      border-color: var(--card-border-hover);
      transform: translateY(-2px);
    }

    .card-title {
      font-size: 0.9rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--text-muted);
      margin-bottom: 0.5rem;
    }

    .card-value {
      font-size: 1.75rem;
      font-weight: 700;
      color: var(--text-main);
    }

    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1.25rem;
    }

    .btn {
      background: linear-gradient(135deg, var(--primary), var(--primary-hover));
      color: #fff;
      border: none;
      padding: 0.6rem 1.2rem;
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
      color: var(--text-main);
      box-shadow: none;
    }

    .btn-secondary:hover {
      background: rgba(255, 255, 255, 0.1);
    }

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
      color: var(--text-muted);
      font-size: 0.8rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      border-bottom: 1px solid var(--card-border);
    }

    td {
      padding: 1rem;
      border-bottom: 1px solid rgba(255, 255, 255, 0.04);
      font-size: 0.925rem;
    }

    tr:last-child td { border-bottom: none; }

    .tag {
      padding: 0.25rem 0.65rem;
      border-radius: 6px;
      font-size: 0.75rem;
      font-weight: 600;
      display: inline-block;
    }

    .tag-created { background: rgba(56, 189, 248, 0.15); color: #38bdf8; border: 1px solid rgba(56, 189, 248, 0.3); }

    .endpoint-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 1rem;
      margin-top: 1rem;
    }

    .endpoint-card {
      background: rgba(2, 6, 23, 0.5);
      border: 1px solid var(--card-border);
      padding: 1rem;
      border-radius: 12px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .method {
      font-weight: 700;
      font-size: 0.75rem;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      margin-right: 0.5rem;
    }
    .method-get { background: rgba(16, 185, 129, 0.2); color: var(--accent-green); }

    pre {
      background: #020617;
      border: 1px solid var(--card-border);
      padding: 1rem;
      border-radius: 10px;
      overflow-x: auto;
      font-family: monospace;
      color: #38bdf8;
      font-size: 0.85rem;
      max-height: 250px;
    }

    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(6px);
      display: none;
      align-items: center;
      justify-content: center;
      z-index: 200;
    }

    .modal {
      background: #0f172a;
      border: 1px solid var(--card-border-hover);
      border-radius: 16px;
      width: 100%;
      max-width: 480px;
      padding: 1.75rem;
      box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6);
    }

    .form-group {
      margin-bottom: 1.2rem;
    }

    label {
      display: block;
      font-size: 0.85rem;
      color: var(--text-muted);
      margin-bottom: 0.4rem;
    }

    input {
      width: 100%;
      background: #020617;
      border: 1px solid var(--card-border);
      color: var(--text-main);
      padding: 0.75rem 1rem;
      border-radius: 8px;
      font-family: inherit;
      outline: none;
    }

    input:focus {
      border-color: var(--primary);
    }
  </style>
</head>
<body>

  <header>
    <div class="brand">
      <div class="brand-icon">⚡</div>
      <div>Enterprise Go Platform</div>
    </div>
    <div class="status-badge">
      <span class="pulse-dot"></span> SYSTEM HEALTHY
    </div>
  </header>

  <main>

    <div class="grid-3">
      <div class="card">
        <div class="card-title">Microservice Name</div>
        <div class="card-value">{{APP_NAME}}</div>
      </div>
      <div class="card">
        <div class="card-title">Active Release Version</div>
        <div class="card-value" style="color: #38bdf8;">{{VERSION}}</div>
      </div>
      <div class="card">
        <div class="card-title">Deployment Environment</div>
        <div class="card-value" style="color: var(--accent-green);">{{ENV}}</div>
      </div>
    </div>

    <div class="card">
      <div class="table-header">
        <div>
          <h2 style="font-size: 1.25rem; font-weight: 600;">Enterprise Orders</h2>
          <p style="font-size: 0.85rem; color: var(--text-muted);">Real-time synchronized domain records</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn btn-secondary" onclick="fetchOrders()">🔄 Refresh</button>
          <button class="btn" onclick="openModal()">+ Create Order</button>
        </div>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Customer Name</th>
              <th>Item</th>
              <th>Quantity</th>
              <th>Total Price</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody id="orders-list">
            <tr><td colspan="6" style="text-align: center; color: var(--text-muted);">Loading live orders...</td></tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card">
      <h2 style="font-size: 1.25rem; font-weight: 600; margin-bottom: 0.25rem;">Interactive API Explorer</h2>
      <p style="font-size: 0.85rem; color: var(--text-muted);">Test live microservice endpoints in real-time</p>

      <div class="endpoint-grid">
        <div class="endpoint-card">
          <div><span class="method method-get">GET</span>/healthz</div>
          <button class="btn btn-secondary" onclick="testEndpoint('/healthz')">Test</button>
        </div>
        <div class="endpoint-card">
          <div><span class="method method-get">GET</span>/livez</div>
          <button class="btn btn-secondary" onclick="testEndpoint('/livez')">Test</button>
        </div>
        <div class="endpoint-card">
          <div><span class="method method-get">GET</span>/readyz</div>
          <button class="btn btn-secondary" onclick="testEndpoint('/readyz')">Test</button>
        </div>
        <div class="endpoint-card">
          <div><span class="method method-get">GET</span>/api/v1/orders</div>
          <button class="btn btn-secondary" onclick="testEndpoint('/api/v1/orders')">Test</button>
        </div>
      </div>

      <div style="margin-top: 1.5rem;">
        <div style="font-size: 0.85rem; color: var(--text-muted); margin-bottom: 0.4rem;">JSON Response Output:</div>
        <pre id="api-output">// Click any "Test" button above to inspect endpoint response</pre>
      </div>
    </div>

  </main>

  <div class="modal-overlay" id="orderModal">
    <div class="modal">
      <h3 style="font-size: 1.2rem; margin-bottom: 1rem;">Create New Enterprise Order</h3>
      <form id="createOrderForm" onsubmit="submitOrder(event)">
        <div class="form-group">
          <label>Customer Name</label>
          <input type="text" id="custName" placeholder="Acme Logistics Inc." required />
        </div>
        <div class="form-group">
          <label>Item Description</label>
          <input type="text" id="itemDesc" placeholder="Managed Kubernetes Node Cluster" required />
        </div>
        <div class="form-group">
          <label>Quantity</label>
          <input type="number" id="itemQty" value="5" min="1" required />
        </div>
        <div class="form-group">
          <label>Unit Price ($)</label>
          <input type="number" step="0.01" id="itemPrice" value="499.00" required />
        </div>
        <div style="display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.5rem;">
          <button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
          <button type="submit" class="btn">Submit Order</button>
        </div>
      </form>
    </div>
  </div>

  <script>
    async function fetchOrders() {
      try {
        const res = await fetch('/api/v1/orders');
        const data = await res.json();
        const tbody = document.getElementById('orders-list');
        tbody.innerHTML = '';

        if (!data || data.length === 0) {
          tbody.innerHTML = '<tr><td colspan="6" style="text-align: center;">No orders found.</td></tr>';
          return;
        }

        data.forEach(function(o) {
          const tr = document.createElement('tr');
          const totalPrice = (o.price * o.quantity).toFixed(2);
          tr.innerHTML =
            '<td style="font-weight: 600; color: #38bdf8;">' + o.id + '</td>' +
            '<td>' + o.customer_name + '</td>' +
            '<td>' + o.item + '</td>' +
            '<td>' + o.quantity + '</td>' +
            '<td>$' + totalPrice + '</td>' +
            '<td><span class="tag tag-created">' + o.status + '</span></td>';
          tbody.appendChild(tr);
        });
      } catch (err) {
        console.error('Failed fetching orders:', err);
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

    function openModal() {
      document.getElementById('orderModal').style.display = 'flex';
    }

    function closeModal() {
      document.getElementById('orderModal').style.display = 'none';
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
          fetchOrders();
          testEndpoint('/api/v1/orders');
        } else {
          alert('Failed creating order');
        }
      } catch (err) {
        alert('Error: ' + err.message);
      }
    }

    fetchOrders();
  </script>
</body>
</html>
`
