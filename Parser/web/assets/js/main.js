let selectedFile = null;

document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');
    const uploadBtn = document.getElementById('uploadBtn');
    const fileInfo = document.getElementById('fileInfo');

    // Drag & Drop события
    dropZone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropZone.classList.add('dragover');
    });

    dropZone.addEventListener('dragleave', () => {
        dropZone.classList.remove('dragover');
    });

    dropZone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropZone.classList.remove('dragover');
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFileSelect(files[0]);
        }
    });

    // Клик по зоне открывает диалог выбора файла
    dropZone.addEventListener('click', () => {
        fileInput.click();
    });

    // Обработка выбора файла через input
    fileInput.addEventListener('change', () => {
        if (fileInput.files.length > 0) {
            handleFileSelect(fileInput.files[0]);
        }
    });

    // Нажатие кнопки загрузки
    uploadBtn.addEventListener('click', () => {
        if (selectedFile) {
            uploadFile(selectedFile);
        }
    });

    function handleFileSelect(file) {
        if (!file.name.endsWith('.go')) {
            alert('Пожалуйста, выберите файл с расширением .go');
            return;
        }
        selectedFile = file;
        fileInfo.innerHTML = `<span class="check-icon">✓</span> Готов к загрузке: ${file.name}`;
        fileInfo.classList.add('ready');
        uploadBtn.disabled = false;
    }
});

async function uploadFile(file) {
    const uploadBtn = document.getElementById('uploadBtn');
    const fileInfo = document.getElementById('fileInfo');
    const output = document.getElementById('output');
    const metricsTable = document.getElementById('metricsTable');
    const extendedMetrics = document.getElementById('extendedMetrics');

    // Блокируем кнопку на время загрузки
    uploadBtn.disabled = true;
    output.textContent = `Загрузка файла ${file.name}...`;

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`Ошибка сервера: ${response.status}`);
        }

        const data = await response.json();

        // Строим таблицу и обновляем расширенные метрики
        buildCombinedTable(data);
        updateExtendedMetrics(data);

        // Показываем таблицу и блок метрик
        metricsTable.style.display = 'table';
        extendedMetrics.style.display = 'block';

        output.textContent = `Файл ${file.name} успешно обработан.`;
    } catch (error) {
        console.error('Ошибка:', error);
        output.textContent = 'Ошибка: ' + error.message;
        // Скрываем таблицу и блок метрик при ошибке
        metricsTable.style.display = 'none';
        extendedMetrics.style.display = 'none';
    } finally {
        // Разблокируем кнопку, чтобы можно было загрузить другой файл
        uploadBtn.disabled = false;
    }
}

function buildCombinedTable(data) {
    const operatorsMap = data.operators || {};
    const operandsMap = data.operands || {};

    // Сортировка по убыванию количества, затем по ключу
    const operators = Object.entries(operatorsMap).sort((a, b) => {
        if (b[1] !== a[1]) return b[1] - a[1];
        return a[0].localeCompare(b[0]);
    });
    const operands = Object.entries(operandsMap).sort((a, b) => {
        if (b[1] !== a[1]) return b[1] - a[1];
        return a[0].localeCompare(b[0]);
    });

    const n1 = data.unique_operators !== undefined ? data.unique_operators : operators.length;
    const N1 = data.operators_total;
    const n2 = data.unique_operands !== undefined ? data.unique_operands : operands.length;
    const N2 = data.operands_total;

    const rowsCount = Math.max(operators.length, operands.length);

    const tbody = document.querySelector('#metricsTable tbody');
    tbody.innerHTML = '';

    for (let i = 0; i < rowsCount; i++) {
        const tr = document.createElement('tr');

        if (i < operators.length) {
            const [op, count] = operators[i];
            tr.innerHTML += `
                <td>${i+1}</td>
                <td><code>${escapeHtml(op)}</code></td>
                <td>${count}</td>
            `;
        } else {
            tr.innerHTML += `<td></td><td></td><td></td>`;
        }

        if (i < operands.length) {
            const [operand, count] = operands[i];
            tr.innerHTML += `
                <td>${i+1}</td>
                <td><code>${escapeHtml(operand)}</code></td>
                <td>${count}</td>
            `;
        } else {
            tr.innerHTML += `<td></td><td></td><td></td>`;
        }

        tbody.appendChild(tr);
    }

    // Итоговая строка с метриками
    const totalRow = document.createElement('tr');
    totalRow.classList.add('total-row');
    totalRow.innerHTML = `
        <td>η<sub>1</sub> = ${n1}</td>
        <td></td>
        <td>N<sub>1</sub> = ${N1}</td>
        <td>η<sub>2</sub> = ${n2}</td>
        <td></td>
        <td>N<sub>2</sub> = ${N2}</td>
    `;
    tbody.appendChild(totalRow);
}

function updateExtendedMetrics(data) {
    const n1 = data.unique_operators || 0;
    const n2 = data.unique_operands || 0;
    const N1 = data.operators_total || 0;
    const N2 = data.operands_total || 0;

    const n = n1 + n2;
    const N = N1 + N2;
    // Вычисляем логарифм по основанию 2, округляем до 2 знаков
    const V = (n > 0) ? (N * Math.log2(n)).toFixed(2) : '0.00';

    document.getElementById('dict-n').innerHTML = `Словарь программы: η = ${n1} + ${n2} = ${n}`;
    document.getElementById('length-N').innerHTML = `Длина программы: N = ${N1} + ${N2} = ${N}`;
    document.getElementById('volume-V').innerHTML = `Объём программы: V = ${N} * log<sub>2</sub>(${n}) = ${V}`;
}

function escapeHtml(unsafe) {
    return unsafe.replace(/[&<>"]/g, function(m) {
        if (m === '&') return '&amp;';
        if (m === '<') return '&lt;';
        if (m === '>') return '&gt;';
        if (m === '"') return '&quot;';
        return m;
    });
}