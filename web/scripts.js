const API_URL = 'http://localhost:8080';

// 1. Inicialização do Sistema
document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    if (loginForm) {
        loginForm.addEventListener('submit', loginInstituicao);
    }

    if (registerForm) {
        registerForm.addEventListener('submit', registerCertificate);
    }

    // Carrega os dados iniciais
    loadDashboard();
});

// 2. Função de Login
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
            document.getElementById('loginSection').style.display = 'none';
        } else {
            resultDiv.innerText = data.error || "Erro ao fazer login";
            resultDiv.style.color = "#ef4444";
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao conectar com o servidor.";
    }
}

// 3. Função de Registro (Com Proteção JWT)
async function registerCertificate(e) {
    e.preventDefault();
    const resultDiv = document.getElementById('registerResult');
    const token = localStorage.getItem('instituicao_token');

    if (!token) {
        resultDiv.innerText = "Erro: Você precisa fazer login primeiro.";
        resultDiv.style.color = "#ef4444";
        return;
    }

    resultDiv.innerHTML = "<strong>Minerando bloco...</strong> Isso pode levar alguns minutos.";
    const formData = new FormData(e.target);

    try {
        const response = await fetch(`${API_URL}/register`, { 
            method: 'POST', 
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData 
        });
        
        if (response.status === 401) {
            resultDiv.innerText = "Sessão expirada. Faça login novamente.";
            localStorage.removeItem('instituicao_token');
            document.getElementById('loginSection').style.display = 'block';
            return;
        }

        const data = await response.json();
        
        if (response.ok) {
            resultDiv.innerHTML = `
                <strong style="color: #10b981;">Sucesso!</strong><br>
                ID: ${data.id}<br>
                Hash: <small>${data.hash}</small>
            `;
            e.target.reset(); 
            loadDashboard(); 
        } else {
            resultDiv.innerText = "Erro: " + (data.error || "Falha no registro");
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao conectar com o servidor.";
        console.error("Erro no fetch:", err);
    }
}

// 4. Função de Verificação (Pública)
async function verifyCertificate() {
    const hash = document.getElementById('hashInput').value.trim();
    const resultDiv = document.getElementById('verifyResult');

    if (!hash) {
        resultDiv.innerText = "Por favor, insira um hash.";
        return;
    }

    try {
        const response = await fetch(`${API_URL}/verify?hash=${encodeURIComponent(hash)}`);
        const data = await response.json();
        
        if (response.ok && !data.error) {
            resultDiv.innerHTML = `
                <div style="border-left: 4px solid #10b981; padding-left: 10px;">
                    <strong style="color: #10b981;">Certificado Autêntico!</strong><br>
                    <strong>Aluno:</strong> ${data.student_name}<br>
                    <strong>Instituição:</strong> ${data.institution}<br>
                    <strong>Curso:</strong> ${data.course}<br>
                    <strong>Data:</strong> ${new Date(data.timestamp * 1000).toLocaleString()}
                </div>
            `;
        } else {
            resultDiv.innerHTML = `<strong style="color: #ef4444;">Certificado não encontrado.</strong>`;
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao consultar servidor.";
    }
}

// 5. Dashboard
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
        <td><code title="${cert.file_hash}">${cert.file_hash.substring(0, 10)}...</code></td>
        <td>
            <div style="display: flex; gap: 5px;">
                <button class="btn btn-primary" style="padding: 5px 10px; font-size: 0.8rem;" 
                    onclick="copyHash('${cert.file_hash}')">Copiar Hash</button>
                
                <a href="${API_URL}/pdfs/cert_${cert.id}.pdf" target="_blank" class="btn" 
                   style="background-color: #10b981; color: white; padding: 5px 10px; font-size: 0.8rem; text-decoration: none; border-radius: 4px;">
                   Baixar PDF
                </a>
            </div>
        </td>
    </tr>


                    <tr>
                        <td><strong>${cert.student_name}</strong></td>
                        <td>${cert.course} <br> <small>${cert.institution}</small></td>
                        <td>${date}</td>
                        <td><code title="${cert.file_hash}">${cert.file_hash.substring(0, 10)}...</code></td>
                        <td>
                            <button class="btn btn-primary" style="padding: 5px 10px; font-size: 0.8rem;" 
                                onclick="copyHash('${cert.file_hash}')">Copiar Hash</button>
                        </td>
                    </tr>
                `;
                tbody.insertAdjacentHTML('beforeend', row);
            });
            if (counter) counter.innerText = `(${certs.length})`;
        } else {
            tbody.innerHTML = "<tr><td colspan='5' style='text-align:center;padding:20px;'>Nenhum certificado registrado.</td></tr>";
            if (counter) counter.innerText = `(0)`;
        }
    } catch (err) {
        console.error("Erro ao carregar dashboard:", err);
        tbody.innerHTML = "<tr><td colspan='5' style='text-align:center;color:red;'>Erro ao carregar dados.</td></tr>";
    }
}

// 6. Utilitários
function filterTable() {
    const input = document.getElementById("searchInput");
    if (!input) return;
    const filter = input.value.toUpperCase();
    const rows = document.getElementById("dashboardBody").getElementsByTagName("tr");

    for (let row of rows) {
        row.style.display = row.innerText.toUpperCase().includes(filter) ? "" : "none";
    }
}

function copyHash(hash) {
    if (!hash) return;
    navigator.clipboard.writeText(hash).then(() => {
        alert("Hash copiado com sucesso!");
    });
}