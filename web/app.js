const API_URL = '/api';

// Загрузка записей при старте
document.addEventListener('DOMContentLoaded', () => {
    loadItems();
    
    // Установка текущей даты по умолчанию
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('date').value = today;
    
    // Установка диапазона для аналитики (последние 30 дней)
    const thirtyDaysAgo = new Date();
    thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
    document.getElementById('analyticsFrom').value = thirtyDaysAgo.toISOString().split('T')[0];
    document.getElementById('analyticsTo').value = today;
});

// Обработка формы добавления
document.getElementById('itemForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const item = {
        type: document.getElementById('type').value,
        amount: parseFloat(document.getElementById('amount').value),
        category: document.getElementById('category').value,
        date: new Date(document.getElementById('date').value).toISOString()
    };
    
    try {
        const response = await fetch(`${API_URL}/items`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        
        if (response.ok) {
            alert('Запись добавлена!');
            document.getElementById('itemForm').reset();
            document.getElementById('date').value = new Date().toISOString().split('T')[0];
            loadItems();
        } else {
            const error = await response.json();
            alert('Ошибка: ' + error.error);
        }
    } catch (error) {
        alert('Ошибка соединения: ' + error.message);
    }
});

// Загрузка списка записей
async function loadItems() {
    const from = document.getElementById('filterFrom').value;
    const to = document.getElementById('filterTo').value;
    
    let url = `${API_URL}/items`;
    const params = new URLSearchParams();
    
    if (from) params.append('from', new Date(from).toISOString());
    if (to) params.append('to', new Date(to).toISOString());
    
    if (params.toString()) url += '?' + params.toString();
    
    try {
        const response = await fetch(url);
        const items = await response.json();
        
        const tbody = document.getElementById('itemsBody');
        tbody.innerHTML = '';
        
        if (!items || items.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Нет записей</td></tr>';
            return;
        }
        
        items.forEach(item => {
            const row = document.createElement('tr');
            const date = new Date(item.date).toLocaleDateString('ru-RU');
            const typeClass = item.type === 'income' ? 'income' : 'expense';
            const typeText = item.type === 'income' ? 'Доход' : 'Расход';
            
            row.innerHTML = `
                <td>${item.id}</td>
                <td class="${typeClass}">${typeText}</td>
                <td>${item.amount.toFixed(2)} ₽</td>
                <td>${item.category}</td>
                <td>${date}</td>
                <td>
                    <button class="btn btn-edit" onclick="editItem(${item.id})">✏️</button>
                    <button class="btn btn-danger" onclick="deleteItem(${item.id})">🗑️</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        alert('Ошибка загрузки данных: ' + error.message);
    }
}

// Загрузка аналитики
async function loadAnalytics() {
    const from = document.getElementById('analyticsFrom').value;
    const to = document.getElementById('analyticsTo').value;
    
    if (!from || !to) {
        alert('Укажите период для аналитики');
        return;
    }
    
    const url = `${API_URL}/analytics?from=${new Date(from).toISOString()}&to=${new Date(to).toISOString()}`;
    
    try {
        const response = await fetch(url);
        const analytics = await response.json();
        
        const container = document.getElementById('analyticsResult');
        container.innerHTML = `
            <div class="analytics-item">
                <h3>Сумма</h3>
                <p>${analytics.sum.toFixed(2)} ₽</p>
            </div>
            <div class="analytics-item">
                <h3>Среднее</h3>
                <p>${analytics.avg.toFixed(2)} ₽</p>
            </div>
            <div class="analytics-item">
                <h3>Количество</h3>
                <p>${analytics.count}</p>
            </div>
            <div class="analytics-item">
                <h3>Медиана</h3>
                <p>${analytics.median.toFixed(2)} ₽</p>
            </div>
            <div class="analytics-item">
                <h3>90-й перцентиль</h3>
                <p>${analytics.percentile_90.toFixed(2)} ₽</p>
            </div>
        `;
    } catch (error) {
        alert('Ошибка загрузки аналитики: ' + error.message);
    }
}

// Редактирование записи
async function editItem(id) {
    try {
        const response = await fetch(`${API_URL}/items/${id}`);
        const item = await response.json();
        
        document.getElementById('editId').value = item.id;
        document.getElementById('editType').value = item.type;
        document.getElementById('editAmount').value = item.amount;
        document.getElementById('editCategory').value = item.category;
        document.getElementById('editDate').value = new Date(item.date).toISOString().split('T')[0];
        
        document.getElementById('editModal').style.display = 'block';
    } catch (error) {
        alert('Ошибка загрузки записи: ' + error.message);
    }
}

// Обработка формы редактирования
document.getElementById('editForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const id = document.getElementById('editId').value;
    const item = {
        type: document.getElementById('editType').value,
        amount: parseFloat(document.getElementById('editAmount').value),
        category: document.getElementById('editCategory').value,
        date: new Date(document.getElementById('editDate').value).toISOString()
    };
    
    try {
        const response = await fetch(`${API_URL}/items/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        
        if (response.ok) {
            alert('Запись обновлена!');
            closeEditModal();
            loadItems();
        } else {
            const error = await response.json();
            alert('Ошибка: ' + error.error);
        }
    } catch (error) {
        alert('Ошибка соединения: ' + error.message);
    }
});

// Удаление записи
async function deleteItem(id) {
    if (!confirm('Удалить эту запись?')) return;
    
    try {
        const response = await fetch(`${API_URL}/items/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('Запись удалена!');
            loadItems();
        } else {
            alert('Ошибка удаления');
        }
    } catch (error) {
        alert('Ошибка соединения: ' + error.message);
    }
}

// Закрытие модального окна
function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
}

// Сброс фильтров
function clearFilters() {
    document.getElementById('filterFrom').value = '';
    document.getElementById('filterTo').value = '';
    loadItems();
}

// Закрытие модального окна при клике вне его
window.onclick = function(event) {
    const modal = document.getElementById('editModal');
    if (event.target === modal) {
        closeEditModal();
    }
}
