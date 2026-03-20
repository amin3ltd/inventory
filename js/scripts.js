// API Configuration
const API_BASE = window.location.origin + '/api';

// Global state
let products = [];
let categories = [];
let orders = [];
let currentOrderType = 'purchase';
let barChart = null;
let areaChart = null;

// SIDEBAR TOGGLE
let sidebarOpen = false;
const sidebar = document.getElementById('sidebar');

function openSidebar() {
  if (!sidebarOpen) {
    sidebar.classList.add('sidebar-responsive');
    sidebarOpen = true;
  }
}

function closeSidebar() {
  if (sidebarOpen) {
    sidebar.classList.remove('sidebar-responsive');
    sidebarOpen = false;
  }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
  loadDashboard();
  loadCategories();
  updateApiUrl();
});

// View Navigation
function showDashboard() {
  hideAllViews();
  document.getElementById('dashboard-view').style.display = 'block';
  loadDashboard();
}

function showProducts() {
  hideAllViews();
  document.getElementById('products-view').style.display = 'block';
  loadProducts();
}

function showOrders(type) {
  currentOrderType = type;
  hideAllViews();
  document.getElementById('orders-view').style.display = 'block';
  document.getElementById('orders-title').textContent = type.toUpperCase() + ' ORDERS';
  loadOrders(type);
}

function showCategories() {
  hideAllViews();
  document.getElementById('categories-view').style.display = 'block';
  loadCategoriesTable();
}

function showSettings() {
  hideAllViews();
  document.getElementById('settings-view').style.display = 'block';
}

function hideAllViews() {
  document.getElementById('dashboard-view').style.display = 'none';
  document.getElementById('products-view').style.display = 'none';
  document.getElementById('orders-view').style.display = 'none';
  document.getElementById('categories-view').style.display = 'none';
  document.getElementById('settings-view').style.display = 'none';
}

// API Functions
async function fetchAPI(endpoint, options = {}) {
  try {
    const response = await fetch(API_BASE + endpoint, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('API Error:', error);
    alert('Failed to connect to API');
    return null;
  }
}

// Dashboard Functions
async function loadDashboard() {
  const statsResponse = await fetchAPI('/dashboard/stats');
  if (statsResponse && statsResponse.success) {
    document.getElementById('total-products').textContent = statsResponse.data.total_products;
    document.getElementById('purchase-orders').textContent = statsResponse.data.purchase_orders;
    document.getElementById('sales-orders').textContent = statsResponse.data.sales_orders;
    document.getElementById('inventory-alerts').textContent = statsResponse.data.inventory_alerts;
  }

  await loadTopProducts();
  await loadOrdersChart();
}

async function loadTopProducts() {
  const response = await fetchAPI('/dashboard/top-products');
  if (response && response.success && response.data.length > 0) {
    const productNames = response.data.map(p => p.name);
    const productQuantities = response.data.map(p => p.quantity);

    if (barChart) {
      barChart.destroy();
    }

    const barChartOptions = {
      series: [{ data: productQuantities }],
      chart: { type: 'bar', height: 350, toolbar: { show: false } },
      colors: ['#246dec', '#cc3c43', '#367952', '#f5b74f', '#4f35a1'],
      plotOptions: {
        bar: { distributed: true, borderRadius: 4, horizontal: false, columnWidth: '40%' },
      },
      dataLabels: { enabled: false },
      legend: { show: false },
      xaxis: { categories: productNames },
      yaxis: { title: { text: 'Count' } },
    };

    barChart = new ApexCharts(document.querySelector('#bar-chart'), barChartOptions);
    barChart.render();
  }
}

async function loadOrdersChart() {
  const ordersResponse = await fetchAPI('/orders');
  if (ordersResponse && ordersResponse.success) {
    const orders = ordersResponse.data;
    const purchaseCount = orders.filter(o => o.type === 'purchase').length;
    const salesCount = orders.filter(o => o.type === 'sales').length;

    if (areaChart) {
      areaChart.destroy();
    }

    const areaChartOptions = {
      series: [
        { name: 'Purchase Orders', data: [purchaseCount, purchaseCount + 2, purchaseCount + 5] },
        { name: 'Sales Orders', data: [salesCount, salesCount + 1, salesCount + 3] },
      ],
      chart: { height: 350, type: 'area', toolbar: { show: false } },
      colors: ['#4f35a1', '#246dec'],
      dataLabels: { enabled: false },
      stroke: { curve: 'smooth' },
      labels: ['Week 1', 'Week 2', 'Week 3'],
      markers: { size: 0 },
      yaxis: [{ title: { text: 'Purchase Orders' } }, { opposite: true, title: { text: 'Sales Orders' } }],
      tooltip: { shared: true, intersect: false },
    };

    areaChart = new ApexCharts(document.querySelector('#area-chart'), areaChartOptions);
    areaChart.render();
  }
}

// Products Functions
async function loadProducts() {
  const response = await fetchAPI('/products');
  if (response && response.success) {
    products = response.data;
    renderProductsTable();
  }
}

function renderProductsTable() {
  const tbody = document.getElementById('products-tbody');
  tbody.innerHTML = '';
  
  products.forEach(product => {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>${product.name}</td>
      <td>${product.sku}</td>
      <td>${product.category}</td>
      <td>${product.quantity}</td>
      <td>$${product.price.toFixed(2)}</td>
      <td>${product.min_stock}</td>
      <td>
        <button class="btn-icon" onclick="editProduct('${product.id}')">✏️</button>
        <button class="btn-icon" onclick="deleteProduct('${product.id}')">🗑️</button>
      </td>
    `;
    tbody.appendChild(row);
  });
}

function showAddProductModal() {
  document.getElementById('product-modal').style.display = 'block';
  populateCategorySelect();
}

function closeProductModal() {
  document.getElementById('product-modal').style.display = 'none';
  document.getElementById('product-form').reset();
}

function populateCategorySelect() {
  const select = document.getElementById('product-category');
  select.innerHTML = '';
  categories.forEach(cat => {
    const option = document.createElement('option');
    option.value = cat.name;
    option.textContent = cat.name;
    select.appendChild(option);
  });
}

document.getElementById('product-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const product = {
    name: document.getElementById('product-name').value,
    sku: document.getElementById('product-sku').value,
    category: document.getElementById('product-category').value,
    quantity: parseInt(document.getElementById('product-quantity').value),
    price: parseFloat(document.getElementById('product-price').value),
    min_stock: parseInt(document.getElementById('product-min-stock').value),
    description: document.getElementById('product-description').value,
  };

  const response = await fetchAPI('/products', {
    method: 'POST',
    body: JSON.stringify(product),
  });

  if (response && response.success) {
    closeProductModal();
    loadProducts();
    loadDashboard();
    alert('Product added successfully!');
  }
});

async function deleteProduct(id) {
  if (!confirm('Are you sure you want to delete this product?')) return;
  
  const response = await fetchAPI('/products/' + id, { method: 'DELETE' });
  if (response && response.success) {
    loadProducts();
    loadDashboard();
    alert('Product deleted successfully!');
  }
}

function editProduct(id) {
  const product = products.find(p => p.id === id);
  if (!product) return;
  
  document.getElementById('product-name').value = product.name;
  document.getElementById('product-sku').value = product.sku;
  document.getElementById('product-category').value = product.category;
  document.getElementById('product-quantity').value = product.quantity;
  document.getElementById('product-price').value = product.price;
  document.getElementById('product-min-stock').value = product.min_stock;
  document.getElementById('product-description').value = product.description;
  
  showAddProductModal();
}

// Orders Functions
async function loadOrders(type) {
  const response = await fetchAPI('/orders');
  if (response && response.success) {
    orders = response.data.filter(o => o.type === type);
    renderOrdersTable();
  }
}

function renderOrdersTable() {
  const tbody = document.getElementById('orders-tbody');
  tbody.innerHTML = '';
  
  orders.forEach(order => {
    const product = products.find(p => p.id === order.product_id);
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>${order.id.substring(0, 8)}</td>
      <td>${order.type}</td>
      <td>${product ? product.name : 'N/A'}</td>
      <td>${order.quantity}</td>
      <td>$${order.total_price.toFixed(2)}</td>
      <td><span class="status-${order.status}">${order.status}</span></td>
      <td>${new Date(order.created_at).toLocaleDateString()}</td>
    `;
    tbody.appendChild(row);
  });
}

function showAddOrderModal() {
  document.getElementById('order-modal').style.display = 'block';
  populateProductSelect();
}

function closeOrderModal() {
  document.getElementById('order-modal').style.display = 'none';
  document.getElementById('order-form').reset();
}

function populateProductSelect() {
  const select = document.getElementById('order-product');
  select.innerHTML = '';
  products.forEach(product => {
    const option = document.createElement('option');
    option.value = product.id;
    option.textContent = product.name + ' ($' + product.price.toFixed(2) + ')';
    select.appendChild(option);
  });
}

document.getElementById('order-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const order = {
    type: document.getElementById('order-type').value,
    product_id: document.getElementById('order-product').value,
    quantity: parseInt(document.getElementById('order-quantity').value),
    status: 'completed',
  };

  const response = await fetchAPI('/orders', {
    method: 'POST',
    body: JSON.stringify(order),
  });

  if (response && response.success) {
    closeOrderModal();
    loadOrders(currentOrderType);
    loadDashboard();
    alert('Order created successfully!');
  }
});

// Categories Functions
async function loadCategories() {
  const response = await fetchAPI('/categories');
  if (response && response.success) {
    categories = response.data;
  }
}

async function loadCategoriesTable() {
  await loadCategories();
  const tbody = document.getElementById('categories-tbody');
  tbody.innerHTML = '';
  
  categories.forEach(category => {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>${category.name}</td>
      <td>${category.description}</td>
      <td>${new Date(category.created_at).toLocaleDateString()}</td>
    `;
    tbody.appendChild(row);
  });
}

function showAddCategoryModal() {
  document.getElementById('category-modal').style.display = 'block';
}

function closeCategoryModal() {
  document.getElementById('category-modal').style.display = 'none';
  document.getElementById('category-form').reset();
}

document.getElementById('category-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const category = {
    name: document.getElementById('category-name').value,
    description: document.getElementById('category-description').value,
  };

  const response = await fetchAPI('/categories', {
    method: 'POST',
    body: JSON.stringify(category),
  });

  if (response && response.success) {
    closeCategoryModal();
    loadCategoriesTable();
    loadCategories();
    alert('Category added successfully!');
  }
});

// Settings Functions
function updateApiUrl() {
  document.getElementById('api-url').textContent = API_BASE;
}

async function testConnection() {
  const response = await fetchAPI('/products');
  if (response) {
    alert('Connection successful! API is working.');
  } else {
    alert('Connection failed! Please check if the server is running.');
  }
}

// Close modals when clicking outside
window.onclick = function(event) {
  if (event.target.classList.contains('modal')) {
    event.target.style.display = 'none';
  }
}
