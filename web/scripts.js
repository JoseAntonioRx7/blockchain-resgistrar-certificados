const API_URL = 'http://localhost:8080';

// 1. Inicialização do Sistema
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    updateCurrentDate();

    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    if (loginForm) loginForm.addEventListener('submit', loginInstituicao);
    if (registerForm) registerForm.addEventListener('submit', registerCertificate);

    loadDashboard();
});

// Atualiza a data no cabeçalho
function updateCurrentDate() {
    const dateElement = document.getElementById('currentDate');
    if (dateElement) {
        const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
        dateElement.innerText = new Date().toLocaleDateString('pt-BR', options);
    }
}

// 2. Controle de Autenticação e Interface
function checkAuth() {
    const token = localStorage.getItem('instituicao_token');
    
    // Elementos de Seção
    const loginSection = document.getElementById('loginSection');
    const dashboardStats = document.getElementById('dashboardStats');
    const mainActions = document.getElementById('mainActions');
    const authHeader = document.getElementById('authHeader');
    const pageTitle = document.getElementById('pageTitle');

    if (token) {
        // Logado
        if (loginSection) loginSection.style.display = 'none';
        if (dashboardStats) dashboardStats.style.display = 'grid';
        if (mainActions) mainActions.style.display = 'grid';
        if (authHeader) authHeader.style.display = 'flex';
        if (pageTitle) pageTitle.innerText = "Visão Geral do Sistema";
    } else {
        // Deslogado
        if (loginSection) loginSection.style.display = 'block';
        if (dashboardStats) dashboardStats.style.display = 'none';
        if (mainActions) mainActions.style.display = 'none';
        if (authHeader) authHeader.style.display = 'none';
        if (pageTitle) pageTitle.innerText = "Acesso Restrito";
    }
}

function logoutInstituicao() {
    localStorage.removeItem('instituicao_token');
    checkAuth();
    const resultDiv = document.getElementById('loginResult');
    if (resultDiv) resultDiv.innerText = "";
}

// 3. Login
async function loginInstituicao(e) {
    e.preventDefault();
    const user = document.getElementById('loginUser').value;
    const pass = document.getElementById('loginPass').value;
    const resultDiv = document.getElementById('loginResult');

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: user, password: pass })
        });

        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('instituicao_token', data.token);
            resultDiv.innerHTML = `<span style="color: #10b981;">Acesso autorizado. Carregando...</span>`;
            
            // Pequeno delay para efeito visual
            setTimeout(() => {
                checkAuth();
                loadDashboard();
            }, 800);
        } else {
            resultDiv.innerText = data.error || "Credenciais inválidas";
            resultDiv.style.color = "#ef4444";
        }
    } catch (err) {
        resultDiv.innerText = "Erro na conexão com o nó blockchain.";
    }
}

// 4. Registro de Certificado
async function registerCertificate(e) {
    e.preventDefault();
    const resultDiv = document.getElementById('registerResult');
    const token = localStorage.getItem('instituicao_token');
    const btn = e.target.querySelector('button');

    if (!token) return;

    // Feedback visual de mineração
    const originalBtnText = btn.innerHTML;
    btn.innerHTML = `<i class="fas fa-circle-notch fa-spin"></i> Minerando Bloco...`;
    btn.disabled = true;

    const formData = new FormData(e.target);

    try {
        const response = await fetch(`${API_URL}/register`, { 
            method: 'POST', 
            headers: { 'Authorization': `Bearer ${token}` },
            body: formData 
        });
        
        if (response.status === 401) {
            logoutInstituicao();
            return;
        }

        const data = await response.json();
        
        if (response.ok) {
            resultDiv.innerHTML = `
                <div style="background: #ecfdf5; color: #059669; padding: 12px; border-radius: 8px; margin-top: 10px;">
                    <i class="fas fa-check-circle"></i> <strong>Sucesso!</strong> Registro minerado com ID: ${data.id.substring(0,8)}
                </div>`;
            e.target.reset(); 
            loadDashboard(); 
        } else {
            resultDiv.innerHTML = `<span style="color: #ef4444;">Erro: ${data.error}</span>`;
        }
    } catch (err) {
        resultDiv.innerText = "Erro fatal no processo de registro.";
    } finally {
        btn.innerHTML = originalBtnText;
        btn.disabled = false;
    }
}

// 5. Verificação Pública
async function verifyCertificate() {
    const hash = document.getElementById('hashInput').value.trim();
    const resultDiv = document.getElementById('verifyResult');

    if (!hash) return;

    resultDiv.innerHTML = "Buscando na rede...";

    try {
        const response = await fetch(`${API_URL}/verify?hash=${encodeURIComponent(hash)}`);
        const data = await response.json();
        
        if (response.ok && !data.error) {
            resultDiv.innerHTML = `
                <div class="result-success" style="border-left: 4px solid #10b981; padding: 15px; background: #ffffff; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.05);">
                    <h4 style="color: #10b981; margin-bottom: 5px;"><i class="fas fa-check-shield"></i> Documento Autêntico</h4>
                    <p style="font-size: 0.9rem; margin: 2px 0;"><strong>Aluno:</strong> ${data.student_name}</p>
                    <p style="font-size: 0.9rem; margin: 2px 0;"><strong>Curso:</strong> ${data.course}</p>
                    <p style="font-size: 0.8rem; color: #666;"><strong>Data de Registro:</strong> ${new Date(data.timestamp * 1000).toLocaleString()}</p>
                </div>
            `;
        } else {
            resultDiv.innerHTML = `
                <div style="padding: 15px; background: #fff5f5; border-radius: 12px; border: 1px solid #feb2b2; color: #c53030;">
                    <i class="fas fa-times-circle"></i> Hash não encontrado ou alterado.
                </div>`;
        }
    } catch (err) {
        resultDiv.innerText = "Erro na rede de validação.";
    }
}

// 6. Atualização do Dashboard e Stats
async function loadDashboard() {
    const tbody = document.getElementById('dashboardBody');
    const counter = document.getElementById('countCerts');
    const statTotal = document.getElementById('statTotal');
    const statBlocks = document.getElementById('statBlocks');

    if (!tbody) return;

    try {
        const response = await fetch(`${API_URL}/list`);
        const certs = await response.json();

        tbody.innerHTML = "";
        
        if (Array.isArray(certs)) {
            // Atualiza Stats do Topo
            if (statTotal) statTotal.innerText = certs.length;
            if (statBlocks) statBlocks.innerText = Math.ceil(certs.length / 1); // Simulação de 1 cert por bloco
            if (counter) counter.innerText = `${certs.length} certificados`;

            certs.reverse().forEach(cert => {
                const date = cert.timestamp ? new Date(cert.timestamp * 1000).toLocaleDateString() : 'N/A';
                const row = `
                    <tr>
                        <td>
                            <div style="display: flex; align-items: center; gap: 10px;">
                                <div style="width: 32px; height: 32px; background: #f3f4f6; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 0.8rem; color: #5e003f; font-weight: bold;">
                                    ${cert.student_name.charAt(0)}
                                </div>
                                <strong>${cert.student_name}</strong>
                            </div>
                        </td>
                        <td>
                            <div style="font-size: 0.9rem;">${cert.course}</div>
                            <div style="font-size: 0.75rem; color: #999;">${cert.institution}</div>
                        </td>
                        <td style="font-size: 0.85rem; color: #666;">${date}</td>
                        <td><code title="${cert.file_hash}">${cert.file_hash.substring(0, 8)}...</code></td>
                        <td style="text-align: right;">
                            <button class="btn" style="padding: 6px 12px; font-size: 0.75rem; background: #f3f4f6; color: #444;" onclick="copyHash('${cert.file_hash}')">
                                <i class="fas fa-copy"></i>
                            </button>
                            <a href="${API_URL}/pdfs/cert_${cert.id}.pdf" target="_blank" class="btn" style="padding: 6px 12px; font-size: 0.75rem; background: #10b981; color: white;">
                                <i class="fas fa-file-pdf"></i> PDF
                            </a>
                        </td>
                    </tr>
                `;
                tbody.insertAdjacentHTML('beforeend', row);
            });
        }
    } catch (err) {
        console.error("Erro ao carregar dashboard:", err);
    }
}

function filterTable() {
    const filter = document.getElementById("searchInput").value.toUpperCase();
    const rows = document.getElementById("dashboardBody").getElementsByTagName("tr");
    for (let row of rows) {
        row.style.display = row.innerText.toUpperCase().includes(filter) ? "" : "none";
    }
}

function copyHash(hash) {
    navigator.clipboard.writeText(hash).then(() => {
        // Feedback simples de cópia
        alert("Hash copiado para a área de transferência!");
    });
}