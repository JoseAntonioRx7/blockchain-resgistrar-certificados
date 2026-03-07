const API_URL = 'http://localhost:8080';

// 1. Inicialização do Sistema
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();

    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    if (loginForm) loginForm.addEventListener('submit', loginInstituicao);
    if (registerForm) registerForm.addEventListener('submit', registerCertificate);

    loadDashboard();
});

// 2. Controle de Autenticação e Visibilidade
function checkAuth() {
    const token = localStorage.getItem('instituicao_token');
    const loginSection = document.getElementById('loginSection');
    const registerSection = document.getElementById('registerSection');
    const authHeader = document.getElementById('authHeader'); // Para o botão de sair

    if (token) {
        // Usuário Logado
        if (loginSection) loginSection.style.display = 'none';
        if (registerSection) registerSection.style.display = 'block';
        if (authHeader) authHeader.style.display = 'block';
    } else {
        // Usuário Deslogado
        if (loginSection) loginSection.style.display = 'block';
        if (registerSection) registerSection.style.display = 'none';
        if (authHeader) authHeader.style.display = 'none';
    }
}

function logoutInstituicao() {
    localStorage.removeItem('instituicao_token');
    checkAuth();
    const resultDiv = document.getElementById('loginResult');
    if (resultDiv) resultDiv.innerText = "";
}

// 3. Função de Login
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
            resultDiv.innerText = "Login efetuado com sucesso!";
            resultDiv.style.color = "#10b981";
            checkAuth(); // Isso agora vai mostrar a registerSection corretamente
        } else {
            resultDiv.innerText = data.error || "Erro ao fazer login";
            resultDiv.style.color = "#ef4444";
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao conectar com o servidor.";
    }
}

// 4. Função de Registro
async function registerCertificate(e) {
    e.preventDefault();
    const resultDiv = document.getElementById('registerResult');
    const token = localStorage.getItem('instituicao_token');

    if (!token) {
        resultDiv.innerText = "Erro: Faça login primeiro.";
        resultDiv.style.color = "#ef4444";
        return;
    }

    resultDiv.innerHTML = "<strong>Minerando bloco...</strong>";
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
            resultDiv.innerHTML = `<strong style="color: #10b981;">Sucesso!</strong> ID: ${data.id}`;
            e.target.reset(); 
            loadDashboard(); 
        } else {
            resultDiv.innerText = "Erro: " + (data.error || "Falha");
        }
    } catch (err) {
        resultDiv.innerText = "Erro no servidor.";
    }
}

// 5. Função de Verificação (Pública)
async function verifyCertificate() {
    const hash = document.getElementById('hashInput').value.trim();
    const resultDiv = document.getElementById('verifyResult');

    if (!hash) return;

    try {
        const response = await fetch(`${API_URL}/verify?hash=${encodeURIComponent(hash)}`);
        const data = await response.json();
        
        if (response.ok && !data.error) {
            resultDiv.innerHTML = `
                <div style="border-left: 4px solid #10b981; padding: 10px; background: #f0fdf4;">
                    <strong style="color: #10b981;">Autêntico!</strong><br>
                    <strong>Aluno:</strong> ${data.student_name}<br>
                    <strong>Data:</strong> ${new Date(data.timestamp * 1000).toLocaleString()}
                </div>
            `;
        } else {
            resultDiv.innerHTML = `<strong style="color: #ef4444;">Não encontrado.</strong>`;
        }
    } catch (err) {
        resultDiv.innerText = "Erro na consulta.";
    }
}

// 6. Dashboard
async function loadDashboard() {
    const tbody = document.getElementById('dashboardBody');
    const counter = document.getElementById('countCerts');
    if (!tbody) return;

    try {
        const response = await fetch(`${API_URL}/list`);
        const certs = await response.json();

        tbody.innerHTML = "";
        if (Array.isArray(certs) && certs.length > 0) {
            certs.forEach(cert => {
                const date = cert.timestamp ? new Date(cert.timestamp * 1000).toLocaleDateString() : 'N/A';
                const row = `
                    <tr>
                        <td><strong>${cert.student_name}</strong></td>
                        <td>${cert.course} <br> <small>${cert.institution}</small></td>
                        <td>${date}</td>
                        <td><code>${cert.file_hash.substring(0, 10)}...</code></td>
                        <td>
                            <button class="btn btn-primary" onclick="copyHash('${cert.file_hash}')">Hash</button>
                            <a href="${API_URL}/pdfs/cert_${cert.id}.pdf" target="_blank" class="btn" style="background: #10b981;">PDF</a>
                        </td>
                    </tr>
                `;
                tbody.insertAdjacentHTML('beforeend', row);
            });
            if (counter) counter.innerText = `(${certs.length})`;
        }
    } catch (err) {
        console.error(err);
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
    navigator.clipboard.writeText(hash).then(() => alert("Copiado!"));
}